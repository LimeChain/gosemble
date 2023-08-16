package extensions

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ChargeTransactionPayment struct {
	fee                 primitives.Balance
	systemModule        system.Module
	txPaymentModule     transaction_payment.Module
	onChargeTransaction hooks.OnChargeTransaction
}

func NewChargeTransactionPayment(module system.Module, txPaymentModule transaction_payment.Module, currencyAdapter primitives.CurrencyAdapter) ChargeTransactionPayment {
	return ChargeTransactionPayment{
		systemModule:        module,
		txPaymentModule:     txPaymentModule,
		onChargeTransaction: newChargeTransaction(currencyAdapter),
	}
}

func (ctp ChargeTransactionPayment) Encode(buffer *bytes.Buffer) {
	sc.Compact(ctp.fee).Encode(buffer)
}

func (ctp *ChargeTransactionPayment) Decode(buffer *bytes.Buffer) {
	ctp.fee = sc.U128(sc.DecodeCompact(buffer))
}

func (ctp ChargeTransactionPayment) Bytes() []byte {
	return sc.EncodedBytes(ctp)
}

func (ctp ChargeTransactionPayment) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(), nil
}

func (ctp ChargeTransactionPayment) Validate(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	finalFee, _, err := ctp.withdrawFee(who, call, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	tip := ctp.fee
	validTransaction := primitives.DefaultValidTransaction()
	validTransaction.Priority = ctp.getPriority(info, length, tip, finalFee)

	return validTransaction, nil
}

func (ctp ChargeTransactionPayment) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctp ChargeTransactionPayment) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, imbalance, err := ctp.withdrawFee(who, call, info, length)
	if err != nil {
		return primitives.Pre{}, err
	}
	return sc.NewVaryingData(ctp.fee, *who, imbalance), nil
}

func (ctp ChargeTransactionPayment) PostDispatch(pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, length sc.Compact, result *primitives.DispatchResult) primitives.TransactionValidityError {
	if pre.HasValue {
		preValue := pre.Value

		tip := preValue[0].(primitives.Balance)
		who := preValue[1].(primitives.Address32)
		imbalance := preValue[2].(sc.Option[primitives.Balance])

		actualFee := ctp.txPaymentModule.ComputeActualFee(sc.U32(length.ToBigInt().Uint64()), *info, *postInfo, tip)
		err := ctp.onChargeTransaction.CorrectAndDepositFee(&who, actualFee, tip, imbalance)
		if err != nil {
			return err
		}

		ctp.systemModule.DepositEvent(transaction_payment.NewEventTransactionFeePaid(who.FixedSequence, actualFee, tip))
	}
	return nil
}

func (ctp ChargeTransactionPayment) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := ctp.ValidateUnsigned(call, info, length)
	return err
}

func (ctp ChargeTransactionPayment) getPriority(info *primitives.DispatchInfo, len sc.Compact, tip primitives.Balance, finalFee primitives.Balance) primitives.TransactionPriority {
	maxBlockWeight := ctp.systemModule.Constants.BlockWeights.MaxBlock.RefTime
	maxDefaultBlockLength := ctp.systemModule.Constants.BlockLength.Max
	maxBlockLength := sc.U64(*maxDefaultBlockLength.Get(info.Class))

	infoWeight := info.Weight.RefTime

	// info_weight.clamp(1, max_block_weight);
	boundedWeight := infoWeight
	if boundedWeight < 1 {
		boundedWeight = 1
	} else if boundedWeight > maxBlockWeight {
		boundedWeight = maxBlockWeight
	}

	// (len as u64).clamp(1, max_block_length);
	boundedLength := sc.U64(len.ToBigInt().Uint64())
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

	bnTip := new(big.Int).Add(tip.ToBigInt(), big.NewInt(1))

	scaledTip := new(big.Int).Mul(bnTip, new(big.Int).SetUint64(uint64(maxTxPerBlock)))

	if info.Class.Is(primitives.DispatchClassNormal) {
		return sc.U64(scaledTip.Uint64())
	} else if info.Class.Is(primitives.DispatchClassMandatory) {
		return sc.U64(scaledTip.Uint64())
	} else if info.Class.Is(primitives.DispatchClassOperational) {
		feeMultiplier := ctp.txPaymentModule.Constants.OperationalFeeMultiplier
		virtualTip := new(big.Int).Mul(finalFee.ToBigInt(), big.NewInt(int64(feeMultiplier)))
		scaledVirtualTip := new(big.Int).Mul(virtualTip, new(big.Int).SetUint64(uint64(maxTxPerBlock)))

		sum := new(big.Int).Add(scaledTip, scaledVirtualTip)

		return sc.U64(sum.Uint64())
	}

	return 0
}

func (ctp ChargeTransactionPayment) withdrawFee(who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Balance, sc.Option[primitives.Balance], primitives.TransactionValidityError) {
	tip := ctp.fee
	fee := ctp.txPaymentModule.ComputeFee(sc.U32(length.ToBigInt().Uint64()), *info, tip)

	imbalance, err := ctp.onChargeTransaction.WithdrawFee(who, _call, info, fee, sc.NewU128FromBigInt(tip.ToBigInt()))
	if err != nil {
		return primitives.Balance{}, sc.NewOption[primitives.Balance](nil), err
	}

	return fee, imbalance, nil
}
