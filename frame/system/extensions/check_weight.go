package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckWeight types.Weight

func (_ CheckWeight) AdditionalSigned() (ok sc.Empty, err types.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (_ CheckWeight) Validate(_who *types.Address32, _call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	return DoValidate(info, length)
}

func (_ CheckWeight) ValidateUnsigned(_call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	return DoValidate(info, length)
}

func (_ CheckWeight) PreDispatch(_who *types.Address32, _call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = DoPreDispatch(info, length)
	return ok, err
}

func (_ CheckWeight) PreDispatchUnsigned(_call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = DoPreDispatch(info, length)
	return ok, err
}

func (_ CheckWeight) PostDispatch(_pre sc.Option[types.Pre], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, _length sc.Compact, _result *types.DispatchResult) (ok types.Pre, err types.TransactionValidityError) {
	unspent := postInfo.CalcUnspent(info)
	if unspent.AnyGt(types.WeightZero()) {
		currentWeight := system.StorageGetBlockWeight()
		currentWeight.Reduce(unspent, info.Class)
		system.StorageSetBlockWeight(currentWeight)
	}
	ok = types.Pre{}
	return ok, err
}

// Do the validate checks. This can be applied to both signed and unsigned.
//
// It only checks that the block weight and length limit will not exceed.
func DoValidate(info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	ok = types.DefaultValidTransaction()

	// ignore the next length. If they return `Ok`, then it is below the limit.
	_, err = checkBlockLength(info, length)
	if err != nil {
		return ok, err
	}

	// during validation we skip block limit check. Since the `validate_transaction`
	// call runs on an empty block anyway, by this we prevent `on_initialize` weight
	// consumption from causing false negatives.
	_, err = checkExtrinsicWeight(info)
	if err != nil {
		return ok, err
	}

	return ok, err
}

func DoPreDispatch(info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	nextLength, err := checkBlockLength(info, length)
	if err != nil {
		return ok, err
	}

	nextWeight, err := checkBlockWeight(info)
	if err != nil {
		return ok, err
	}

	_, err = checkExtrinsicWeight(info)
	if err != nil {
		return ok, err
	}

	system.StorageSetAllExtrinsicsLen(nextLength)
	system.StorageSetBlockWeight(nextWeight)

	return ok, err
}

// Checks if the current extrinsic can fit into the block with respect to block length limits.
//
// Upon successes, it returns the new block length as a `Result`.
func checkBlockLength(info *types.DispatchInfo, length sc.Compact) (ok sc.U32, err types.TransactionValidityError) {
	lengthLimit := system.DefaultBlockLength()
	currentLen := system.StorageGetAllExtrinsicsLen()
	addedLen := sc.U32(sc.U128(length).ToBigInt().Uint64())

	nextLen := currentLen.SaturatingAdd(addedLen)

	var maxLimit sc.U32
	if info.Class.Is(types.DispatchClassNormal) {
		maxLimit = lengthLimit.Max.Normal
	} else if info.Class.Is(types.DispatchClassOperational) {
		maxLimit = lengthLimit.Max.Operational
	} else if info.Class.Is(types.DispatchClassMandatory) {
		maxLimit = lengthLimit.Max.Mandatory
	} else {
		log.Critical("invalid DispatchClass type in CheckBlockLength()")
	}

	if nextLen > maxLimit {
		err = types.NewTransactionValidityError(types.NewInvalidTransactionExhaustsResources())
	} else {
		ok = sc.U32(sc.ToCompact(nextLen).ToBigInt().Uint64())
	}

	return ok, err
}

// Checks if the current extrinsic can fit into the block with respect to block weight limits.
//
// Upon successes, it returns the new block weight as a `Result`.
func checkBlockWeight(info *types.DispatchInfo) (ok types.ConsumedWeight, err types.TransactionValidityError) {
	maximumWeight := system.DefaultBlockWeights()
	allWeight := system.StorageGetBlockWeight()
	return CalculateConsumedWeight(maximumWeight, allWeight, info)
}

// Checks if the current extrinsic does not exceed the maximum weight a single extrinsic
// with given `DispatchClass` can have.
func checkExtrinsicWeight(info *types.DispatchInfo) (ok sc.Empty, err types.TransactionValidityError) {
	max := system.DefaultBlockWeights().Get(info.Class).MaxExtrinsic

	if max.HasValue {
		if info.Weight.AnyGt(max.Value) {
			err = types.NewTransactionValidityError(types.NewInvalidTransactionExhaustsResources())
		} else {
			ok = sc.Empty{}
		}
	}

	return ok, err
}

func CalculateConsumedWeight(maximumWeight system.BlockWeights, allWeight types.ConsumedWeight, info *types.DispatchInfo) (ok types.ConsumedWeight, err types.TransactionValidityError) {
	extrinsicWeight := info.Weight.SaturatingAdd(maximumWeight.Get(info.Class).BaseExtrinsic)
	limitPerClass := maximumWeight.Get(info.Class)

	// add the weight. If class is unlimited, use saturating add instead of checked one.
	if !limitPerClass.MaxTotal.HasValue && !limitPerClass.Reserved.HasValue {
		allWeight.SaturatingAdd(extrinsicWeight, info.Class)
	} else {
		// TODO:
		_, e := allWeight.CheckedAccrue(extrinsicWeight, info.Class)
		if e != nil {
			err = types.NewTransactionValidityError(types.NewInvalidTransactionExhaustsResources())
			return ok, err
		}
	}

	perClass := allWeight.Get(info.Class)

	// Check if we don't exceed per-class allowance
	switch limitPerClass.MaxTotal.HasValue {
	case true:
		max := limitPerClass.MaxTotal.Value
		if perClass.AnyGt(max) {
			err = types.NewTransactionValidityError(types.NewInvalidTransactionExhaustsResources())
			return ok, err
		}
	case false:
		// There is no `max_total` limit (`None`),
		// or we are below the limit.
		// TODO:
	}

	// In cases total block weight is exceeded, we need to fall back
	// to `reserved` pool if there is any.
	if allWeight.Total().AnyGt(maximumWeight.MaxBlock) {
		if limitPerClass.Reserved.HasValue {
			// We are over the limit in reserved pool.
			reserved := limitPerClass.Reserved.Value
			if perClass.AnyGt(reserved) {
				err = types.NewTransactionValidityError(types.NewInvalidTransactionExhaustsResources())
				return ok, err
			}
		} else {
			// There is either no limit in reserved pool (`None`),
			// or we are below the limit.
			// TODO:
		}
	}

	ok = allWeight

	return ok, err
}
