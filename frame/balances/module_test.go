package balances

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Module_GetIndex(t *testing.T) {
	assert.Equal(t, sc.U8(moduleId), setupModule().GetIndex())
}

func Test_Module_name(t *testing.T) {
	assert.Equal(t, name, setupModule().name())
}

func Test_Module_Functions(t *testing.T) {
	target := setupModule()
	functions := target.Functions()

	assert.Equal(t, 6, len(functions))
}

func Test_Module_PreDispatch(t *testing.T) {
	target := setupModule()

	result, err := target.PreDispatch(setupCallTransfer())

	assert.Nil(t, err)
	assert.Equal(t, sc.Empty{}, result)
}

func Test_Module_ValidateUnsigned(t *testing.T) {
	target := setupModule()

	result, err := target.ValidateUnsigned(primitives.TransactionSource{}, setupCallTransfer())

	assert.Equal(t, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator()), err)
	assert.Equal(t, primitives.ValidTransaction{}, result)
}

func Test_Module_DepositIntoExisting_Success(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	tryMutateResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewOption[sc.U128](nil), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{Value: targetValue}),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(tryMutateResult)

	result, err := target.DepositIntoExisting(fromAddress.AsAddress32(), targetValue)

	assert.Equal(t, targetValue, result)
	assert.Nil(t, err)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockTotalIssuance.AssertNotCalled(t, "Get")
	mockTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Module_DepositIntoExisting_ZeroValue(t *testing.T) {
	target := setupModule()

	result, err := target.DepositIntoExisting(fromAddress.AsAddress32(), sc.NewU128(0))

	assert.Equal(t, sc.NewU128(0), result)
	assert.Nil(t, err)
	mockStoredMap.AssertNotCalled(t, "TryMutateExists", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_DepositIntoExisting_TryMutateAccount_Fails(t *testing.T) {
	target := setupModule()
	expectError := primitives.NewDispatchErrorCannotLookup()
	mockReturn := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectError,
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(mockReturn)

	result, err := target.DepositIntoExisting(fromAddress.AsAddress32(), targetValue)

	assert.Equal(t, sc.U128{}, result)
	assert.Equal(t, expectError, err)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_Withdraw_Success(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	tryMutateResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewOption[sc.U128](nil), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(tryMutateResult)

	result, err := target.Withdraw(fromAddress.AsAddress32(), targetValue, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, targetValue, result)
	assert.Nil(t, err)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockTotalIssuance.AssertNotCalled(t, "Get")
	mockTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Module_Withdraw_ZeroValue(t *testing.T) {
	target := setupModule()

	result, err := target.Withdraw(fromAddress.AsAddress32(), sc.NewU128(0), sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, sc.NewU128(0), result)
	assert.Nil(t, err)
	mockStoredMap.AssertNotCalled(t, "TryMutateExists", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_Withdraw_TryMutateAccount_Fails(t *testing.T) {
	target := setupModule()
	expectError := primitives.NewDispatchErrorCannotLookup()
	mockReturn := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectError,
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(mockReturn)

	result, err := target.Withdraw(fromAddress.AsAddress32(), targetValue, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, sc.U128{}, result)
	assert.Equal(t, expectError, err)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_ensureCanWithdraw_Success(t *testing.T) {
	target := setupModule()

	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)

	result := target.ensureCanWithdraw(fromAddress.AsAddress32(), targetValue, primitives.ReasonsFee, sc.NewU128(5))

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
}

func Test_Module_ensureCanWithdraw_ZeroAmount(t *testing.T) {
	target := setupModule()

	result := target.ensureCanWithdraw(fromAddress.AsAddress32(), sc.NewU128(0), primitives.ReasonsFee, sc.NewU128(5))

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStoredMap.AssertNotCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
}

func Test_Module_ensureCanWithdraw_LiquidityRestrictions(t *testing.T) {
	target := setupModule()
	expected := primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Error:   sc.U32(ErrorLiquidityRestrictions),
		Message: sc.NewOption[sc.Str](nil),
	})
	frozenAccountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{
			MiscFrozen: sc.NewU128(10),
			FeeFrozen:  sc.NewU128(11),
		},
	}
	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(frozenAccountInfo)

	result := target.ensureCanWithdraw(fromAddress.AsAddress32(), targetValue, primitives.ReasonsFee, sc.NewU128(5))

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
}

func Test_Module_tryMutateAccount_Success(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	tryMutateResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewOption[sc.U128](nil), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}
	expected := sc.Result[sc.Encodable]{}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(tryMutateResult)

	result := target.tryMutateAccount(fromAddress.AsAddress32(), func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] { return sc.Result[sc.Encodable]{} })

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
}

func Test_Module_tryMutateAccount_TryMutateAccountWithDust_Fails(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorCannotLookup(),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(expected)

	result := target.tryMutateAccount(fromAddress.AsAddress32(), func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] { return sc.Result[sc.Encodable]{} })

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_tryMutateAccountWithDust_Success(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	tryMutateResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewOption[sc.U128](nil), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}
	expected := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.Result[sc.Encodable]{}, newDustCleaner(moduleId, fromAddress.AsAddress32(), sc.NewOption[negativeImbalance](nil), mockStoredMap)),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(tryMutateResult)

	result := target.tryMutateAccountWithDust(fromAddress.AsAddress32(), func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] { return sc.Result[sc.Encodable]{} })

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
}

