package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	signerOption = sc.NewOption[types.Address32](constants.OneAddress)
	emptySigner  = sc.NewOption[types.Address32](nil)

	txSource     = types.NewTransactionSourceExternal()
	dispatchInfo = &types.DispatchInfo{
		Weight:  types.WeightFromParts(4, 5),
		Class:   types.NewDispatchClassMandatory(),
		PaysFee: types.NewPaysNo(),
	}
	length           = sc.ToCompact(5)
	postDispatchInfo = types.PostDispatchInfo{
		ActualWeight: sc.NewOption[types.Weight](types.WeightFromParts(2, 3)),
		PaysFee:      types.PaysYes,
	}
	errPostDispatchInfo = types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
		PostInfo: types.PostDispatchInfo{
			ActualWeight: sc.NewOption[types.Weight](nil),
			PaysFee:      types.PaysNo,
		},
		Error: types.NewDispatchErrorCorruption(),
	}
	pre            = sc.Sequence[types.Pre]{sc.NewVaryingData(sc.U32(1)), sc.NewVaryingData(sc.U32(7))}
	optionPre      = sc.NewOption[sc.Sequence[types.Pre]](pre)
	emptyOptionPre = sc.NewOption[sc.Sequence[types.Pre]](nil)
)

var (
	mockTransactional     *mocks.IoTransactional[types.PostDispatchInfo, types.DispatchError]
	mockUnsignedValidator *mocks.UnsignedValidator

	mockWithStorageLayer = mock.AnythingOfType("func() (types.PostDispatchInfo, goscale.VaryingData)")
)

func Test_CheckedExtrinsic_Function(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	assert.Equal(t, mockCall, target.Function())
}

func Test_CheckedExtrinsic_Apply_Signed_Success(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		Ok: postDispatchInfo,
	}
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatch", signerOption.Value, mockCall, dispatchInfo, length).
		Return(pre, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(postDispatchInfo, nil)
	mockSignedExtra.
		On("PostDispatch", optionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult).
		Return(types.TransactionValidityError(nil))

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatch", signerOption.Value, mockCall, dispatchInfo, length)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", optionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Signed_PreDispatch_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())

	mockSignedExtra.
		On("PreDispatch", signerOption.Value, mockCall, dispatchInfo, length).
		Return(pre, expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatch", signerOption.Value, mockCall, dispatchInfo, length)
	mockTransactional.AssertNotCalled(t, "WithStorageLayer", mock.Anything)
}

