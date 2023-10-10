package core

import (
	"bytes"

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

type Module struct {
	executive      executive.Module
	decoder        types.RuntimeDecoder
	runtimeVersion *primitives.RuntimeVersion
	memUtils       utils.WasmMemoryTranslator
}

func New(module executive.Module, decoder types.RuntimeDecoder, runtimeVersion *primitives.RuntimeVersion) Module {
	return Module{
		executive:      module,
		decoder:        decoder,
		runtimeVersion: runtimeVersion,
		memUtils:       utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// Version returns a pointer-size SCALE-encoded Runtime version.
// [Specification](https://spec.polkadot.network/#defn-rt-core-version)
func (m Module) Version() int64 {
	buffer := &bytes.Buffer{}
	m.runtimeVersion.Encode(buffer)
	return m.memUtils.BytesToOffsetAndSize(buffer.Bytes())
}

// InitializeBlock starts the execution of a particular block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded header of the block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-initialize-block)
func (m Module) InitializeBlock(dataPtr int32, dataLen int32) {
	data := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)
	header := primitives.DecodeHeader(buffer)
	m.executive.InitializeBlock(header)
}

// ExecuteBlock executes the provided block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded block.
// [Specification](https://spec.polkadot.network/#sect-rte-core-execute-block)
func (m Module) ExecuteBlock(dataPtr int32, dataLen int32) {
	data := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)
	block := m.decoder.DecodeBlock(buffer)
	m.executive.ExecuteBlock(block)
}
