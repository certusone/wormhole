package ethereum

import (
	"context"
	"fmt"
	"github.com/certusone/wormhole/bridge/pkg/p2p"
	gossipv1 "github.com/certusone/wormhole/bridge/pkg/proto/gossip/v1"
	"math/big"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	"github.com/certusone/wormhole/bridge/pkg/common"
	"github.com/certusone/wormhole/bridge/pkg/ethereum/abi"
	"github.com/certusone/wormhole/bridge/pkg/readiness"
	"github.com/certusone/wormhole/bridge/pkg/supervisor"
	"github.com/certusone/wormhole/bridge/pkg/vaa"
)

var (
	ethConnectionErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wormhole_eth_connection_errors_total",
			Help: "Total number of Ethereum connection errors (either during initial connection or while watching)",
		}, []string{"reason"})

	ethMessagesObserved = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "wormhole_eth_messages_observed_total",
			Help: "Total number of Eth messages observed (pre-confirmation)",
		})
	ethMessagesConfirmed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "wormhole_eth_messages_confirmed_total",
			Help: "Total number of Eth messages verified (post-confirmation)",
		})
	guardianSetChangesConfirmed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "wormhole_eth_guardian_set_changes_confirmed_total",
			Help: "Total number of guardian set changes verified (we only see confirmed ones to begin with)",
		})
	currentEthHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "wormhole_eth_current_height",
			Help: "Current Ethereum block height",
		})
	queryLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "wormhole_eth_query_latency",
			Help: "Latency histogram for Ethereum calls (note that most interactions are streaming queries, NOT calls, and we cannot measure latency for those",
		}, []string{"operation"})
)

func init() {
	prometheus.MustRegister(ethConnectionErrors)
	prometheus.MustRegister(ethMessagesObserved)
	prometheus.MustRegister(ethMessagesConfirmed)
	prometheus.MustRegister(guardianSetChangesConfirmed)
	prometheus.MustRegister(currentEthHeight)
	prometheus.MustRegister(queryLatency)
}

type (
	EthBridgeWatcher struct {
		url    string
		bridge eth_common.Address

		pendingLocks      map[eth_common.Hash]*pendingLock
		pendingLocksGuard sync.Mutex

		lockChan chan *common.MessagePublication
		setChan  chan *common.GuardianSet
	}

	pendingLock struct {
		lock   *common.MessagePublication
		height uint64
	}
)

func NewEthBridgeWatcher(url string, bridge eth_common.Address, lockEvents chan *common.MessagePublication, setEvents chan *common.GuardianSet) *EthBridgeWatcher {
	return &EthBridgeWatcher{url: url, bridge: bridge, lockChan: lockEvents, setChan: setEvents, pendingLocks: map[eth_common.Hash]*pendingLock{}}
}

