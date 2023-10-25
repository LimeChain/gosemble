package blockbuilder

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "BlockBuilder"
	apiVersion    = 6
)

type BlockBuilder interface {
	ApplyExtrinsic(dataPtr int32, dataLen int32) int64
	FinalizeBlock() int64
	InherentExtrinsics(dataPtr int32, dataLen int32) int64
	CheckInherents(dataPtr int32, dataLen int32) int64
}

type Module struct {
	runtimeExtrinsic extrinsic.RuntimeExtrinsic
	executive        executive.Module
	decoder          types.RuntimeDecoder
	memUtils         utils.WasmMemoryTranslator
}

func New(runtimeExtrinsic extrinsic.RuntimeExtrinsic, executive executive.Module, decoder types.RuntimeDecoder) Module {
	return Module{
		runtimeExtrinsic: runtimeExtrinsic,
		executive:        executive,
		decoder:          decoder,
		memUtils:         utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// ApplyExtrinsic applies an extrinsic to a particular block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded unchecked extrinsic.
// Returns a pointer-size of the SCALE-encoded result, which specifies if this extrinsic is included in this block or not.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-apply-extrinsic)
func (m Module) ApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	uxt, err := m.decoder.DecodeUncheckedExtrinsic(buffer)
	if err != nil {
		log.Critical(err.Error())
		return 0
	}

	ok, errApplyExtr := m.executive.ApplyExtrinsic(uxt)
	var applyExtrinsicResult primitives.ApplyExtrinsicResult
	if errApplyExtr != nil {
		applyExtrinsicResult = primitives.NewApplyExtrinsicResult(errApplyExtr)
	} else {
		applyExtrinsicResult = primitives.NewApplyExtrinsicResult(ok)
	}

	buffer.Reset()
	applyExtrinsicResult.Encode(buffer)

	return m.memUtils.BytesToOffsetAndSize(buffer.Bytes())
}

// FinalizeBlock finalizes the state changes for the current block.
// Returns a pointer-size of the SCALE-encoded header for this block.
// [Specification](https://spec.polkadot.network/#defn-rt-blockbuilder-finalize-block)
func (m Module) FinalizeBlock() int64 {
	header, err := m.executive.FinalizeBlock()
	if err != nil {
		log.Critical(err.Error())
		return 0
	}
	encodedHeader := header.Bytes()
	return m.memUtils.BytesToOffsetAndSize(encodedHeader)
}

// InherentExtrinsics generates inherent extrinsics. Inherent data varies depending on chain configuration.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded inherent data.
// Returns a pointer-size of the SCALE-encoded timestamp extrinsic.
// [Specification](https://spec.polkadot.network/#defn-rt-builder-inherent-extrinsics)
func (m Module) InherentExtrinsics(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	inherentData, err := primitives.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
		return 0
	}

	result, err := m.runtimeExtrinsic.CreateInherents(*inherentData)
	if err != nil {
		log.Critical(err.Error())
		return 0
	}

	return m.memUtils.BytesToOffsetAndSize(result)
}

// CheckInherents checks the inherents are valid.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded inherent data.
// Returns a pointer-size of the SCALE-encoded result, specifying if all inherents are valid.
// [Specification](https://spec.polkadot.network/#id-blockbuilder_check_inherents)
func (m Module) CheckInherents(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	block, err := m.decoder.DecodeBlock(buffer)
	if err != nil {
		log.Critical(err.Error())
		return 0
	}

	inherentData, err := primitives.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	result := m.runtimeExtrinsic.CheckInherents(*inherentData, block)

	return m.memUtils.BytesToOffsetAndSize(result.Bytes())
}