func Test_Module_tryMutateAccountWithDust_Success_Endowed(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	tryMutateResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewOption[sc.U128](targetValue), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}
	expected := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.Result[sc.Encodable]{}, newDustCleaner(moduleId, fromAddress.AsAddress32(), sc.NewOption[negativeImbalance](nil), mockStoredMap)),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(tryMutateResult)
	mockStoredMap.On("DepositEvent", newEventEndowed(moduleId, fromAddress.AsAddress32().FixedSequence, targetValue))

	result := target.tryMutateAccountWithDust(fromAddress.AsAddress32(), func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] { return sc.Result[sc.Encodable]{} })

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertCalled(t, "DepositEvent", newEventEndowed(moduleId, fromAddress.AsAddress32().FixedSequence, targetValue))
}

func Test_Module_tryMutateAccountWithDust_TryMutateExists_Fail(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorCannotLookup(),
	}

	mockStoredMap.On("TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]")).Return(expected)

	result := target.tryMutateAccountWithDust(fromAddress.AsAddress32(), func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] { return sc.Result[sc.Encodable]{} })

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "TryMutateExists", fromAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_mutateAccount_Success(t *testing.T) {
	target := setupModule()
	target.storage.TotalIssuance = new(mocks.StorageValue[sc.U128])
	maybeAccount := &primitives.AccountData{}
	expected := sc.Result[sc.Encodable]{
		HasError: false,
		Value:    sc.NewVaryingData(sc.NewOption[sc.U128](sc.NewU128(0)), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}

	result := target.
		mutateAccount(
			maybeAccount,
			func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
				return sc.Result[sc.Encodable]{}
			},
		)

	assert.Equal(t, expected, result)
}

func Test_Module_mutateAccount_f_result(t *testing.T) {
	target := setupModule()
	target.storage.TotalIssuance = new(mocks.StorageValue[sc.U128])
	maybeAccount := &primitives.AccountData{
		Free: sc.NewU128(2),
	}
	e := primitives.NewDispatchErrorBadOrigin()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    e,
	}

	result := target.
		mutateAccount(
			maybeAccount,
			func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value:    e,
				}
			},
		)

	assert.Equal(t, expected, result)
}

func Test_Module_mutateAccount_Success_NotNewAccount(t *testing.T) {
	target := setupModule()
	target.storage.TotalIssuance = new(mocks.StorageValue[sc.U128])
	maybeAccount := &primitives.AccountData{
		Free: sc.NewU128(2),
	}
	expected := sc.Result[sc.Encodable]{
		HasError: false,
		Value:    sc.NewVaryingData(sc.NewOption[sc.U128](nil), sc.NewOption[negativeImbalance](nil), sc.Result[sc.Encodable]{}),
	}

	result := target.
		mutateAccount(
			maybeAccount,
			func(who *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
				return sc.Result[sc.Encodable]{}
			},
		)

	assert.Equal(t, expected, result)
}

func Test_Module_postMutation_Success(t *testing.T) {
	target := setupModule()

	accOption, imbalance := target.postMutation(*fromAccountData)

	assert.Equal(t, sc.NewOption[primitives.AccountData](*fromAccountData), accOption)
	assert.Equal(t, sc.NewOption[negativeImbalance](nil), imbalance)
}

func Test_Module_postMutation_ZeroTotal(t *testing.T) {
	target := setupModule()

	fromAccountData.Free = sc.NewU128(0)

	accOption, imbalance := target.postMutation(*fromAccountData)

	assert.Equal(t, sc.NewOption[primitives.AccountData](nil), accOption)
	assert.Equal(t, sc.NewOption[negativeImbalance](nil), imbalance)
}

