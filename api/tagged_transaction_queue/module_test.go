package tagged_transaction_queue

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
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

	validTx = primitives.DefaultValidTransaction()
	// txValidityError, _       = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	// todo
	txValidityError           = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	txValidityUnexpectedError = primitives.NewTransactionValidityError(primitives.NewUnexpectedError(fmt.Errorf("panic error")))
	validitySuccessResult, _  = primitives.NewTransactionValidityResult(validTx)
	validityFailResult, _     = primitives.NewTransactionValidityResult(txValidityError)
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

func Test_Module_ValidateTransaction_UnexpectedError_Panics(t *testing.T) {
	target := setup()

	data := append(txSource.Bytes(), blockHash.Bytes()...)
	expectBuffer := bytes.NewBuffer(data)
	_, err := expectBuffer.ReadByte()
	assert.Nil(t, err)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(data)
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", expectBuffer).Return(mockUxt, nil)
	mockExecutive.On("ValidateTransaction", txSource, mockUxt, blockHash).
		Return(primitives.ValidTransaction{}, txValidityUnexpectedError)
	mockMemoryUtils.On("BytesToOffsetAndSize", validityFailResult.Bytes()).Return(ptrAndSize)

	assert.PanicsWithValue(t, "panic error", func() { target.ValidateTransaction(dataPtr, dataLen) })
}

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	expect := primitives.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[primitives.RuntimeApiMethodMetadata]{
			primitives.RuntimeApiMethodMetadata{
				Name: "validate_transaction",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "source",
						Type: sc.ToCompact(metadata.TypesTransactionSource),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "tx",
						Type: sc.ToCompact(metadata.UncheckedExtrinsic),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "block_hash",
						Type: sc.ToCompact(metadata.TypesH256),
					},
				},
				Output: sc.ToCompact(metadata.TypesResultValidityTransaction),
				Docs: sc.Sequence[sc.Str]{" Validate the transaction.",
					"",
					" This method is invoked by the transaction pool to learn details about given transaction.",
					" The implementation should make sure to verify the correctness of the transaction",
					" against current state. The given `block_hash` corresponds to the hash of the block",
					" that is used as current state.",
					"",
					" Note that this call may be performed by the pool multiple times and transactions",
					" might be verified in any possible order."},
			},
		},
		Docs: sc.Sequence[sc.Str]{" The `TaggedTransactionQueue` api trait for interfering with the transaction queue."},
	}

	assert.Equal(t, expect, target.Metadata())
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
