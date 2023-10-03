package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	dustCleanerAccount  = constants.ZeroAddress
	expectedDustCleaner = dustCleaner{
		accountId: dustCleanerAccount,
		negativeImbalance: sc.NewOption[negativeImbalance](negativeImbalance{
			Balance: issuanceBalance,
		}),
		eventDepositor: nil,
	}
	issuanceBalance = sc.NewU128(123)

	mockStorageTotalIssuance *mocks.StorageValue[sc.U128]
	mockEventDepositor       *mocks.EventDepositor
)

func Test_NegativeImbalance_New(t *testing.T) {
	target := setupNegativeImbalance()

	assert.Equal(t, negativeImbalance{issuanceBalance, mockStorageTotalIssuance}, target)
}

func Test_NegativeImbalance_Drop(t *testing.T) {
	target := setupNegativeImbalance()
	mockStorageTotalIssuance.On("Get").Return(sc.NewU128(5))
	mockStorageTotalIssuance.On("Put", sc.NewU128(0)).Return()

	target.Drop()

	mockStorageTotalIssuance.AssertCalled(t, "Get")
	mockStorageTotalIssuance.AssertCalled(t, "Put", sc.NewU128(0))
}

func Test_PositiveImbalance_New(t *testing.T) {
	target := setupPositiveImbalance()

	assert.Equal(t, positiveImbalance{issuanceBalance, mockStorageTotalIssuance}, target)
}

func Test_PositiveImbalance_Drop(t *testing.T) {
	target := setupPositiveImbalance()
	mockStorageTotalIssuance.On("Get").Return(sc.NewU128(5))
	mockStorageTotalIssuance.On("Put", sc.NewU128(128)).Return()

	target.Drop()

	mockStorageTotalIssuance.AssertCalled(t, "Get")
	mockStorageTotalIssuance.AssertCalled(t, "Put", sc.NewU128(128))
}

func Test_DustCleanerValue_New(t *testing.T) {
	target := setupDustCleanerValue()
	expected := dustCleaner{
		moduleIndex: moduleId,
		accountId:   constants.ZeroAddress,
		negativeImbalance: sc.NewOption[negativeImbalance](negativeImbalance{
			Balance:       issuanceBalance,
			totalIssuance: mockStorageTotalIssuance,
		}),
		eventDepositor: mockEventDepositor,
	}

	assert.Equal(t, expected, target)
}

func Test_DustCleanerValue_Encode(t *testing.T) {
	target := setupDustCleanerValue()

	buffer := &bytes.Buffer{}
	target.Encode(buffer)

	assert.Equal(t, expectedDustCleaner.Bytes(), buffer.Bytes())
}

func Test_DustCleanerValue_Bytes(t *testing.T) {
	target := setupDustCleanerValue()

	assert.Equal(t, expectedDustCleaner.Bytes(), target.Bytes())
}

func Test_DustCleanerValue_Drop(t *testing.T) {
	expectedEvent := newEventDustLost(moduleId, dustCleanerAccount.FixedSequence, issuanceBalance)
	target := setupDustCleanerValue()
	mockEventDepositor.On("DepositEvent", expectedEvent).Return()
	mockStorageTotalIssuance.On("Get").Return(sc.NewU128(5))
	mockStorageTotalIssuance.On("Put", sc.NewU128(0))

	target.Drop()

	mockEventDepositor.AssertCalled(t, "DepositEvent", expectedEvent)
	mockStorageTotalIssuance.AssertCalled(t, "Get")
	mockStorageTotalIssuance.AssertCalled(t, "Put", sc.NewU128(0))
}

func setupNegativeImbalance() negativeImbalance {
	mockStorageTotalIssuance = new(mocks.StorageValue[sc.U128])
	return newNegativeImbalance(issuanceBalance, mockStorageTotalIssuance)
}

func setupPositiveImbalance() positiveImbalance {
	mockStorageTotalIssuance = new(mocks.StorageValue[sc.U128])
	return newPositiveImbalance(issuanceBalance, mockStorageTotalIssuance)
}

func setupDustCleanerValue() dustCleaner {
	mockEventDepositor = new(mocks.EventDepositor)
	return newDustCleaner(moduleId, dustCleanerAccount, sc.NewOption[negativeImbalance](setupNegativeImbalance()), mockEventDepositor)
}
