/*
BlockBuilder - Version 4.
*/
package blockbuilder

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/inherent"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type BlockBuilder interface {
	ApplyExtrinsic(dataPtr int32, dataLen int32) int64
	FinalizeBlock(dataPtr int32, dataLen int32) int64
	InherentExtrinisics(dataPtr int32, dataLen int32) int64
	CheckInherents(dataPtr int32, dataLen int32) int64
	RandomSeed(dataPtr int32, dataLen int32) int64
}

/*
https://spec.polkadot.network/#sect-rte-apply-extrinsic

SCALE encoded arguments (extrinsic types.Extrinsic) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func ApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	buffer := &bytes.Buffer{}
	buffer.Write(utils.ToWasmMemorySlice(dataPtr, dataLen))
	uxt := types.DecodeUncheckedExtrinsic(buffer)

	ok, err := executive.ApplyExtrinsic(uxt)
	var applyExtrinsicResult types.ApplyExtrinsicResult
	if err != nil {
		applyExtrinsicResult = types.NewApplyExtrinsicResult(err)
	} else {
		applyExtrinsicResult = types.NewApplyExtrinsicResult(ok)
	}

	buffer.Reset()
	applyExtrinsicResult.Encode(buffer)

	return utils.BytesToOffsetAndSize(buffer.Bytes())
}

/*
https://spec.polkadot.network/#defn-rt-blockbuilder-finalize-block

SCALE encoded arguments () allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded (types.Header) data.
*/

// FinalizeBlock finalizes block - it is up the caller to ensure that all header fields are valid
// except state-root.
func FinalizeBlock() int64 {
	system.NoteFinishedExtrinsics()

	blockNumber := system.StorageGetBlockNumber()

	executive.IdleAndFinalizeHook(blockNumber)

	header := system.Finalize()
	encodedHeader := header.Bytes()

	return utils.BytesToOffsetAndSize(encodedHeader)
}

/*
https://spec.polkadot.network/#defn-rt-builder-inherent-extrinsics

SCALE encoded arguments (data types.InherentsData) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]types.Extrinsic) data.
*/
func InherentExtrinisics(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)

	buffer := &bytes.Buffer{}
	buffer.Write(b)

	inherentData, err := types.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	result := timestamp.CreateInherent(*inherentData)
	result = append(sc.ToCompact(1).Bytes(), result...)

	return utils.BytesToOffsetAndSize(result)
}

/*
https://spec.polkadot.network/#id-blockbuilder_check_inherents

SCALE encoded arguments (block types.Block, data types.InherentsData) allocated in the Wasm VM memory, passed as:

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func CheckInherents(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)

	buffer := &bytes.Buffer{}
	buffer.Write(b)

	block := types.DecodeBlock(buffer)

	inherentData, err := types.DecodeInherentData(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	buffer.Reset()

	checkInherentsResult := inherent.CheckInherents(*inherentData, block)

	checkInherentsResult.Encode(buffer)
	return utils.BytesToOffsetAndSize(buffer.Bytes())
}
