package transaction_payment

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var DefaultMultiplierValue = sc.NewU128FromUint64(1)
var DefaultTip = sc.NewU128FromUint64(0)

func QueryInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := types.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	partialFee := sc.NewU128FromUint64(0)
	if ext.IsSigned() {
		partialFee = computeFee(length, dispatchInfo, DefaultTip)
	}

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

func QueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := types.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	var feeDetails primitives.FeeDetails
	if ext.IsSigned() {
		feeDetails = computeFeeDetails(length, dispatchInfo, DefaultTip)
	} else {
		feeDetails = primitives.FeeDetails{
			InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		}
	}

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}

func QueryCallInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := types.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	partialFee := computeFee(length, dispatchInfo, DefaultTip)

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

func QueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := types.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	feeDetails := computeFeeDetails(length, dispatchInfo, DefaultTip)

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}

func computeFee(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.Balance {
	return computeFeeDetails(len, info, tip).FinalFee()
}

func computeFeeDetails(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.FeeDetails {
	return computeFeeRaw(len, info.Weight, tip, info.PaysFee, info.Class)
}

func computeFeeRaw(len sc.U32, weight primitives.Weight, tip primitives.Balance, paysFee primitives.Pays, class primitives.DispatchClass) primitives.FeeDetails {
	if paysFee[0] == primitives.PaysYes { // TODO: type safety
		unadjustedWeightFee := weightToFee(weight)
		multiplier := storageNextFeeMultiplier()

		bnAdjustedWeightFee := new(big.Int).Mul(multiplier.ToBigInt(), unadjustedWeightFee.ToBigInt())
		adjustedWeightFee := sc.NewU128FromBigInt(bnAdjustedWeightFee)

		lenFee := lengthToFee(len)
		baseFee := weightToFee(system.DefaultBlockWeights().Get(class).BaseExtrinsic)

		inclusionFee := sc.NewOption[primitives.InclusionFee](primitives.NewInclusionFee(baseFee, lenFee, adjustedWeightFee))

		return primitives.FeeDetails{
			InclusionFee: inclusionFee,
			Tip:          tip,
		}
	}

	return primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		Tip:          tip,
	}
}

func lengthToFee(length sc.U32) primitives.Balance {
	return constants.LengthToFee.WeightToFee(primitives.WeightFromParts(sc.U64(length), 0))
}

func weightToFee(weight primitives.Weight) primitives.Balance {
	cappedWeight := weight.Min(system.DefaultBlockWeights().MaxBlock)

	return constants.WeightToFee.WeightToFee(cappedWeight)
}

func storageNextFeeMultiplier() sc.U128 {
	txPaymentHash := hashing.Twox128(constants.KeyTransactionPayment)
	nextFeeMultiplierHash := hashing.Twox128(constants.KeyNextFeeMultiplier)
	key := append(txPaymentHash, nextFeeMultiplierHash...)

	return storage.GetDecodeOnEmpty(key, sc.DecodeU128, DefaultMultiplierValue)
}
