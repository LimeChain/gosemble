package tagged_transaction_queue

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(5)

	txSource  = primitives.NewTransactionSourceLocal()
	blockHash = primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}

	validTx                  = primitives.DefaultValidTransaction()
	txValidityError, _       = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	validitySuccessResult, _ = primitives.NewTransactionValidityResult(validTx)
	validityFailResult, _    = primitives.NewTransactionValidityResult(txValidityError)
)

var (
	mockExecutive      *mocks.Executive
	mockRuntimeDecoder *mocks.RuntimeDecoder
	mockMemoryUtils    *mocks.MemoryTranslator
	mockUxt            *mocks.UncheckedExtrinsic
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName, apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_ValidateTransaction_Success(t *testing.T) {
	target := setup()

	data := append(txSource.Bytes(), blockHash.Bytes()...)
	expectBuffer := bytes.NewBuffer(data)
	_, err := expectBuffer.ReadByte()
	assert.Nil(t, err)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(data)
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", expectBuffer).Return(mockUxt, nil)
	mockExecutive.On("ValidateTransaction", txSource, mockUxt, blockHash).Return(validTx, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", validitySuccessResult.Bytes()).Return(ptrAndSize)

	result := target.ValidateTransaction(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ValidateTransaction", txSource, mockUxt, blockHash)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", validitySuccessResult.Bytes())
}

func Test_Module_ValidateTransaction_Fails(t *testing.T) {
	target := setup()

	data := append(txSource.Bytes(), blockHash.Bytes()...)
	expectBuffer := bytes.NewBuffer(data)
	_, err := expectBuffer.ReadByte()
	assert.Nil(t, err)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(data)
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", expectBuffer).Return(mockUxt, nil)
	mockExecutive.On("ValidateTransaction", txSource, mockUxt, blockHash).
		Return(primitives.ValidTransaction{}, txValidityError)
	mockMemoryUtils.On("BytesToOffsetAndSize", validityFailResult.Bytes()).Return(ptrAndSize)

	result := target.ValidateTransaction(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ValidateTransaction", txSource, mockUxt, blockHash)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", validityFailResult.Bytes())
}

func setup() Module {
	mockExecutive = new(mocks.Executive)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockMemoryUtils = new(mocks.MemoryTranslator)
	mockUxt = new(mocks.UncheckedExtrinsic)

	target := New(mockExecutive, mockRuntimeDecoder)
	target.memUtils = mockMemoryUtils

	return target
}
