package transaction_payment_call

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(2)

	length = sc.U32(5)

	baseWeight         = primitives.WeightFromParts(1, 2)
	dispatchInfoWeight = primitives.WeightFromParts(2, 3)
	dispatchInfoClass  = primitives.NewDispatchClassNormal()
	dispatchInfoPays   = primitives.NewPaysYes()

	dispatchInfo = primitives.DispatchInfo{
		Weight:  dispatchInfoWeight,
		Class:   dispatchInfoClass,
		PaysFee: dispatchInfoPays,
	}
)

var (
	mockTransactionPayment *mocks.TransactionPaymentModule
	mockRuntimeDecoder     *mocks.RuntimeDecoder
	mockMemoryUtils        *mocks.MemoryTranslator
	mockCall               *mocks.Call
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

func Test_Module_QueryCallInfo(t *testing.T) {
	target := setup()

	partialFee := sc.NewU128(10)
	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfoWeight,
		Class:      dispatchInfoClass,
		PartialFee: partialFee,
	}
	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFee", length, dispatchInfo, constants.DefaultTip).Return(partialFee)
	mockMemoryUtils.On("BytesToOffsetAndSize", runtimeDispatchInfo.Bytes()).Return(ptrAndSize)

	result := target.QueryCallInfo(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockTransactionPayment.AssertCalled(t, "ComputeFee", length, dispatchInfo, constants.DefaultTip)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", runtimeDispatchInfo.Bytes())
}

func Test_Module_QueryCallFeeDetails(t *testing.T) {
	target := setup()

	feeDetails := primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](
			primitives.NewInclusionFee(
				sc.NewU128(9),
				sc.NewU128(8),
				sc.NewU128(7),
			)),
		Tip: constants.DefaultTip,
	}
	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip).Return(feeDetails)
	mockMemoryUtils.On("BytesToOffsetAndSize", feeDetails.Bytes()).Return(ptrAndSize)

	result := target.QueryCallFeeDetails(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockTransactionPayment.AssertCalled(t, "ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", feeDetails.Bytes())
}

func setup() Module {
	mockTransactionPayment = new(mocks.TransactionPaymentModule)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockMemoryUtils = new(mocks.MemoryTranslator)
	mockCall = new(mocks.Call)

	target := New(mockRuntimeDecoder, mockTransactionPayment)
	target.memUtils = mockMemoryUtils

	return target
}