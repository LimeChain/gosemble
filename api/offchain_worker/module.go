package offchain_worker

import (
	"bytes"

	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "OffchainWorkerApi"
	apiVersion    = 2
)

type Module struct {
	executive executive.Module
	memUtils  utils.WasmMemoryTranslator
}

func New(executive executive.Module) Module {
	return Module{
		executive: executive,
		memUtils:  utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() types.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return types.NewApiItem(hash, apiVersion)
}

// OffchainWorker starts an off-chain task for an imported block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-offchainworkerapi_offchain_worker)
func (m Module) OffchainWorker(dataPtr int32, dataLen int32) {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	header, err := primitives.DecodeHeader(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	m.executive.OffchainWorker(header)
}