func Test_Module_postMutation_LessExistentialDeposit(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance
	target.Constants.ExistentialDeposit = sc.NewU128(6)

	accOption, imbalance := target.postMutation(*fromAccountData)

	assert.Equal(t, sc.NewOption[primitives.AccountData](nil), accOption)
	assert.Equal(t, sc.NewOption[negativeImbalance](newNegativeImbalance(fromAccountData.Total(), target.storage.TotalIssuance)), imbalance)
}

func Test_Module_withdraw_Success(t *testing.T) {
	target := setupModule()
	value := sc.NewU128(3)

	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("DepositEvent", newEventWithdraw(moduleId, fromAddress.AsAddress32().FixedSequence, value))

	result := target.withdraw(fromAddress.AsAddress32(), value, fromAccountData, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, sc.Result[sc.Encodable]{Value: value}, result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	assert.Equal(t, sc.NewU128(2), fromAccountData.Free)
	mockStoredMap.AssertCalled(t, "DepositEvent", newEventWithdraw(moduleId, fromAddress.AsAddress32().FixedSequence, value))
}

func Test_Module_withdraw_InsufficientBalance(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorInsufficientBalance),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}
	value := sc.NewU128(10)

	result := target.withdraw(fromAddress.AsAddress32(), value, fromAccountData, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_withdraw_KeepAlive(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorKeepAlive),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}

	result := target.withdraw(fromAddress.AsAddress32(), targetValue, fromAccountData, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_withdraw_CannotWithdraw(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}
	value := sc.NewU128(3)

	frozenAccountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{
			MiscFrozen: sc.NewU128(10),
			FeeFrozen:  sc.NewU128(11),
		},
	}
	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(frozenAccountInfo)

	result := target.withdraw(fromAddress.AsAddress32(), value, fromAccountData, sc.U8(primitives.ReasonsFee), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_deposit_Success(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		Value: targetValue,
	}

	mockStoredMap.On("DepositEvent", newEventDeposit(moduleId, toAddress.AsAddress32().FixedSequence, targetValue))

	result := target.deposit(toAddress.AsAddress32(), toAccountData, false, targetValue)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockStoredMap.AssertCalled(t, "DepositEvent", newEventDeposit(moduleId, toAddress.AsAddress32().FixedSequence, targetValue))
}

