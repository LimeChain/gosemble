package offchain_worker

import (
	"bytes"

	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Module struct {
	executive executive.Module
}

func New(executive executive.Module) Module {
	return Module{executive: executive}
}

// OffchainWorker starts an off-chain task for an imported block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-offchainworkerapi_offchain_worker)
func (m Module) OffchainWorker(dataPtr int32, dataLen int32) {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	header := types.DecodeHeader(buffer)

	m.executive.OffchainWorker(header)
}
