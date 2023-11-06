package extensions

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	mockCurrencyAdapter *mocks.CurrencyAdapter
	target              chargeTransaction

	who               = constants.ZeroAddressAccountId
	fee               = sc.NewU128(5)
	imbalance         = sc.NewU128(1)
	expectedImbalance = sc.NewOption[sc.U128](imbalance)
	tip               = sc.NewU128(0)
	reasons           = sc.U8(primitives.WithdrawReasonsTransactionPayment)

	correctedFee     = sc.NewU128(10)
	alreadyWithdrawn = sc.NewOption[sc.U128](sc.NewU128(11))
	refundAmount     = sc.NewU128(1)

	expectedError, _ = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionPayment())
)

func Test_ChargeTransaction_WithdrawFee_Success(t *testing.T) {
	setUp()
	mockCurrencyAdapter.On("Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive).Return(imbalance, nil)

	result, err := target.WithdrawFee(who, nil, nil, fee, tip)

	assert.Nil(t, err)
	assert.Equal(t, expectedImbalance, result)
	mockCurrencyAdapter.AssertCalled(t, "Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive)
}

func Test_ChargeTransaction_WithdrawFee_ZeroFee(t *testing.T) {
	setUp()
	expectedImbalance := sc.NewOption[sc.U128](nil)

	result, err := target.WithdrawFee(who, nil, nil, tip, tip)

	assert.Nil(t, err)
	assert.Equal(t, expectedImbalance, result)
	mockCurrencyAdapter.AssertNotCalled(t, "Withdraw")
}

func Test_ChargeTransaction_WithdrawFee_WithTip(t *testing.T) {
	setUp()

	reasons := sc.U8(primitives.WithdrawReasonsTip)
	mockCurrencyAdapter.On("Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive).Return(imbalance, nil)
	tip := fee

	result, err := target.WithdrawFee(who, nil, nil, fee, tip)

	assert.Nil(t, err)
	assert.Equal(t, expectedImbalance, result)
	mockCurrencyAdapter.AssertCalled(t, "Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive)
}

func Test_ChargeTransaction_WithdrawFee_Fail(t *testing.T) {
	setUp()
	expectedImbalance := sc.NewOption[sc.U128](nil)

	mockError := primitives.NewDispatchErrorBadOrigin()
	mockCurrencyAdapter.On("Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive).Return(imbalance, mockError)

	result, err := target.WithdrawFee(who, nil, nil, fee, tip)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, expectedImbalance, result)
	mockCurrencyAdapter.AssertCalled(t, "Withdraw", who, fee, reasons, primitives.ExistenceRequirementKeepAlive)
}

func Test_ChargeTransaction_CorrectAndDepositFee_AlreadyWithdrawn_Success(t *testing.T) {
	setUp()
	mockCurrencyAdapter.On("DepositIntoExisting", who, refundAmount).Return(refundAmount, nil)

	result := target.CorrectAndDepositFee(who, correctedFee, tip, alreadyWithdrawn)

	assert.Nil(t, result)
	mockCurrencyAdapter.AssertCalled(t, "DepositIntoExisting", who, refundAmount)
}

func Test_ChargeTransaction_CorrectAndDepositFee_NotWithdrawn(t *testing.T) {
	setUp()
	alreadyWithdrawn := sc.NewOption[sc.U128](nil)

	result := target.CorrectAndDepositFee(who, correctedFee, tip, alreadyWithdrawn)

	assert.Nil(t, result)
	mockCurrencyAdapter.AssertNotCalled(t, "DepositIntoExisting")
}

func Test_ChargeTransaction_CorrectAndDepositFee_AlreadyWithdrawn_DepositIntoExisting_Fail(t *testing.T) {
	setUp()
	mockCurrencyAdapter.On("DepositIntoExisting", who, refundAmount).Return(imbalance, primitives.NewDispatchErrorBadOrigin())

	result := target.CorrectAndDepositFee(who, correctedFee, tip, alreadyWithdrawn)

	assert.Equal(t, expectedError, result)
	mockCurrencyAdapter.AssertCalled(t, "DepositIntoExisting", who, refundAmount)
}

func Test_ChargeTransaction_CorrectAndDepositFee_AlreadyWithdrawn_Fail(t *testing.T) {
	setUp()
	positiveImbalance := sc.NewU128(50)
	mockCurrencyAdapter.On("DepositIntoExisting", who, refundAmount).Return(positiveImbalance, nil)

	result := target.CorrectAndDepositFee(who, correctedFee, tip, alreadyWithdrawn)

	assert.Equal(t, expectedError, result)
	mockCurrencyAdapter.AssertCalled(t, "DepositIntoExisting", who, refundAmount)
}

func setUp() {
	mockCurrencyAdapter = new(mocks.CurrencyAdapter)
	target = newChargeTransaction(mockCurrencyAdapter)
}
