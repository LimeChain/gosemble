package core

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
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
	encoded := m.runtimeVersion.Bytes()

	return m.memUtils.BytesToOffsetAndSize(encoded)
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
	header, err := primitives.DecodeHeader(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	err = m.executive.InitializeBlock(header)
	if err != nil {
		log.Critical(err.Error())
	}
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
	block, err := m.decoder.DecodeBlock(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	err = m.executive.ExecuteBlock(block)
	if err != nil {
		log.Critical(err.Error())
	}
}

func (m Module) Metadata() primitives.RuntimeApiMetadata {
	methods := sc.Sequence[primitives.RuntimeApiMethodMetadata]{
		primitives.RuntimeApiMethodMetadata{
			Name:   "version",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(metadata.TypesRuntimeVersion),
			Docs:   sc.Sequence[sc.Str]{" Returns the version of the runtime."},
		},
		primitives.RuntimeApiMethodMetadata{
			Name: "execute_block",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
				primitives.RuntimeApiMethodParamMetadata{
					Name: "block",
					Type: sc.ToCompact(metadata.TypesBlock),
				},
			},
			Output: sc.ToCompact(metadata.TypesEmptyTuple),
			Docs:   sc.Sequence[sc.Str]{" Execute the given block."},
		},
		primitives.RuntimeApiMethodMetadata{
			Name: "initialize_block",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
				primitives.RuntimeApiMethodParamMetadata{
					Name: "header",
					Type: sc.ToCompact(metadata.Header),
				},
			},
			Output: sc.ToCompact(metadata.TypesEmptyTuple),
			Docs:   sc.Sequence[sc.Str]{" Initialize a block with the given header."},
		},
	}

	return primitives.RuntimeApiMetadata{
		Name:    ApiModuleName,
		Methods: methods,
		Docs:    sc.Sequence[sc.Str]{" The `Core` runtime api that every Substrate runtime needs to implement."},
	}
}
