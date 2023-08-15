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
	apiModuleName = "BlockBuilder"
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
	decoder          types.ModuleDecoder
}

func New(runtimeExtrinsic extrinsic.RuntimeExtrinsic, executive executive.Module, decoder types.ModuleDecoder) Module {
	return Module{
		runtimeExtrinsic,
		executive,
		decoder,
	}
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(apiModuleName))
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
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	uxt := m.decoder.DecodeUncheckedExtrinsic(buffer)

	ok, err := m.executive.ApplyExtrinsic(uxt)
	var applyExtrinsicResult primitives.ApplyExtrinsicResult
	if err != nil {
		applyExtrinsicResult = primitives.NewApplyExtrinsicResult(err)
	} else {
		applyExtrinsicResult = primitives.NewApplyExtrinsicResult(ok)
	}

	buffer.Reset()
	applyExtrinsicResult.Encode(buffer)

	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

// FinalizeBlock finalizes the state changes for the current block.
// Returns a pointer-size of the SCALE-encoded header for this block.
// [Specification](https://spec.polkadot.network/#defn-rt-blockbuilder-finalize-block)
func (m Module) FinalizeBlock() int64 {
	header := m.executive.FinalizeBlock()
	encodedHeader := header.Bytes()

	return utils.BytesToOffsetAndSize(encodedHeader)
}

// InherentExtrinsics generates inherent extrinsics. Inherent data varies depending on chain configuration.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded inherent data.
// Returns a pointer-size of the SCALE-encoded timestamp extrinsic.
// [Specification](https://spec.polkadot.network/#defn-rt-builder-inherent-extrinsics)
func (m Module) InherentExtrinsics(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	inherentData, err := primitives.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	result := m.runtimeExtrinsic.CreateInherents(*inherentData)

	return utils.BytesToOffsetAndSize(result)
}

// CheckInherents checks the inherents are valid.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded inherent data.
// Returns a pointer-size of the SCALE-encoded result, specifying if all inherents are valid.
// [Specification](https://spec.polkadot.network/#id-blockbuilder_check_inherents)
func (m Module) CheckInherents(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	block := m.decoder.DecodeBlock(buffer)

	inherentData, err := primitives.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	result := m.runtimeExtrinsic.CheckInherents(*inherentData, block)

	return utils.BytesToOffsetAndSize(result.Bytes())
}
