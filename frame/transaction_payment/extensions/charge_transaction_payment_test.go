package extensions

import (
	"bytes"
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	extLen       = sc.U32(0)
	txTip        = sc.NewU128(10)
	txFee        = sc.NewU128(10)
	txImbalance  = sc.NewOption[types.Balance](sc.NewU128(0))
	whoAccountId = constants.ZeroAccountId
	postResult   = types.DispatchResult{}

	info = types.DispatchInfo{
		Weight:  types.WeightFromParts(0, 0),
		Class:   types.NewDispatchClassOperational(),
		PaysFee: types.PaysNo,
	}

	postInfo = types.PostDispatchInfo{
		ActualWeight: sc.NewOption[types.Weight](types.WeightFromParts(0, 0)),
		PaysFee:      0,
	}

	blockWeights = types.BlockWeights{
		BaseBlock: types.WeightFromParts(1, 2),
		MaxBlock:  types.WeightFromParts(3, 4),
		PerClass: types.PerDispatchClass[types.WeightsPerClass]{
			Normal: types.WeightsPerClass{
				BaseExtrinsic: types.WeightFromParts(5, 6),
			},
			Operational: types.WeightsPerClass{
				BaseExtrinsic: types.WeightFromParts(1, 1),
			},
			Mandatory: types.WeightsPerClass{
				BaseExtrinsic: types.WeightFromParts(9, 10),
			},
		},
	}

	blockLength = types.BlockLength{
		Max: types.PerDispatchClass[sc.U32]{
			Normal:      1,
			Operational: 2,
			Mandatory:   3,
		},
	}

	invalidTransactionPaymentError = types.NewTransactionValidityError(types.NewInvalidTransactionPayment())
)

var (
	targetChargeTxPayment                 ChargeTransactionPayment
	mockSystemModule                      *mocks.SystemModule
	mockTxPaymentModule                   *mocks.TransactionPaymentModule
	mockOnChargeTransaction               *mocks.OnChargeTransaction
	mockCurrencyAdapterForChargeTxPayment *mocks.CurrencyAdapter
	mockCall                              *mocks.Call
)

func setup(fee types.Balance) {
	mockSystemModule = new(mocks.SystemModule)
	mockTxPaymentModule = new(mocks.TransactionPaymentModule)
	mockOnChargeTransaction = new(mocks.OnChargeTransaction)
	mockCurrencyAdapterForChargeTxPayment = new(mocks.CurrencyAdapter)
	mockCall = new(mocks.Call)

	targetChargeTxPayment = ChargeTransactionPayment{
		systemModule:        mockSystemModule,
		txPaymentModule:     mockTxPaymentModule,
		onChargeTransaction: newChargeTransaction(mockCurrencyAdapterForChargeTxPayment),
	}

	targetChargeTxPayment.onChargeTransaction = mockOnChargeTransaction
	targetChargeTxPayment.fee = fee
}

func Test_Encode(t *testing.T) {
	setup(sc.NewU128(16383))

	buffer := bytes.NewBuffer([]byte{})

	err := targetChargeTxPayment.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0xfd, 0xff}, buffer.Bytes())
}

func Test_Decode(t *testing.T) {
	setup(sc.NewU128(0))

	buffer := bytes.NewBuffer([]byte{0xfd, 0xff})

	targetChargeTxPayment.Decode(buffer)

	assert.Equal(t, sc.NewU128(16383), targetChargeTxPayment.fee)
}

func Test_Bytes(t *testing.T) {
	setup(sc.NewU128(16383))

	assert.Equal(t, []byte{0xfd, 0xff}, targetChargeTxPayment.Bytes())
}

func Test_AdditionalSigned(t *testing.T) {
	setup(sc.NewU128(0))

	additionalSigned, err := targetChargeTxPayment.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, sc.NewVaryingData(), additionalSigned)
}

