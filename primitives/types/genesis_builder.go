package types

// todo

//! Substrate genesis config builder
//!
//! This Runtime API allows to construct `GenesisConfig`, in particular:
//! - serialize the runtime default `GenesisConfig` struct into json format,
//! - put the GenesisConfig struct into the storage. Internally this operation calls
//!   `GenesisBuild::build` function for all runtime pallets, which is typically provided by
//!   pallet's author.
//! - deserialize the `GenesisConfig` from given json blob and put `GenesisConfig` into the state
//!   storage. Allows to build customized configuration.
//!
//! Providing externalities with empty storage and putting `GenesisConfig` into storage allows to
//! catch and build the raw storage of `GenesisConfig` which is the foundation for genesis block.
//
/// The result type alias, used in build methods. `Err` contains formatted error message.
///
// pub type Result = core::result::Result<(), sp_runtime::RuntimeString>;

// #[pallet::genesis_config]
// 	#[derive(DefaultNoBound)]
// 	pub struct GenesisConfig<T: Config<I>, I: 'static = ()> {
// 		/// Initial pallet operating mode.
// 		pub operating_mode: MessagesOperatingMode,
// 		/// Initial pallet owner.
// 		pub owner: Option<T::AccountId>,
// 		/// Dummy marker.
// 		pub phantom: sp_std::marker::PhantomData<I>,
// 	}

// todo Hooks and will be implemented by all frames (see other hooks)
//
// /// A trait to define the build function of a genesis config for both runtime and pallets.
// ///
// /// Replaces deprecated [`GenesisBuild<T,I>`].
// pub trait BuildGenesisConfig: sp_runtime::traits::MaybeSerializeDeserialize {
// 	/// The build function puts initial `GenesisConfig` keys/values pairs into the storage.
// 	fn build(&self);
// }
