package extrinsic

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	moduleId = sc.U8(3)
	txSource = primitives.NewTransactionSourceExternal()
)

var (
	mockCall             *mocks.Call
	mockRuntimeExtrinsic *mocks.RuntimeExtrinsic
	mockModule           *mocks.Module
)

func Test_UnsignedValidatorForChecked_PreDispatch(t *testing.T) {
	target := setupUnsignedValidator()

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, true)
	mockModule.On("PreDispatch", mockCall).Return(sc.Empty{}, nil)

	result, err := target.PreDispatch(mockCall)

	assert.Equal(t, sc.Empty{}, result)
	assert.Nil(t, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertCalled(t, "PreDispatch", mockCall)
}

func Test_UnsignedValidatorForChecked_PreDispatch_NotFound(t *testing.T) {
	target := setupUnsignedValidator()

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, false)

	result, err := target.PreDispatch(mockCall)

	assert.Equal(t, sc.Empty{}, result)
	assert.Nil(t, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertNotCalled(t, "PreDispatch", mock.Anything)
}

func Test_UnsignedValidatorForChecked_PreDispatch_Err(t *testing.T) {
	target := setupUnsignedValidator()
	expect := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, true)
	mockModule.On("PreDispatch", mockCall).Return(sc.Empty{}, expect)

	result, err := target.PreDispatch(mockCall)

	assert.Equal(t, sc.Empty{}, result)
	assert.Equal(t, expect, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertCalled(t, "PreDispatch", mockCall)
}

func Test_UnsignedValidatorForChecked_ValidateUnsigned(t *testing.T) {
	target := setupUnsignedValidator()
	expect := primitives.DefaultValidTransaction()

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, true)
	mockModule.On("ValidateUnsigned", txSource, mockCall).Return(primitives.DefaultValidTransaction(), nil)

	result, err := target.ValidateUnsigned(txSource, mockCall)

	assert.Equal(t, expect, result)
	assert.Nil(t, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertCalled(t, "ValidateUnsigned", txSource, mockCall)
}

func Test_UnsignedValidatorForChecked_ValidateUnsigned_NotFound(t *testing.T) {
	target := setupUnsignedValidator()
	expect := primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, false)

	result, err := target.ValidateUnsigned(txSource, mockCall)

	assert.Equal(t, primitives.ValidTransaction{}, result)
	assert.Equal(t, expect, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertNotCalled(t, "ValidateUnsigned", mock.Anything, mock.Anything)
}

func Test_UnsignedValidatorForChecked_ValidateUnsigned_Err(t *testing.T) {
	target := setupUnsignedValidator()
	expect := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())

	mockCall.On("ModuleIndex").Return(moduleId)
	mockRuntimeExtrinsic.On("Module", moduleId).Return(mockModule, true)
	mockModule.On("ValidateUnsigned", txSource, mockCall).Return(primitives.ValidTransaction{}, expect)

	result, err := target.ValidateUnsigned(txSource, mockCall)

	assert.Equal(t, primitives.ValidTransaction{}, result)
	assert.Equal(t, expect, err)
	mockCall.AssertCalled(t, "ModuleIndex")
	mockRuntimeExtrinsic.AssertCalled(t, "Module", moduleId)
	mockModule.AssertCalled(t, "ValidateUnsigned", txSource, mockCall)
}

func setupUnsignedValidator() primitives.UnsignedValidator {
	mockRuntimeExtrinsic = new(mocks.RuntimeExtrinsic)
	mockModule = new(mocks.Module)
	mockCall = new(mocks.Call)

	return NewUnsignedValidatorForChecked(mockRuntimeExtrinsic)
}
