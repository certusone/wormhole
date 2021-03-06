use borsh::{BorshDeserialize, BorshSerialize};
use solana_program::pubkey::Pubkey;
use solitaire::{
    pack_type,
    processors::seeded::{AccountOwner, Owned},
};
use spl_token::state::{Account, Mint};

pub type Address = [u8; 32];
pub type ChainID = u16;

#[derive(Default, Clone, Copy, BorshDeserialize, BorshSerialize)]
pub struct Config {
    pub wormhole_bridge: Pubkey,
    pub fees: FeeStructure,
}

impl Owned for Config {
    fn owner(&self) -> AccountOwner {
        AccountOwner::This
    }
}

#[derive(Default, Clone, Copy, BorshDeserialize, BorshSerialize)]
pub struct FeeStructure {
    pub usd_ephemeral: u64,
    pub usd_persistent: u64,
}

impl Owned for FeeStructure {
    fn owner(&self) -> AccountOwner {
        AccountOwner::This
    }
}

#[derive(Default, Clone, Copy, BorshDeserialize, BorshSerialize)]
pub struct EndpointRegistration {
    pub chain: ChainID,
    pub contract: Address,
}

impl Owned for EndpointRegistration {
    fn owner(&self) -> AccountOwner {
        AccountOwner::This
    }
}

#[derive(Default, Clone, Copy, BorshDeserialize, BorshSerialize)]
pub struct WrappedMeta {
    pub chain: ChainID,
    pub token_address: Address,
}

impl Owned for WrappedMeta {
    fn owner(&self) -> AccountOwner {
        AccountOwner::This
    }
}

pack_type!(SplMint, Mint, AccountOwner::Other(spl_token::id()));
pack_type!(SplAccount, Account, AccountOwner::Other(spl_token::id()));
