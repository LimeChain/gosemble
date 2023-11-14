package transaction_payment

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/transaction_payment/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	mockUxt                *mocks.UncheckedExtrinsic
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
	expect := primitives.NewApiItem(hexName, apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_QueryInfo_Signed(t *testing.T) {
	target := setup()

	partialFee := sc.NewU128(10)
	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfoWeight,
		Class:      dispatchInfoClass,
		PartialFee: partialFee,
	}

	bufferUxt := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(mockUxt, nil)
	mockUxt.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockUxt.On("IsSigned").Return(true)
	mockTransactionPayment.On("ComputeFee", length, dispatchInfo, constants.DefaultTip).Return(partialFee, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", runtimeDispatchInfo.Bytes()).Return(ptrAndSize)

	result := target.QueryInfo(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockUxt.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockUxt.AssertCalled(t, "IsSigned")
	mockTransactionPayment.AssertCalled(t, "ComputeFee", length, dispatchInfo, constants.DefaultTip)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", runtimeDispatchInfo.Bytes())
}

func Test_Module_QueryInfo_Unsigned(t *testing.T) {
	target := setup()

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfoWeight,
		Class:      dispatchInfoClass,
		PartialFee: constants.Zero,
	}

	bufferUxt := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(mockUxt, nil)
	mockUxt.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockUxt.On("IsSigned").Return(false)
	mockMemoryUtils.On("BytesToOffsetAndSize", runtimeDispatchInfo.Bytes()).Return(ptrAndSize)

	result := target.QueryInfo(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockUxt.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockUxt.AssertCalled(t, "IsSigned")
	mockTransactionPayment.AssertNotCalled(t, "ComputeFee", mock.Anything, mock.Anything, mock.Anything)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", runtimeDispatchInfo.Bytes())
}

func Test_Module_QueryFeeDetails_Signed(t *testing.T) {
	target := setup()

	feeDetails := types.FeeDetails{
		InclusionFee: sc.NewOption[types.InclusionFee](
			types.NewInclusionFee(
				sc.NewU128(9),
				sc.NewU128(8),
				sc.NewU128(7),
			)),
		Tip: constants.DefaultTip,
	}
	bufferUxt := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(mockUxt, nil)
	mockUxt.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockUxt.On("IsSigned").Return(true)
	mockTransactionPayment.On("ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip).Return(feeDetails, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", feeDetails.Bytes()).Return(ptrAndSize)

	result := target.QueryFeeDetails(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockUxt.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockUxt.AssertCalled(t, "IsSigned")
	mockTransactionPayment.AssertCalled(t, "ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", feeDetails.Bytes())
}

func Test_Module_QueryFeeDetails_Unsigned(t *testing.T) {
	target := setup()

	feeDetails := types.FeeDetails{
		InclusionFee: sc.NewOption[types.InclusionFee](nil),
	}
	bufferUxt := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeUncheckedExtrinsic", bufferUxt).Return(mockUxt, nil)
	mockUxt.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockUxt.On("IsSigned").Return(false)
	mockMemoryUtils.On("BytesToOffsetAndSize", feeDetails.Bytes()).Return(ptrAndSize)

	result := target.QueryFeeDetails(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockUxt.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockUxt.AssertCalled(t, "IsSigned")
	mockTransactionPayment.AssertNotCalled(t, "ComputeFeeDetails", mock.Anything, mock.Anything, mock.Anything)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", feeDetails.Bytes())
}

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	expect := primitives.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[primitives.RuntimeApiMethodMetadata]{
			primitives.RuntimeApiMethodMetadata{
				Name: "query_info",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "uxt",
						Type: sc.ToCompact(metadata.UncheckedExtrinsic),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "len",
						Type: sc.ToCompact(metadata.PrimitiveTypesU32),
					},
				},
				Output: sc.ToCompact(metadata.TypesTransactionPaymentRuntimeDispatchInfo),
				Docs:   sc.Sequence[sc.Str]{},
			},
			primitives.RuntimeApiMethodMetadata{
				Name: "query_fee_details",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "uxt",
						Type: sc.ToCompact(metadata.UncheckedExtrinsic),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "len",
						Type: sc.ToCompact(metadata.PrimitiveTypesU32),
					},
				},
				Output: sc.ToCompact(metadata.TypesTransactionPaymentFeeDetails),
				Docs:   sc.Sequence[sc.Str]{},
			},
		},
		Docs: sc.Sequence[sc.Str]{},
	}

	assert.Equal(t, expect, target.Metadata())
}

func setup() Module {
	mockTransactionPayment = new(mocks.TransactionPaymentModule)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockMemoryUtils = new(mocks.MemoryTranslator)
	mockUxt = new(mocks.UncheckedExtrinsic)
	mockCall = new(mocks.Call)

	target := New(mockRuntimeDecoder, mockTransactionPayment)
	target.memUtils = mockMemoryUtils

	return target
}
