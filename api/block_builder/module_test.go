package blockbuilder

import (
	"bytes"
	"io"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(5)

	bUxt = []byte("unchecked_extrinsic")
	uxt  = new(mocks.UncheckedExtrinsic)
)

var (
	mockExecutive        *mocks.Executive
	mockRuntimeDecoder   *mocks.RuntimeDecoder
	mockRuntimeExtrinsic *mocks.RuntimeExtrinsic
	mockMemoryUtils      *mocks.MemoryTranslator
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName[:], apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_ApplyExtrinsic_Success(t *testing.T) {
	target := setup()

	bufferUxt := bytes.NewBuffer(bUxt)
	outcome := primitives.NewDispatchOutcome(sc.Empty{})
	bExtrinsicResult := primitives.NewApplyExtrinsicResult(outcome).Bytes()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bUxt)
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(uxt)
	mockExecutive.On("ApplyExtrinsic", uxt).Return(outcome, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", bExtrinsicResult).Return(ptrAndSize)

	result := target.ApplyExtrinsic(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ApplyExtrinsic", uxt)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", bExtrinsicResult)
}

func Test_Module_ApplyExtrinsic_Fails(t *testing.T) {
	target := setup()

	bufferUxt := bytes.NewBuffer(bUxt)
	outcome := primitives.NewDispatchOutcome(sc.Empty{})
	validityError := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	bExtrinsicResult := primitives.NewApplyExtrinsicResult(validityError).Bytes()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bUxt)
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(uxt)
	mockExecutive.On("ApplyExtrinsic", uxt).Return(outcome, validityError)
	mockMemoryUtils.On("BytesToOffsetAndSize", bExtrinsicResult).Return(ptrAndSize)

	result := target.ApplyExtrinsic(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ApplyExtrinsic", uxt)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", bExtrinsicResult)
}

func Test_Module_FinalizeBlock(t *testing.T) {
	target := setup()

	parentHash := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355c").ToBytes()
	stateRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b").ToBytes()
	extrinsicsRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a").ToBytes()
	header := primitives.Header{
		ParentHash: primitives.Blake2bHash{
			FixedSequence: sc.BytesToFixedSequenceU8(parentHash)},
		Number:         5,
		StateRoot:      primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(stateRoot)},
		ExtrinsicsRoot: primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(extrinsicsRoot)},
	}

	mockExecutive.On("FinalizeBlock").Return(header)
	mockMemoryUtils.On("BytesToOffsetAndSize", header.Bytes()).Return(ptrAndSize)

	result := target.FinalizeBlock()

	assert.Equal(t, ptrAndSize, result)
	mockExecutive.AssertCalled(t, "FinalizeBlock")
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", header.Bytes())
}

func Test_Module_InherentExtrinsics_Success(t *testing.T) {
	target := setup()

	bInherentData := sc.ToCompact(0).Bytes()
	inherentData := primitives.NewInherentData()
	bCreate := []byte{1}

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bInherentData)
	mockRuntimeExtrinsic.On("CreateInherents", *inherentData).Return(bCreate)
	mockMemoryUtils.On("BytesToOffsetAndSize", bCreate).Return(ptrAndSize)

	result := target.InherentExtrinsics(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeExtrinsic.AssertCalled(t, "CreateInherents", *inherentData)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", bCreate)
}

func Test_Module_InherentExtrinsics_InvalidInherentData(t *testing.T) {
	target := setup()

	bInvalidInherentData := sc.ToCompact(1).Bytes()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bInvalidInherentData)

	assert.PanicsWithValue(t,
		io.EOF.Error(),
		func() { target.InherentExtrinsics(dataPtr, dataLen) },
	)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeExtrinsic.AssertNotCalled(t, "CreateInherents", mock.Anything)
	mockMemoryUtils.AssertNotCalled(t, "BytesToOffsetAndSize", mock.Anything)
}

func Test_Module_CheckInherents_Success(t *testing.T) {
	target := setup()

	block := types.NewBlock(primitives.Header{Number: 1}, sc.Sequence[primitives.UncheckedExtrinsic]{})
	bInherentData := sc.ToCompact(0).Bytes()
	bufferData := bytes.NewBuffer(bInherentData)
	inherentData := primitives.NewInherentData()
	checkResult := primitives.NewCheckInherentsResult()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bInherentData)
	mockRuntimeDecoder.On("DecodeBlock", bufferData).Return(block)
	mockRuntimeExtrinsic.On("CheckInherents", *inherentData, block).Return(checkResult)
	mockMemoryUtils.On("BytesToOffsetAndSize", checkResult.Bytes()).Return(ptrAndSize)

	result := target.CheckInherents(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockRuntimeExtrinsic.AssertCalled(t, "CheckInherents", *inherentData, block)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", checkResult.Bytes())
}

func Test_Module_CheckInherents_InvalidInherentData(t *testing.T) {
	target := setup()

	block := types.NewBlock(primitives.Header{Number: 1}, sc.Sequence[primitives.UncheckedExtrinsic]{})
	bytesInvalidInherentData := sc.ToCompact(1).Bytes()
	bufferData := bytes.NewBuffer(bytesInvalidInherentData)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bytesInvalidInherentData)
	mockRuntimeDecoder.On("DecodeBlock", bufferData).Return(block)

	assert.PanicsWithValue(t,
		io.EOF.Error(),
		func() { target.CheckInherents(dataPtr, dataLen) },
	)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockRuntimeExtrinsic.AssertNotCalled(t, "CheckInherents", mock.Anything)
	mockMemoryUtils.AssertNotCalled(t, "BytesToOffsetAndSize", mock.Anything)
}

func setup() Module {
	mockExecutive = new(mocks.Executive)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockRuntimeExtrinsic = new(mocks.RuntimeExtrinsic)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New(mockRuntimeExtrinsic, mockExecutive, mockRuntimeDecoder)
	target.memUtils = mockMemoryUtils

	return target
}
