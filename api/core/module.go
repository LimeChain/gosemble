package core

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "Core"
	apiVersion    = 4
)

type Core interface {
	Version() int64
	ExecuteBlock(dataPtr int32, dataLen int32)
	InitializeBlock(dataPtr int32, dataLen int32)
}

type Module[N sc.Numeric] struct {
	executive      executive.Module[N]
	decoder        types.ModuleDecoder[N]
	runtimeVersion *primitives.RuntimeVersion
}

func New[N sc.Numeric](module executive.Module[N], decoder types.ModuleDecoder[N], runtimeVersion *primitives.RuntimeVersion) Module[N] {
	return Module[N]{
		module,
		decoder,
		runtimeVersion,
	}
}

func (m Module[N]) Name() string {
	return ApiModuleName
}

func (m Module[N]) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// Version returns a pointer-size SCALE-encoded Runtime version.
// [Specification](https://spec.polkadot.network/#defn-rt-core-version)
func (m Module[N]) Version() int64 {
	buffer := &bytes.Buffer{}
	m.runtimeVersion.Encode(buffer)

	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

// InitializeBlock starts the execution of a particular block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-initialize-block)
func (m Module[N]) InitializeBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)
	header := primitives.DecodeHeader[N](buffer)
	m.executive.InitializeBlock(header)
}

// ExecuteBlock executes the provided block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-execute-block)
func (m Module[N]) ExecuteBlock(dataPtr int32, dataLen int32) {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)
	block := m.decoder.DecodeBlock(buffer)
	m.executive.ExecuteBlock(block)
}
