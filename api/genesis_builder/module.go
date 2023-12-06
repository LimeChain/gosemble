package genesisbuilder

import (
	"bytes"
	"encoding/json"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GenesisBuilder"
	apiVersion    = 1
)

type GenesisConfigMap = map[string]json.RawMessage

type Module struct {
	modules  []primitives.GenesisBuilder
	memUtils utils.WasmMemoryTranslator
}

func New(modules []primitives.GenesisBuilder) Module {
	return Module{
		modules:  modules,
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (m Module) CreateDefaultConfig() int64 {
	gcMap := make(GenesisConfigMap)
	for _, m := range m.modules {
		gcBz, err := m.CreateDefaultConfig()
		if err != nil {
			log.Critical(err.Error())
		}

		gcMap[m.ConfigModuleKey()] = gcBz
	}

	gcBz, err := json.Marshal(gcMap)
	if err != nil {
		log.Critical(err.Error())
	}

	return m.memUtils.BytesToOffsetAndSize(sc.BytesToSequenceU8(gcBz).Bytes())
}

func (m Module) BuildConfig(dataPtr int32, dataLen int32) int64 {
	gcBz := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	gcDecoded, err := sc.DecodeSequence[sc.U8](bytes.NewBuffer(gcBz))
	if err != nil {
		log.Critical(err.Error())
	}

	gcMap := make(GenesisConfigMap)
	if err := json.Unmarshal(sc.SequenceU8ToBytes(gcDecoded), &gcMap); err != nil {
		log.Critical(err.Error())
	}

	for _, module := range m.modules {
		gcBz, _ := gcMap[module.ConfigModuleKey()]
		// todo
		// if !ok {
		// 	continue
		// }

		if err := module.BuildConfig(gcBz); err != nil {
			log.Critical(err.Error())
		}
	}

	// todo should we return the error instead log.Critical()? double check substrate logic
	return m.memUtils.BytesToOffsetAndSize([]byte{0})
}

func (m Module) Metadata() primitives.RuntimeApiMetadata {
	// todo metadata
	return primitives.RuntimeApiMetadata{}
}

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

// sp_api::decl_runtime_apis! {
// 	/// API to interact with GenesisConfig for the runtime
// 	pub trait GenesisBuilder {
// 		/// Creates the default `GenesisConfig` and returns it as a JSON blob.
// 		///
// 		/// This function instantiates the default `GenesisConfig` struct for the runtime and serializes it into a JSON
// 		/// blob. It returns a `Vec<u8>` containing the JSON representation of the default `GenesisConfig`.
// 		fn create_default_config() -> sp_std::vec::Vec<u8>;
//
// 		/// Build `GenesisConfig` from a JSON blob not using any defaults and store it in the storage.
// 		///
// 		/// This function deserializes the full `GenesisConfig` from the given JSON blob and puts it into the storage.
// 		/// If the provided JSON blob is incorrect or incomplete or the deserialization fails, an error is returned.
// 		/// It is recommended to log any errors encountered during the process.
// 		///
// 		/// Please note that provided json blob must contain all `GenesisConfig` fields, no defaults will be used.
// 		fn build_config(json: sp_std::vec::Vec<u8>) -> Result;
// 	}
// }