func Test_Validate_Error(t *testing.T) {
	setup(txFee)

	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, nil)
	mockOnChargeTransaction.On("WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee).
		Return(sc.NewOption[types.Balance](nil), invalidTransactionPaymentError)

	res, err := targetChargeTxPayment.Validate(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockTxPaymentModule.AssertCalled(t, "ComputeFee", extLen, info, txTip)
	mockOnChargeTransaction.AssertCalled(t, "WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee)

	assert.Equal(t, types.ValidTransaction{}, res)
	assert.Equal(t, invalidTransactionPaymentError, err)
}

func Test_Validate_Mandatory(t *testing.T) {
	setup(txFee)

	info := types.DispatchInfo{
		Weight:  types.WeightFromParts(0, 0),
		Class:   types.NewDispatchClassMandatory(),
		PaysFee: types.PaysNo,
	}
	expectedValidTransaction := types.DefaultValidTransaction()
	expectedValidTransaction.Priority = sc.U64(33)

	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, nil)
	mockOnChargeTransaction.On("WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee).
		Return(sc.NewOption[types.Balance](sc.NewU128(1)), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("BlockLength").Return(blockLength)

	res, err := targetChargeTxPayment.Validate(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockSystemModule.AssertCalled(t, "BlockWeights")
	mockSystemModule.AssertCalled(t, "BlockLength")
	mockTxPaymentModule.AssertNotCalled(t, "OperationalFeeMultiplier")
	assert.Nil(t, err)
	assert.Equal(t, expectedValidTransaction, res)

}

func Test_Validate_Operational_NoError(t *testing.T) {
	setup(txFee)

	expectedValidTransaction := types.DefaultValidTransaction()
	expectedValidTransaction.Priority = sc.U64(42)

	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, nil)
	mockOnChargeTransaction.On("WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee).
		Return(sc.NewOption[types.Balance](sc.NewU128(1)), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("BlockLength").Return(blockLength)

	mockTxPaymentModule.On("OperationalFeeMultiplier").Return(sc.U8(1))

	res, err := targetChargeTxPayment.Validate(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockSystemModule.AssertCalled(t, "BlockWeights")
	mockSystemModule.AssertCalled(t, "BlockLength")
	mockTxPaymentModule.AssertCalled(t, "OperationalFeeMultiplier")
	assert.Nil(t, err)
	assert.Equal(t, expectedValidTransaction, res)
}

func Test_ValidateUnsigned(t *testing.T) {
	setup(txFee)

	res, err := targetChargeTxPayment.ValidateUnsigned(mockCall, &types.DispatchInfo{}, sc.ToCompact(sc.U32(0)))

	assert.Equal(t, types.DefaultValidTransaction(), res)
	assert.Nil(t, err)
}

func Test_PreDispatch_Success(t *testing.T) {
	setup(txFee)

	imbalance := sc.NewOption[types.Balance](sc.NewU128(1))
	expectedResult := sc.NewVaryingData(txFee, whoAccountId, imbalance)

	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, nil)
	mockOnChargeTransaction.On("WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee).Return(imbalance, nil)

	res, err := targetChargeTxPayment.PreDispatch(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockTxPaymentModule.AssertCalled(t, "ComputeFee", extLen, info, txTip)
	mockOnChargeTransaction.AssertCalled(t, "WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee)
	assert.Nil(t, err)
	assert.Equal(t, expectedResult, res)
}

func Test_PreDispatch_Error(t *testing.T) {
	setup(txFee)

	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, nil)
	mockOnChargeTransaction.On("WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee).
		Return(sc.NewOption[types.Balance](nil), invalidTransactionPaymentError)

	res, err := targetChargeTxPayment.PreDispatch(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockTxPaymentModule.AssertCalled(t, "ComputeFee", extLen, info, txTip)
	mockOnChargeTransaction.AssertCalled(t, "WithdrawFee", whoAccountId, mockCall, &info, txTip, txFee)
	assert.Equal(t, invalidTransactionPaymentError, err)
	assert.Equal(t, types.Pre{}, res)
}

func Test_PreDispatch_CumputeFeeError(t *testing.T) {
	setup(txFee)
	expectedErr := errors.New("error")
	mockTxPaymentModule.On("ComputeFee", extLen, info, txTip).Return(txFee, expectedErr)

	_, err := targetChargeTxPayment.PreDispatch(whoAccountId, mockCall, &info, sc.ToCompact(extLen))

	mockTxPaymentModule.AssertCalled(t, "ComputeFee", extLen, info, txTip)
	assert.Equal(t, expectedErr, err)
}

func Test_PostDispatch_None(t *testing.T) {
	setup(txFee)

	pre := sc.NewOption[types.Pre](nil)

	err := targetChargeTxPayment.PostDispatch(pre, &info, &postInfo, sc.ToCompact(extLen), &postResult)

	mockTxPaymentModule.AssertNotCalled(t, "ComputeActualFee", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockSystemModule.AssertNotCalled(t, "DepositEvent", mock.Anything)
	assert.Nil(t, err)
}

func Test_PostDispatch_Some(t *testing.T) {
	setup(txFee)

	pre := sc.NewOption[types.Pre](
		sc.NewVaryingData(
			txFee,
			whoAccountId,
			txImbalance,
		),
	)

	actualFee := sc.NewU128(1)
	mockTxPaymentModule.On("ComputeActualFee", extLen, info, postInfo, txTip).Return(actualFee, nil)
	mockOnChargeTransaction.On("CorrectAndDepositFee", whoAccountId, actualFee, txTip, txImbalance).Return(nil)
	mockTxPaymentModule.On("GetIndex").Return(sc.U8(0))
	mockSystemModule.On("DepositEvent", mock.Anything)

	err := targetChargeTxPayment.PostDispatch(pre, &info, &postInfo, sc.ToCompact(extLen), &postResult)

	mockTxPaymentModule.AssertCalled(t, "ComputeActualFee", extLen, info, postInfo, txTip)
	mockOnChargeTransaction.AssertCalled(t, "CorrectAndDepositFee", whoAccountId, actualFee, txTip, txImbalance)
	mockSystemModule.AssertCalled(t, "DepositEvent", mock.Anything)
	assert.Nil(t, err)
}

func Test_PostDispatch_CorrectAndDepositFeeError(t *testing.T) {
	setup(txFee)

	pre := sc.NewOption[types.Pre](
		sc.NewVaryingData(
			txFee,
			whoAccountId,
			txImbalance,
		),
	)

	actualFee := sc.NewU128(1)
	mockTxPaymentModule.On("ComputeActualFee", extLen, info, postInfo, txTip).Return(actualFee, nil)
	mockOnChargeTransaction.On("CorrectAndDepositFee", whoAccountId, actualFee, txTip, txImbalance).
		Return(invalidTransactionPaymentError)

	err := targetChargeTxPayment.PostDispatch(pre, &info, &postInfo, sc.ToCompact(extLen), &postResult)

	mockTxPaymentModule.AssertCalled(t, "ComputeActualFee", extLen, info, postInfo, txTip)
	mockOnChargeTransaction.AssertCalled(t, "CorrectAndDepositFee", whoAccountId, actualFee, txTip, txImbalance)
	mockSystemModule.AssertNotCalled(t, "DepositEvent", mock.Anything)
	assert.Equal(t, invalidTransactionPaymentError, err)
}

func Test_PostDispatch_ComputeActualFeeError(t *testing.T) {
	setup(txFee)

	pre := sc.NewOption[types.Pre](
		sc.NewVaryingData(
			txFee,
			whoAccountId,
			txImbalance,
		),
	)
	expectedErr := errors.New("error")

	actualFee := sc.NewU128(1)
	mockTxPaymentModule.On("ComputeActualFee", extLen, info, postInfo, txTip).Return(actualFee, expectedErr)

	err := targetChargeTxPayment.PostDispatch(pre, &info, &postInfo, sc.ToCompact(extLen), &postResult)

	mockTxPaymentModule.AssertCalled(t, "ComputeActualFee", extLen, info, postInfo, txTip)
	assert.Equal(t, expectedErr, err)
}

func Test_PreDispatchUnsigned(t *testing.T) {
	setup(txFee)

	err := targetChargeTxPayment.PreDispatchUnsigned(mockCall, &types.DispatchInfo{}, sc.ToCompact(sc.U32(0)))

	assert.Nil(t, err)
}

func Test_Metadata(t *testing.T) {
	setup(txFee)

	metadataType, metadataSignedExtension := targetChargeTxPayment.Metadata()

	expectedMetadataType := types.NewMetadataTypeWithParam(
		metadata.ChargeTransactionPayment,
		"ChargeTransactionPayment",
		sc.Sequence[sc.Str]{"pallet_transaction_payment", "ChargeTransactionPayment"},
		types.NewMetadataTypeDefinitionComposite(
			sc.Sequence[types.MetadataTypeDefinitionField]{
				types.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU128, "BalanceOf<T>"),
			},
		),
		types.NewMetadataEmptyTypeParameter("T"),
	)

	expectedMetadataSignedExtension := types.NewMetadataSignedExtension(
		"ChargeTransactionPayment", metadata.ChargeTransactionPayment, metadata.TypesEmptyTuple,
	)

	assert.Equal(t, expectedMetadataType, metadataType)
	assert.Equal(t, expectedMetadataSignedExtension, metadataSignedExtension)
}

func Test_getPriority(t *testing.T) {
	setup(txFee)

	info := types.DispatchInfo{
		Weight:  types.WeightFromParts(7, 0),
		Class:   types.NewDispatchClassNormal(),
		PaysFee: types.PaysYes,
	}

	extLen = sc.U32(5)

	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("BlockLength").Return(blockLength)

	priority := targetChargeTxPayment.getPriority(&info, sc.ToCompact(extLen), txTip, txFee)

	mockSystemModule.AssertCalled(t, "BlockWeights")
	mockSystemModule.AssertCalled(t, "BlockLength")
	mockTxPaymentModule.AssertNotCalled(t, "OperationalFeeMultiplier")
	assert.Equal(t, sc.U64(11), priority)
}
