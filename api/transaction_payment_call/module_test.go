package transaction_payment_call

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/transaction_payment/types"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/log"
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
	dispatchInfoPays   = primitives.PaysYes

	dispatchInfo = primitives.DispatchInfo{
		Weight:  dispatchInfoWeight,
		Class:   dispatchInfoClass,
		PaysFee: dispatchInfoPays,
	}

	errPanic = errors.New("panic")
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
	expect := primitives.NewApiItem(hexName, apiVersion)

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
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFee", length, dispatchInfo, constants.DefaultTip).Return(partialFee, nil)
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

func Test_Module_QueryCallInfo_DecodeCall_Panics(t *testing.T) {
	target := setup()

	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.QueryCallInfo(dataPtr, dataLen) },
	)
}

func Test_Module_QueryCallInfo_DecodeU32_Panics(t *testing.T) {
	target := setup()

	bufferCall := bytes.NewBuffer([]byte{})

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return([]byte{})
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)

	assert.PanicsWithValue(t,
		io.EOF.Error(),
		func() { target.QueryCallInfo(dataPtr, dataLen) },
	)
}

func Test_Module_QueryCallInfo_ComputeFee_Panics(t *testing.T) {
	target := setup()

	partialFee := sc.NewU128(10)
	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFee", length, dispatchInfo, constants.DefaultTip).Return(partialFee, errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.QueryCallInfo(dataPtr, dataLen) },
	)
}

func Test_Module_QueryCallFeeDetails(t *testing.T) {
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
	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip).Return(feeDetails, nil)
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

func Test_Module_QueryCallFeeDetails_DecodeCall_Panics(t *testing.T) {
	target := setup()

	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.QueryCallFeeDetails(dataPtr, dataLen) },
	)
}

func Test_Module_QueryCallFeeDetails_DecodeU32_Panics(t *testing.T) {
	target := setup()

	bufferCall := bytes.NewBuffer([]byte{})

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return([]byte{})
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)

	assert.PanicsWithValue(t,
		io.EOF.Error(),
		func() { target.QueryCallFeeDetails(dataPtr, dataLen) },
	)
}

func Test_Module_QueryCallFeeDetails_ComputeFeeDetails_Panics(t *testing.T) {
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
	bufferCall := bytes.NewBuffer(length.Bytes())

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(length.Bytes())
	mockRuntimeDecoder.On("DecodeCall", bufferCall).Return(mockCall, nil)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfoWeight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfoClass)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfoPays)
	mockTransactionPayment.On("ComputeFeeDetails", length, dispatchInfo, constants.DefaultTip).Return(feeDetails, errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.QueryCallFeeDetails(dataPtr, dataLen) },
	)
}

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	expect := primitives.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[primitives.RuntimeApiMethodMetadata]{
			primitives.RuntimeApiMethodMetadata{
				Name: "query_call_info",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "call",
						Type: sc.ToCompact(metadata.RuntimeCall),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "len",
						Type: sc.ToCompact(metadata.PrimitiveTypesU32),
					},
				},
				Output: sc.ToCompact(metadata.TypesTransactionPaymentRuntimeDispatchInfo),
				Docs:   sc.Sequence[sc.Str]{" Query information of a dispatch class, weight, and fee of a given encoded `Call`."},
			},
			primitives.RuntimeApiMethodMetadata{
				Name: "query_call_fee_details",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "call",
						Type: sc.ToCompact(metadata.RuntimeCall),
					},
					primitives.RuntimeApiMethodParamMetadata{
						Name: "len",
						Type: sc.ToCompact(metadata.PrimitiveTypesU32),
					},
				},
				Output: sc.ToCompact(metadata.TypesTransactionPaymentFeeDetails),
				Docs:   sc.Sequence[sc.Str]{" Query fee details of a given encoded `Call`."},
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
	mockCall = new(mocks.Call)

	target := New(mockRuntimeDecoder, mockTransactionPayment, log.NewLogger())
	target.memUtils = mockMemoryUtils

	return target
}
