package genesisbuilder

// todo

/// Get the default `GenesisConfig` as a JSON blob. For more info refer to
/// [`sp_genesis_builder::GenesisBuilder::create_default_config`]
///
// pub fn create_default_config<GC>() -> sp_std::vec::Vec<u8>
// where
// 	GC: BuildGenesisConfig + Default,
// {
// 	serde_json::to_string(&GC::default())
// 		.expect("serialization to json is expected to work. qed.")
// 		.into_bytes()
// }

/// Build `GenesisConfig` from a JSON blob not using any defaults and store it in the storage. For
/// more info refer to [`sp_genesis_builder::GenesisBuilder::build_config`].
///
// pub fn build_config<GC: BuildGenesisConfig>(json: sp_std::vec::Vec<u8>) -> BuildResult {
// 	let gc = serde_json::from_slice::<GC>(&json)
// 		.map_err(|e| format_runtime_string!("Invalid JSON blob: {}", e))?;
// 	<GC as BuildGenesisConfig>::build(&gc);
// 	Ok(())
// }
