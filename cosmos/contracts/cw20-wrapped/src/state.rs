use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cosmwasm_std::{CanonicalAddr, ReadonlyStorage, Storage};
use cosmwasm_storage::{singleton, singleton_read, ReadonlySingleton, Singleton};

pub const KEY_WRAPPED_ASSET: &[u8] = b"wrappedAsset";

// Created at initialization and reference original asset and bridge address
#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct WrappedAssetInfo {
    pub asset_chain: u8,            // Asset chain id
    pub asset_address: Vec<u8>,     // Asset smart contract address on the original chain
    pub bridge: CanonicalAddr,      // Bridge address, authorized to mint and burn wrapped tokens
}

pub fn wrapped_asset_info<S: Storage>(storage: &mut S) -> Singleton<S, WrappedAssetInfo> {
    singleton(storage, KEY_WRAPPED_ASSET)
}

pub fn wrapped_asset_info_read<S: ReadonlyStorage>(
    storage: &S,
) -> ReadonlySingleton<S, WrappedAssetInfo> {
    singleton_read(storage, KEY_WRAPPED_ASSET)
}