func (e *EthBridgeWatcher) Run(ctx context.Context) error {
	// Initialize gossip metrics (we want to broadcast the address even if we're not yet syncing)
	p2p.DefaultRegistry.SetNetworkStats(vaa.ChainIDEthereum, &gossipv1.Heartbeat_Network{
		BridgeAddress: e.bridge.Hex(),
	})

	timeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	c, err := ethclient.DialContext(timeout, e.url)
	if err != nil {
		ethConnectionErrors.WithLabelValues("dial_error").Inc()
		return fmt.Errorf("dialing eth client failed: %w", err)
	}

	f, err := abi.NewAbiFilterer(e.bridge, c)
	if err != nil {
		return fmt.Errorf("could not create wormhole bridge filter: %w", err)
	}

	caller, err := abi.NewAbiCaller(e.bridge, c)
	if err != nil {
		panic(err)
	}

	// Timeout for initializing subscriptions
	timeout, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Subscribe to new message publications
	messageC := make(chan *abi.AbiLogMessagePublished, 2)
	messageSub, err := f.WatchLogMessagePublished(&bind.WatchOpts{Context: timeout}, messageC, nil)
	if err != nil {
		ethConnectionErrors.WithLabelValues("subscribe_error").Inc()
		return fmt.Errorf("failed to subscribe to message publication events: %w", err)
	}

	// Subscribe to guardian set changes
	guardianSetC := make(chan *abi.AbiGuardianSetAdded, 2)
	guardianSetEvent, err := f.WatchGuardianSetAdded(&bind.WatchOpts{Context: timeout}, guardianSetC, nil)
	if err != nil {
		ethConnectionErrors.WithLabelValues("subscribe_error").Inc()
		return fmt.Errorf("failed to subscribe to guardian set events: %w", err)
	}

	errC := make(chan error)
	logger := supervisor.Logger(ctx)

	// Get initial validator set from Ethereum. We could also fetch it from Solana,
	// because both sets are synchronized, we simply made an arbitrary decision to use Ethereum.
	timeout, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	idx, gs, err := FetchCurrentGuardianSet(timeout, e.url, e.bridge)
	if err != nil {
		ethConnectionErrors.WithLabelValues("guardian_set_fetch_error").Inc()
		return fmt.Errorf("failed requesting guardian set from Ethereum: %w", err)
	}
	logger.Info("initial guardian set fetched", zap.Any("value", gs), zap.Uint32("index", idx))
	e.setChan <- &common.GuardianSet{
		Keys:  gs.Keys,
		Index: idx,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-messageSub.Err():
				ethConnectionErrors.WithLabelValues("subscription_error").Inc()
				errC <- fmt.Errorf("error while processing message publication subscription: %w", e)
				return
			case e := <-guardianSetEvent.Err():
				ethConnectionErrors.WithLabelValues("subscription_error").Inc()
				errC <- fmt.Errorf("error while processing guardian set subscription: %w", e)
				return
			case ev := <-messageC:
				// Request timestamp for block
				msm := time.Now()
				timeout, cancel = context.WithTimeout(ctx, 15*time.Second)
				b, err := c.BlockByNumber(timeout, big.NewInt(int64(ev.Raw.BlockNumber)))
				cancel()
				queryLatency.WithLabelValues("block_by_number").Observe(time.Since(msm).Seconds())

				if err != nil {
					ethConnectionErrors.WithLabelValues("block_by_number_error").Inc()
					errC <- fmt.Errorf("failed to request timestamp for block %d: %w", ev.Raw.BlockNumber, err)
					return
				}

				lock := &common.MessagePublication{
					TxHash:           ev.Raw.TxHash,
					Timestamp:        time.Unix(int64(b.Time()), 0),
					Nonce:            ev.Nonce,
					Sequence:         ev.Sequence,
					EmitterChain:     vaa.ChainIDEthereum,
					EmitterAddress:   PadAddress(ev.Sender),
					Payload:          ev.Payload,
					ConsistencyLevel: ev.ConsistencyLevel,
				}

				logger.Info("found new message publication transaction", zap.Stringer("tx", ev.Raw.TxHash),
					zap.Uint64("block", ev.Raw.BlockNumber))

				ethMessagesObserved.Inc()

				e.pendingLocksGuard.Lock()
				e.pendingLocks[ev.Raw.TxHash] = &pendingLock{
					lock:   lock,
					height: ev.Raw.BlockNumber,
				}
				e.pendingLocksGuard.Unlock()
			case ev := <-guardianSetC:
				logger.Info("guardian set has changed, fetching new value",
					zap.Uint32("new_index", ev.Index))

				guardianSetChangesConfirmed.Inc()

				msm := time.Now()
				timeout, cancel = context.WithTimeout(ctx, 15*time.Second)
				gs, err := caller.GetGuardianSet(&bind.CallOpts{Context: timeout}, ev.Index)
				cancel()
				queryLatency.WithLabelValues("get_guardian_set").Observe(time.Since(msm).Seconds())
				if err != nil {
					// We failed to process the guardian set update and are now out of sync with the chain.
					// Recover by crashing the runnable, which causes the guardian set to be re-fetched.
					errC <- fmt.Errorf("error requesting new guardian set value for %d: %w", ev.Index, err)
					return
				}

				logger.Info("new guardian set fetched", zap.Any("value", gs), zap.Uint32("index", ev.Index))
				e.setChan <- &common.GuardianSet{
					Keys:  gs.Keys,
					Index: ev.Index,
				}
			}
		}
	}()

	// Watch headers
	headSink := make(chan *types.Header, 2)
	headerSubscription, err := c.SubscribeNewHead(ctx, headSink)
	if err != nil {
		return fmt.Errorf("failed to subscribe to header events: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-headerSubscription.Err():
				errC <- fmt.Errorf("error while processing header subscription: %w", e)
				return
			case ev := <-headSink:
				start := time.Now()
				logger.Info("processing new header", zap.Stringer("block", ev.Number))
				currentEthHeight.Set(float64(ev.Number.Int64()))
				readiness.SetReady(common.ReadinessEthSyncing)
				p2p.DefaultRegistry.SetNetworkStats(vaa.ChainIDEthereum, &gossipv1.Heartbeat_Network{
					Height:        ev.Number.Int64(),
					BridgeAddress: e.bridge.Hex(),
				})

				e.pendingLocksGuard.Lock()

				blockNumberU := ev.Number.Uint64()
				for hash, pLock := range e.pendingLocks {

					// Transaction was dropped and never picked up again
					if pLock.height+4*uint64(pLock.lock.ConsistencyLevel) <= blockNumberU {
						logger.Debug("observation timed out", zap.Stringer("tx", pLock.lock.TxHash),
							zap.Stringer("block", ev.Number))
						delete(e.pendingLocks, hash)
						continue
					}

					// Transaction is now ready
					if pLock.height+uint64(pLock.lock.ConsistencyLevel) <= ev.Number.Uint64() {
						logger.Debug("observation confirmed", zap.Stringer("tx", pLock.lock.TxHash),
							zap.Stringer("block", ev.Number))
						delete(e.pendingLocks, hash)
						e.lockChan <- pLock.lock
						ethMessagesConfirmed.Inc()
					}
				}

				e.pendingLocksGuard.Unlock()
				logger.Info("processed new header", zap.Stringer("block", ev.Number),
					zap.Duration("took", time.Since(start)))
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

// Fetch the current guardian set ID and guardian set from the chain.
func FetchCurrentGuardianSet(ctx context.Context, rpcURL string, bridgeContract eth_common.Address) (uint32, *abi.StructsGuardianSet, error) {
	c, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return 0, nil, fmt.Errorf("dialing eth client failed: %w", err)
	}

	caller, err := abi.NewAbiCaller(bridgeContract, c)
	if err != nil {
		panic(err)
	}

	opts := &bind.CallOpts{Context: ctx}

	currentIndex, err := caller.GetCurrentGuardianSetIndex(opts)
	if err != nil {
		return 0, nil, fmt.Errorf("error requesting current guardian set index: %w", err)
	}

	gs, err := caller.GetGuardianSet(opts, currentIndex)
	if err != nil {
		return 0, nil, fmt.Errorf("error requesting current guardian set value: %w", err)
	}

	return currentIndex, &gs, nil
}
