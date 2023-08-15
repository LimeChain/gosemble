package offchain_worker

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "OffchainWorkerApi"
	apiVersion    = 2
)

type Module[N sc.Numeric] struct {
	executive executive.Module[N]
}

func New[N sc.Numeric](executive executive.Module[N]) Module[N] {
	return Module[N]{executive: executive}
}

func (m Module[N]) Name() string {
	return ApiModuleName
}

func (m Module[N]) Item() types.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return types.NewApiItem(hash, apiVersion)
}

// OffchainWorker starts an off-chain task for an imported block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/chap-runtime-api#id-offchainworkerapi_offchain_worker)
func (m Module[N]) OffchainWorker(dataPtr int32, dataLen int32) {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)
	header := primitives.DecodeHeader[N](buffer)
	m.executive.OffchainWorker(header)
}
