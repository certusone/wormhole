package ethereum

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/certusone/wormhole/bridge/pkg/common"
	"github.com/certusone/wormhole/bridge/pkg/supervisor"
	"github.com/certusone/wormhole/bridge/pkg/vaa"
	"github.com/dfuse-io/solana-go"
	"github.com/dfuse-io/solana-go/rpc"
	"github.com/dfuse-io/solana-go/rpc/ws"
	eth_common "github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type SolanaWatcher struct {
	bridge    solana.PublicKey
	url       string
	lockEvent chan *common.ChainLock
}

func NewSolanaWatcher(wsUrl string, bridgeAddress solana.PublicKey, lockEvents chan *common.ChainLock) *SolanaWatcher {
	return &SolanaWatcher{bridge: bridgeAddress, url: wsUrl, lockEvent: lockEvents}
}

func (s *SolanaWatcher) Run(ctx context.Context) error {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	c, err := ws.Dial(tCtx, s.url)
	if err != nil {
		return fmt.Errorf("failed to connect to solana ws: %w", err)
	}
	defer c.Close()

	logger := supervisor.Logger(ctx)

	sub, err := c.ProgramSubscribe(s.bridge, rpc.CommitmentMax)
	if err != nil {
		return fmt.Errorf("failed to subscribe to program: %w", err)
	}
	go func() {
		<-ctx.Done()
		sub.Unsubscribe()
	}()

	logger.Info("watching for on-chain events")
	for {
		updateRaw, err := sub.Recv()
		if err != nil {
			return err
		}
		programUpdate := updateRaw.(*ws.ProgramResult)
		data := programUpdate.Value.Account.Data

		// 1184 is the size of a TransferOutProposal as determined by Rust code `size_of::<TransferOutProposal>`
		if len(data) != 1184 {
			logger.Debug(
				"saw update to non-transfer-proposal wormhole account",
				zap.Stringer("account", programUpdate.Value.PubKey),
				zap.Uint64("slot", programUpdate.Context.Slot),
			)
			continue
		}

		proposal, err := ParseTransferOutProposal(data)
		if err != nil {
			logger.Warn(
				"failed to parse transfer proposal",
				zap.Stringer("account", programUpdate.Value.PubKey),
				zap.Uint64("slot", programUpdate.Context.Slot),
				zap.Error(err),
			)
			continue
		}

		// VAA submitted
		if proposal.VaaTime.Unix() == 0 {
			continue
		}

		var txHash eth_common.Hash
		copy(txHash[:], programUpdate.Value.PubKey[:])

		lock := &common.ChainLock{
			TxHash:        txHash,
			Timestamp:     proposal.LockupTime,
			Nonce:         proposal.Nonce,
			SourceAddress: proposal.SourceAddress,
			TargetAddress: proposal.ForeignAddress,
			SourceChain:   vaa.ChainIDSolana,
			TargetChain:   proposal.ToChainID,
			TokenChain:    proposal.Asset.Chain,
			TokenAddress:  proposal.Asset.Address,
			TokenDecimals: proposal.Asset.Decimals,
			Amount:        proposal.Amount,
		}
		logger.Info("found new lockup transaction", zap.Stringer("lockup_address", programUpdate.Value.PubKey))
		s.lockEvent <- lock
	}
}

type (
	TransferOutProposal struct {
		Amount           *big.Int
		ToChainID        vaa.ChainID
		SourceAddress    vaa.Address
		ForeignAddress   vaa.Address
		Asset            vaa.AssetMeta
		Nonce            uint32
		VAA              [1001]byte
		VaaTime          time.Time
		LockupTime       time.Time
		PokeCounter      uint8
		SignatureAccount solana.PublicKey
	}
)

func ParseTransferOutProposal(data []byte) (*TransferOutProposal, error) {
	prop := &TransferOutProposal{}
	r := bytes.NewBuffer(data)

	var amountBytes [32]byte
	if n, err := r.Read(amountBytes[:]); err != nil || n != 32 {
		return nil, fmt.Errorf("failed to read amount: %w", err)
	}
	// Reverse (little endian -> big endian)
	for i := 0; i < len(amountBytes)/2; i++ {
		amountBytes[i], amountBytes[len(amountBytes)-i-1] = amountBytes[len(amountBytes)-i-1], amountBytes[i]
	}
	prop.Amount = new(big.Int).SetBytes(amountBytes[:])

	if err := binary.Read(r, binary.LittleEndian, &prop.ToChainID); err != nil {
		return nil, fmt.Errorf("failed to read to chain id: %w", err)
	}

	if n, err := r.Read(prop.SourceAddress[:]); err != nil || n != 32 {
		return nil, fmt.Errorf("failed to read source address: %w", err)
	}

	if n, err := r.Read(prop.ForeignAddress[:]); err != nil || n != 32 {
		return nil, fmt.Errorf("failed to read source address: %w", err)
	}

	assetMeta := vaa.AssetMeta{}
	if n, err := r.Read(assetMeta.Address[:]); err != nil || n != 32 {
		return nil, fmt.Errorf("failed to read asset meta address: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &assetMeta.Chain); err != nil {
		return nil, fmt.Errorf("failed to read asset meta chain: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &assetMeta.Decimals); err != nil {
		return nil, fmt.Errorf("failed to read asset meta decimals: %w", err)
	}
	prop.Asset = assetMeta

	// Skip alignment byte
	r.Next(1)

	if err := binary.Read(r, binary.LittleEndian, &prop.Nonce); err != nil {
		return nil, fmt.Errorf("failed to read nonce: %w", err)
	}

	if n, err := r.Read(prop.VAA[:]); err != nil || n != 1001 {
		return nil, fmt.Errorf("failed to read vaa: %w", err)
	}

	// Skip alignment bytes
	r.Next(3)

	var vaaTime uint32
	if err := binary.Read(r, binary.LittleEndian, &vaaTime); err != nil {
		return nil, fmt.Errorf("failed to read vaa time: %w", err)
	}
	prop.VaaTime = time.Unix(int64(vaaTime), 0)

	var lockupTime uint32
	if err := binary.Read(r, binary.LittleEndian, &lockupTime); err != nil {
		return nil, fmt.Errorf("failed to read lockup time: %w", err)
	}
	prop.LockupTime = time.Unix(int64(lockupTime), 0)

	if err := binary.Read(r, binary.LittleEndian, &prop.PokeCounter); err != nil {
		return nil, fmt.Errorf("failed to read poke counter: %w", err)
	}

	if n, err := r.Read(prop.SignatureAccount[:]); err != nil || n != 32 {
		return nil, fmt.Errorf("failed to read signature account: %w", err)
	}

	return prop, nil
}
