package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ChargeTransactionPayment struct {
	fee                 primitives.Balance
	systemModule        system.Module
	txPaymentModule     transaction_payment.Module
	onChargeTransaction hooks.OnChargeTransaction
}

func NewChargeTransactionPayment(module system.Module, txPaymentModule transaction_payment.Module, currencyAdapter primitives.CurrencyAdapter) primitives.SignedExtension {
	return &ChargeTransactionPayment{
		systemModule:        module,
		txPaymentModule:     txPaymentModule,
		onChargeTransaction: newChargeTransaction(currencyAdapter),
	}
}

func (ctp ChargeTransactionPayment) Encode(buffer *bytes.Buffer) error {
	return sc.Compact(ctp.fee).Encode(buffer)
}

func (ctp *ChargeTransactionPayment) Decode(buffer *bytes.Buffer) error {
	fee, err := sc.DecodeCompact(buffer)
	if err != nil {
		return err
	}

	ctp.fee = sc.U128(fee)
	return nil
}

func (ctp ChargeTransactionPayment) Bytes() []byte {
	return sc.EncodedBytes(ctp)
}

func (ctp ChargeTransactionPayment) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return sc.NewVaryingData(), nil
}

func (ctp ChargeTransactionPayment) Validate(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	finalFee, _, err := ctp.withdrawFee(who, call, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	tip := ctp.fee
	validTransaction := primitives.DefaultValidTransaction()
	validTransaction.Priority = ctp.getPriority(info, length, tip, finalFee)

	return validTransaction, nil
}

func (ctp ChargeTransactionPayment) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctp ChargeTransactionPayment) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, error) {
	_, imbalance, err := ctp.withdrawFee(who, call, info, length)
	if err != nil {
		return primitives.Pre{}, err
	}
	return sc.NewVaryingData(ctp.fee, who, imbalance), nil
}

func (ctp ChargeTransactionPayment) PostDispatch(pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, length sc.Compact, result *primitives.DispatchResult) error {
	if pre.HasValue {
		preValue := pre.Value

		tip := preValue[0].(primitives.Balance)
		who := preValue[1].(primitives.AccountId)
		imbalance := preValue[2].(sc.Option[primitives.Balance])

		actualFee, err := ctp.txPaymentModule.ComputeActualFee(sc.U32(length.ToBigInt().Uint64()), *info, *postInfo, tip)
		if err != nil {
			return err
		}

		errFee := ctp.onChargeTransaction.CorrectAndDepositFee(who, actualFee, tip, imbalance)
		if errFee != nil {
			return errFee
		}

		ctp.systemModule.DepositEvent(
			transaction_payment.NewEventTransactionFeePaid(
				ctp.txPaymentModule.GetIndex(),
				who,
				actualFee,
				tip,
			),
		)
	}
	return nil
}

func (ctp ChargeTransactionPayment) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) error {
	_, err := ctp.ValidateUnsigned(call, info, length)
	return err
}

//func (ctp ChargeTransactionPayment) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
//	return primitives.NewMetadataTypeWithParam(
//			metadata.ChargeTransactionPayment,
//			"ChargeTransactionPayment",
//			sc.Sequence[sc.Str]{"pallet_transaction_payment", "ChargeTransactionPayment"},
//			primitives.NewMetadataTypeDefinitionComposite(
//				sc.Sequence[primitives.MetadataTypeDefinitionField]{
//					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU128, "BalanceOf<T>"),
//				},
//			),
//			primitives.NewMetadataEmptyTypeParameter("T"),
//		),
//		primitives.NewMetadataSignedExtension("ChargeTransactionPayment", metadata.ChargeTransactionPayment, metadata.TypesEmptyTuple)
//}

func (ctp ChargeTransactionPayment) getPriority(info *primitives.DispatchInfo, len sc.Compact, tip primitives.Balance, finalFee primitives.Balance) primitives.TransactionPriority {
	maxBlockWeight := ctp.systemModule.BlockWeights().MaxBlock.RefTime
	maxDefaultBlockLength := ctp.systemModule.BlockLength().Max

	value, err := maxDefaultBlockLength.Get(info.Class)
	if err != nil {
		log.Critical(err.Error())
	}
	maxBlockLength := sc.U64(*value)

	infoWeight := info.Weight.RefTime

	// info_weight.clamp(1, max_block_weight);
	boundedWeight := infoWeight // TODO: clamp
	if boundedWeight < 1 {
		boundedWeight = 1
	} else if boundedWeight > maxBlockWeight {
		boundedWeight = maxBlockWeight
	}

	// (len as u64).clamp(1, max_block_length);
	boundedLength := sc.U64(len.ToBigInt().Uint64()) // TODO: clamp
	if boundedLength < 1 {
		boundedLength = 1
	} else if boundedLength > maxBlockLength {
		boundedLength = maxBlockLength
	}

	maxTxPerBlockWeight := maxBlockWeight / boundedWeight
	maxTxPerBlockLength := maxBlockLength / boundedLength

	maxTxPerBlock := maxTxPerBlockWeight
	if maxTxPerBlockWeight > maxTxPerBlockLength {
		maxTxPerBlock = maxTxPerBlockLength
	}

	bnTip := tip.Add(sc.NewU128(1))

	scaledTip := bnTip.Mul(sc.NewU128(maxTxPerBlock))

	isNormal, infoClassErr := info.Class.Is(primitives.DispatchClassNormal)
	if infoClassErr != nil {
		log.Critical(infoClassErr.Error())
	}
	if isNormal {
		return sc.U64(scaledTip.ToBigInt().Uint64())
	}

	isMandatory, infoClassErr := info.Class.Is(primitives.DispatchClassMandatory)
	if infoClassErr != nil {
		log.Critical(infoClassErr.Error())
	}
	if isMandatory {
		return sc.U64(scaledTip.ToBigInt().Uint64())
	}

	isOperational, infoClassErr := info.Class.Is(primitives.DispatchClassOperational)
	if infoClassErr != nil {
		log.Critical(infoClassErr.Error())
	}
	if isOperational {
		feeMultiplier := ctp.txPaymentModule.OperationalFeeMultiplier()
		virtualTip := finalFee.Mul(sc.NewU128(feeMultiplier))
		scaledVirtualTip := virtualTip.Mul(sc.NewU128(maxTxPerBlock))

		sum := scaledTip.Add(scaledVirtualTip)

		return sc.U64(sum.ToBigInt().Uint64())
	}

	return 0
}

func (ctp ChargeTransactionPayment) withdrawFee(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Balance, sc.Option[primitives.Balance], error) {
	tip := ctp.fee
	fee, err := ctp.txPaymentModule.ComputeFee(sc.U32(length.ToBigInt().Uint64()), *info, tip)
	if err != nil {
		return primitives.Balance{}, sc.NewOption[primitives.Balance](nil), err
	}

	imbalance, errWithdraw := ctp.onChargeTransaction.WithdrawFee(who, call, info, fee, tip)
	if errWithdraw != nil {
		return primitives.Balance{}, sc.NewOption[primitives.Balance](nil), errWithdraw
	}

	return fee, imbalance, nil
}
