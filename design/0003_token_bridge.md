# Token Bridge App

[TOC]

## Objective

To use the Wormhole message passing protocol to transfer tokens between different connected chains.

## Background

The decentralized finance ecosystem is developing into a direction where different chains with different strengths
become the home for various protocols. However, a token is usually only minted on a single chain and therefore
disconnected from the ecosystem and protocols on other chains.

Each chain typically has one de-facto standard for token issuance, like ERC-20 on Ethereum and SPL on Solana. Those
standards, while not identical, all implement a similar interface with the same concepts like owning, minting,
transferring and burning tokens.

To connect chains, a token would ideally have a native mint - the original token on the chain it was originally created
on - and a wrapped version on other chains that represents ownership of the native token.

While the Wormhole messaging protocol provides a way to attest and transfer messages between chains which could
technically be used to implement bridging for individual tokens, this would require manual engineering effort for each
token and create incompatible protocols with bad UX.

## Goals

We want to implement a generalized token bridge using the Wormhole message passing protocol that is able to bridge any
standards-compliant token between chains, creating unique wrapped representations on each connected chain on demand.

* Allow transfer of standards-compliant tokens between chains.
* Allow creation of wrapped assets.
* Use a universal token representation that is compatible with most VM data types.

## Non-Goals

* Support fee-burning / rebasing / non-standard tokens.
* Manage chain-specific token metadata that isn't broadly applicable to all chains.
* Automatically relay token transfer messages to the target chain.

## Overview

On each chain of the token bridge network there will be a token bridge endpoint program.

These programs will manage authorization of payloads (emitter filtering), wrapped representations of foreign chain
tokens ("Wrapped Assets") and custody locked tokens.

## Detailed Design

For outbound transfers, the contracts will have a lock method that either locks up a native token and produces a
respective Transfer message that is posted to Wormhole, or burns a wrapped token and produces/posts said message.

For inbound transfers they can consume, verify and process Wormhole messages containing a token bridge payload.

There will be four different payloads:

* `Transfer` - Will trigger the release of locked tokens or minting of wrapped tokens.
* `AssetMeta` - Attests asset metadata (required before the first transfer).
* `RegisterChain` - Register the token bridge contract (emitter address) for a foreign chain.
* `UpgradeContract` - Upgrade the contract.

Since anyone can use Wormhole to publish messages that match the payload format of the token bridge, an authorization
payload needs to be implemented. This is done using an `(emitter_chain, emitter_address)` tuple. Every endpoint of the
token bridge needs to know the addresses of the respective other endpoints on other chains. This registration of token
bridge endpoints is implemented via `RegisterChain` where a `(chain_id, emitter_address)` tuple can be registered. Only
one endpoint can be registered per chain. Endpoints are immutable. This payload will only be accepted if the emitter is
the hardcoded governance contract.

In order to transfer assets to another chain, a user needs to call the `transfer` method of the bridge contract with the
recipient details and respective fee they are willing to pay. The contract will either hold the tokens in a custody
account (in case it is a native token) or burn wrapped assets. Wrapped assets can be burned because they can be freely
minted once tokens are transferred back and this way the total supply can indicate the total amount of tokens currently
held on this chain. After the lockup the contract will post a `Transfer` payload message to Wormhole. Once the message
has been signed by the guardians, it can be posted to the target chain of the transfer. The target chain will then
either release native tokens from custody or mint a wrapped asset depending on whether it's a native token there. The
program will keep track of consumed message digests for replay prevention.

Since the method for posting a VAA to the token bridge is authorized by the message signature itself, anyone can post
any message. The `completeTransfer` method will accept a fee recipient. In case that field is set, the fee amount
specified will be sent to the fee recipient and the remainder of the amount to the intended receiver of the transfer.
This allows transfers to be completed by independent relayers to improve UX for users that will only need to send a
single transaction for as long as the fee is sufficient and the token accepted by anyone acting as relayer.

In order to keep `Transfer` messages small, they don't carry all metadata of the token. However, this means that before
a token can be transferred to a new chain for the first time, the metadata needs to be bridged and the wrapped asset
created. Metadata in this case includes the amount of decimals which is a core requirement for instantiating a token.

The metadata of a token can be attested by calling `attestToken` on its respective native chain which will produce a
`AssetMeta` wormhole message. This message can be used to attest state and initialize a WrappedAsset on any chain in the
wormhole network using the details. A token is identified by the tuple `(chain_id, chain_address)` and metadata should
be mapped to this identifier. A wrapped asset may only ever be created once for a given identifier and not updated.

### API / database schema

Proposed bridge interface:

`attestToken(address token)` - Produce a `AssetMeta` message for a given token

`transfer(address token, uint256 amount, uint16 recipient_chain, bytes32 recipient, uint256 fee)` - Initiate
a `Transfer`

`createWrapped(Message assetMeta)` - Creates a wrapped asset using `AssetMeta`

`completeTransfer(Message transfer)` - Execute a `Transfer` message

`registerChain(Message registerChain)` - Execute a `RegisterChain` governance message

`upgrade(Message upgrade)` - Execute a `UpgradeContract` governance message

---
**Payloads**:

Transfer:

```
PayloadID uint8 = 1
// Amount being transferred (big-endian uint256)
Amount [32]uint8
// Address of the token. Left-zero-padded if shorter than 32 bytes
TokenAddress [32]uint8
// Chain ID of the token
TokenChain uint16
// Address of the recipient. Left-zero-padded if shorter than 32 bytes
To [32]uint8
// Chain ID of the recipient
ToChain uint16
// Amount of tokens (big-endian uint256) that the user is willing to pay as relayer fee. Must be <= Amount.
Fee [32]uint8
```

AssetMeta:

```
PayloadID uint8 = 2
// Address of the token. Left-zero-padded if shorter than 32 bytes
TokenAddress [32]uint8
// Chain ID of the token
TokenChain uint16
// Number of decimals of the token
Decimals uint8
// Symbol of the token (UTF-8)
Symbol [32]uint8
// Name of the token (UTF-8)
Name [32]uint8
```

RegisterChain:

```
PayloadID uint8 = 3
// Chain ID
ChainID uint16
// Emitter address. Left-zero-padded if shorter than 32 bytes
EmitterAddress [32]uint8
```

UpgradeContract:

```
PayloadID uint8 = 4
// Address of the new contract
NewContract [32]uint8
```

## Caveats

There is no guarantee for completion of transfers. If a user initiates a transfer and doesn't call `completeTransfer`
on the target chain, a transfer might not be completed. In case a guardian set change happens in-between and the
original signer guardian set expires, the transfer will be stuck indefinitely.

Since there is no way for a token bridge endpoint to know which other chain already has wrapped assets set up for the
native asset on its chain, there may be transfers initiated for assets that don't have wrapped assets set up yet on the
target chain. However, the transfer will become executable once the wrapped asset is set up (which can be done any time)
.