func Test_CheckedExtrinsic_Apply_Signed_WithStorageLayerErr(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: true,
		Err:      errPostDispatchInfo,
	}
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatch", signerOption.Value, mockCall, dispatchInfo, length).
		Return(pre, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(errPostDispatchInfo.PostInfo, errPostDispatchInfo.Error)
	mockSignedExtra.
		On("PostDispatch", optionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult).
		Return(types.TransactionValidityError(nil))

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatch", signerOption.Value, mockCall, dispatchInfo, length)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", optionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Signed_WithStorageLayerErr_PostDispatchErr(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: true,
		Err:      errPostDispatchInfo,
	}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatch", signerOption.Value, mockCall, dispatchInfo, length).
		Return(pre, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(errPostDispatchInfo.PostInfo, errPostDispatchInfo.Error)
	mockSignedExtra.
		On("PostDispatch", optionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult).
		Return(expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatch", signerOption.Value, mockCall, dispatchInfo, length)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", optionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Signed_PostDispatchFails(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		Ok: postDispatchInfo,
	}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatch", signerOption.Value, mockCall, dispatchInfo, length).
		Return(pre, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(postDispatchInfo, nil)
	mockSignedExtra.
		On("PostDispatch", optionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult).
		Return(expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatch", signerOption.Value, mockCall, dispatchInfo, length)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", optionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Unsigned_Success(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		Ok: postDispatchInfo,
	}
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(types.TransactionValidityError(nil))
	mockUnsignedValidator.On("PreDispatch", mockCall).Return(sc.Empty{}, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(postDispatchInfo, nil)
	mockSignedExtra.
		On("PostDispatch", emptyOptionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult).
		Return(types.TransactionValidityError(nil))

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertCalled(t, "PreDispatch", mockCall)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", emptyOptionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Unsigned_PreDispatchUnsigned_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertNotCalled(t, "PreDispatch", mock.Anything)
	mockTransactional.AssertNotCalled(t, "WithStorageLayer", mock.Anything)
}

func Test_CheckedExtrinsic_Apply_Unsigned_PreDispatch_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(types.TransactionValidityError(nil))
	mockUnsignedValidator.On("PreDispatch", mockCall).Return(sc.Empty{}, expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertCalled(t, "PreDispatch", mockCall)
	mockTransactional.AssertNotCalled(t, "WithStorageLayer", mock.Anything)
}

func Test_CheckedExtrinsic_Apply_Unsigned_WithStorageLayerErr(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: true,
		Err:      errPostDispatchInfo,
	}
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(types.TransactionValidityError(nil))
	mockUnsignedValidator.On("PreDispatch", mockCall).Return(sc.Empty{}, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(errPostDispatchInfo.PostInfo, errPostDispatchInfo.Error)
	mockSignedExtra.
		On("PostDispatch", emptyOptionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult).
		Return(types.TransactionValidityError(nil))

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertCalled(t, "PreDispatch", mockCall)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", emptyOptionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Unsigned_WithStorageLayerErr_PostDispatch_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: true,
		Err:      errPostDispatchInfo,
	}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(types.TransactionValidityError(nil))
	mockUnsignedValidator.On("PreDispatch", mockCall).Return(sc.Empty{}, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(errPostDispatchInfo.PostInfo, errPostDispatchInfo.Error)
	mockSignedExtra.
		On("PostDispatch", emptyOptionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult).
		Return(expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertCalled(t, "PreDispatch", mockCall)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", emptyOptionPre, dispatchInfo, &errPostDispatchInfo.PostInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Apply_Unsigned_PostDispatch_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		Ok: postDispatchInfo,
	}
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionStale())
	dispatchResult := types.NewDispatchResult(expect.Err)

	mockSignedExtra.
		On("PreDispatchUnsigned", mockCall, dispatchInfo, length).
		Return(types.TransactionValidityError(nil))
	mockUnsignedValidator.On("PreDispatch", mockCall).Return(sc.Empty{}, nil)
	mockTransactional.On("WithStorageLayer", mockWithStorageLayer).Return(postDispatchInfo, nil)
	mockSignedExtra.
		On("PostDispatch", emptyOptionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult).
		Return(expectError)

	result, err := target.Apply(mockUnsignedValidator, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "PreDispatchUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.AssertCalled(t, "PreDispatch", mockCall)
	mockTransactional.AssertCalled(t, "WithStorageLayer", mockWithStorageLayer)
	mockSignedExtra.
		AssertCalled(t, "PostDispatch", emptyOptionPre, dispatchInfo, &postDispatchInfo, length, &dispatchResult)
}

func Test_CheckedExtrinsic_Validate_Signed_Success(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	mockSignedExtra.
		On("Validate", signerOption.Value, mockCall, dispatchInfo, length).
		Return(types.DefaultValidTransaction(), nil)

	result, err := target.Validate(mockUnsignedValidator, txSource, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, types.DefaultValidTransaction(), result)
	mockSignedExtra.
		AssertCalled(t, "Validate", signerOption.Value, mockCall, dispatchInfo, length)
}

func Test_CheckedExtrinsic_Validate_Unsigned_Success(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expect := types.DefaultValidTransaction().CombineWith(types.DefaultValidTransaction())

	mockSignedExtra.
		On("ValidateUnsigned", mockCall, dispatchInfo, length).
		Return(types.DefaultValidTransaction(), nil)
	mockUnsignedValidator.
		On("ValidateUnsigned", txSource, mockCall).
		Return(types.DefaultValidTransaction(), nil)

	result, err := target.Validate(mockUnsignedValidator, txSource, dispatchInfo, length)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockSignedExtra.
		AssertCalled(t, "ValidateUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.
		AssertCalled(t, "ValidateUnsigned", txSource, mockCall)
}

func Test_CheckedExtrinsic_Validate_Unsigned_ValidateUnsigned_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)
	expectError := types.NewTransactionValidityError(types.NewInvalidTransactionPayment())

	mockSignedExtra.
		On("ValidateUnsigned", mockCall, dispatchInfo, length).
		Return(types.ValidTransaction{}, expectError)

	result, err := target.Validate(mockUnsignedValidator, txSource, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, types.ValidTransaction{}, result)
	mockSignedExtra.
		AssertCalled(t, "ValidateUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.
		AssertNotCalled(t, "ValidateUnsigned", mock.Anything, mock.Anything)
}

func Test_CheckedExtrinsic_Validate_Unsigned_ValidatorValidateUnsigned_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(emptySigner)

	expectError := types.NewTransactionValidityError(types.NewUnknownTransactionNoUnsignedValidator())

	mockSignedExtra.
		On("ValidateUnsigned", mockCall, dispatchInfo, length).
		Return(types.DefaultValidTransaction(), nil)
	mockUnsignedValidator.
		On("ValidateUnsigned", txSource, mockCall).
		Return(types.ValidTransaction{}, expectError)

	result, err := target.Validate(mockUnsignedValidator, txSource, dispatchInfo, length)

	assert.Equal(t, expectError, err)
	assert.Equal(t, types.ValidTransaction{}, result)
	mockSignedExtra.
		AssertCalled(t, "ValidateUnsigned", mockCall, dispatchInfo, length)
	mockUnsignedValidator.
		AssertCalled(t, "ValidateUnsigned", txSource, mockCall)
}

func Test_CheckedExtrinsic_dispatch_Success(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	args := sc.NewVaryingData(sc.U32(1))
	dispatchResult := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		Ok: postDispatchInfo,
	}

	mockCall.On("Args").Return(args)
	mockCall.On("Dispatch", types.RawOriginFrom(signerOption), args).Return(dispatchResult)

	res, err := target.dispatch(signerOption)

	assert.Nil(t, err)
	assert.Equal(t, postDispatchInfo, res)
	mockCall.AssertCalled(t, "Args")
	mockCall.AssertCalled(t, "Dispatch", types.RawOriginFrom(signerOption), args)
}

func Test_CheckedExtrinsic_dispatch_Fails(t *testing.T) {
	target := setupCheckedExtrinsic(signerOption)

	args := sc.NewVaryingData(sc.U32(1))
	dispatchResult := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: true,
		Ok:       postDispatchInfo,
		Err:      errPostDispatchInfo,
	}

	mockCall.On("Args").Return(args)
	mockCall.On("Dispatch", types.RawOriginFrom(signerOption), args).Return(dispatchResult)

	res, err := target.dispatch(signerOption)

	assert.Equal(t, errPostDispatchInfo.Error, err)
	assert.Equal(t, errPostDispatchInfo.PostInfo, res)
	mockCall.AssertCalled(t, "Args")
	mockCall.AssertCalled(t, "Dispatch", types.RawOriginFrom(signerOption), args)
}

func setupCheckedExtrinsic(signer sc.Option[types.Address32]) checkedExtrinsic {
	mockCall = new(mocks.Call)
	mockSignedExtra = new(mocks.SignedExtra)
	mockTransactional = new(mocks.IoTransactional[types.PostDispatchInfo, types.DispatchError])
	mockUnsignedValidator = new(mocks.UnsignedValidator)

	target := NewCheckedExtrinsic(signer, mockCall, mockSignedExtra).(checkedExtrinsic)
	target.transactional = mockTransactional

	return target
}