func Test_Module_deposit_DeadAccount(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorDeadAccount),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}

	result := target.deposit(toAddress.AsAddress32(), toAccountData, true, targetValue)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(1), toAccountData.Free)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_deposit_ArithmeticOverflow(t *testing.T) {
	target := setupModule()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorArithmetic(primitives.NewArithmeticErrorOverflow()),
	}
	toAccountData.Free = sc.MaxU128()

	result := target.deposit(toAddress.AsAddress32(), toAccountData, false, targetValue)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.MaxU128(), toAccountData.Free)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Module_Metadata(t *testing.T) {
	target := setupModule()

	expectMetadataTypes := sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesBalancesEvent, "pallet_balances pallet Event", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Endowed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free_balance", "T::Balance"),
					},
					EventEndowed,
					"Event.Endowed"),
				primitives.NewMetadataDefinitionVariant(
					"DustLost",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventDustLost,
					"Events.DustLost"),
				primitives.NewMetadataDefinitionVariant(
					"Transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventTransfer,
					"Events.Transfer"),
				primitives.NewMetadataDefinitionVariant(
					"BalanceSet",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "T::Balance"),
					},
					EventBalanceSet,
					"Events.BalanceSet"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventReserved,
					"Events.Reserved"),
				primitives.NewMetadataDefinitionVariant(
					"Unreserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventUnreserved,
					"Events.Unreserved"),
				primitives.NewMetadataDefinitionVariant(
					"ReserveRepatriated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesBalanceStatus, "destination_status", "Status"),
					},
					EventReserveRepatriated,
					"Events.ReserveRepatriated"),
				primitives.NewMetadataDefinitionVariant(
					"Deposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventDeposit,
					"Event.Deposit"),
				primitives.NewMetadataDefinitionVariant(
					"Withdraw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventWithdraw,
					"Event.Withdraw"),
				primitives.NewMetadataDefinitionVariant(
					"Slashed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventSlashed,
					"Event.Slashed"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesBalanceStatus,
			"BalanceStatus",
			sc.Sequence[sc.Str]{"frame_support", "traits", "tokens", "misc", "BalanceStatus"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Free",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.BalanceStatusFree,
						"BalanceStatus.Free"),
					primitives.NewMetadataDefinitionVariant(
						"Reserved",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.BalanceStatusReserved,
						"BalanceStatus.Reserved"),
				})),

		primitives.NewMetadataTypeWithParams(metadata.TypesBalancesErrors,
			"pallet_balances pallet Error",
			sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Error"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"VestingBalance",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorVestingBalance,
						"Vesting balance too high to send value"),
					primitives.NewMetadataDefinitionVariant(
						"LiquidityRestrictions",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorLiquidityRestrictions,
						"Account liquidity restrictions prevent withdrawal"),
					primitives.NewMetadataDefinitionVariant(
						"InsufficientBalance",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorInsufficientBalance,
						"Balance too low to send value."),
					primitives.NewMetadataDefinitionVariant(
						"ExistentialDeposit",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorExistentialDeposit,
						"Value too low to create account due to existential deposit"),
					primitives.NewMetadataDefinitionVariant(
						"KeepAlive",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorKeepAlive,
						"Transfer/payment would kill account"),
					primitives.NewMetadataDefinitionVariant(
						"ExistingVestingSchedule",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorExistingVestingSchedule,
						"A vesting schedule already exists for this account"),
					primitives.NewMetadataDefinitionVariant(
						"DeadAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorDeadAccount,
						"Beneficiary account must pre-exist"),
					primitives.NewMetadataDefinitionVariant(
						"TooManyReserves",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorTooManyReserves,
						"Number of named reserves exceed MaxReserves"),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),

		primitives.NewMetadataTypeWithParams(metadata.BalancesCalls, "Balances calls", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionTransferIndex,
					"Transfer some liquid free balance to another account."),
				primitives.NewMetadataDefinitionVariant(
					"set_balance",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "who", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "new_free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "new_reserved", "T::Balance"),
					},
					functionSetBalanceIndex,
					"Set the balances of a given account."),
				primitives.NewMetadataDefinitionVariant(
					"force_transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "source", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionForceTransferIndex,
					"Exactly as `transfer`, except the origin must be root and the source account may be specified."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_keep_alive",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionTransferKeepAliveIndex,
					"Same as the [`transfer`] call, but with a check that the transfer will not kill the origin account."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_all",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesBool, "keep_alive", "bool"),
					},
					functionTransferAllIndex,
					"Transfer the entire transferable balance from the caller account."),
				primitives.NewMetadataDefinitionVariant(
					"force_unreserve",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "who", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					functionForceFreeIndex,
					"Unreserve some balance from a user by force."),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),
	}
	expectMetadataModule := primitives.MetadataModule{
		Name: name,
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: name,
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"TotalIssuance",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesU128)),
					"The total units issued in the system."),
			},
		}),
		Call: sc.NewOption[sc.Compact](sc.ToCompact(metadata.BalancesCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				name,
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.BalancesCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Balances, Runtime>"),
				},
				moduleId,
				"Call.Balances"),
		),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesEvent)),
		EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				name,
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesBalancesEvent, "pallet_balances::Event<Runtime>"),
				},
				moduleId,
				"Events.Balances"),
		),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"ExistentialDeposit",
				sc.ToCompact(metadata.PrimitiveTypesU128),
				sc.BytesToSequenceU8(existentialDeposit.Bytes()),
				"The minimum amount required to keep an account open. MUST BE GREATER THAN ZERO!",
			),
			primitives.NewMetadataModuleConstant(
				"MaxLocks",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(maxLocks.Bytes()),
				"The maximum number of locks that should exist on an account.  Not strictly enforced, but used for weight estimation.",
			),
			primitives.NewMetadataModuleConstant(
				"MaxReserves",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(maxReserves.Bytes()),
				"The maximum number of named reserves that can exist on an account.",
			),
		},
		Error: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesErrors)),
		Index: moduleId,
	}

	resultTypes, resultMetadataModule := target.Metadata()

	assert.Equal(t, expectMetadataTypes, resultTypes)
	assert.Equal(t, expectMetadataModule, resultMetadataModule)
}

func setupModule() Module {
	mockStoredMap = new(mocks.StoredMap)
	config := NewConfig(dbWeight, maxLocks, maxReserves, existentialDeposit, mockStoredMap)

	fromAccountData = &primitives.AccountData{
		Free: sc.NewU128(5),
	}

	toAccountData = &primitives.AccountData{
		Free: sc.NewU128(1),
	}

	return New(moduleId, config)
}
