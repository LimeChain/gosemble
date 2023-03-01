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

func (w CheckWeight) PreDispatchUnsigned(_call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = DoPreDispatch(info, length)
	return ok, err
}

func (_ CheckWeight) PostDispatch(_pre sc.Option[types.Pre], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, _length sc.Compact, _result *types.DispatchResult) (ok types.Pre, err types.TransactionValidityError) {
	// TODO:

	// let unspent = post_info.calc_unspent(info);
	// if unspent.any_gt(Weight::zero()) {
	// 	crate::BlockWeight::<T>::mutate(|current_weight| {
	// 		current_weight.reduce(unspent, info.class);
	// 	})
	// }
	// Ok(())

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
	// TODO:
	// let next_len = Self::check_block_length(info, len)?;
	// let next_weight = Self::check_block_weight(info)?;
	// Self::check_extrinsic_weight(info)?;

	// crate::AllExtrinsicsLen::<T>::put(next_len);
	// crate::BlockWeight::<T>::put(next_weight);
	// Ok(())
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
	if info.Class.Is(types.NormalDispatch) {
		maxLimit = lengthLimit.Max.Normal
	} else if info.Class.Is(types.OperationalDispatch) {
		maxLimit = lengthLimit.Max.Operational
	} else if info.Class.Is(types.MandatoryDispatch) {
		maxLimit = lengthLimit.Max.Mandatory
	} else {
		log.Critical("invalid DispatchClass type in CheckBlockLength()")
	}

	if nextLen > maxLimit {
		err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.ExhaustsResourcesError))
	} else {
		ok = sc.U32(sc.ToCompact(nextLen).ToBigInt().Uint64())
	}

	return ok, err
}

// Checks if the current extrinsic does not exceed the maximum weight a single extrinsic
// with given `DispatchClass` can have.
func checkExtrinsicWeight(info *types.DispatchInfo) (ok sc.Empty, err types.TransactionValidityError) {
	max := system.DefaultBlockWeights().Get(info.Class).MaxExtrinsic

	if max.HasValue {
		if info.Weight.AnyGt(max.Value) {
			err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.ExhaustsResourcesError))
		} else {
			ok = sc.Empty{}
		}
	}

	return ok, err
}
