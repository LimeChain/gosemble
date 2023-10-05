package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckedExtrinsic interface {
	Signed() sc.Option[primitives.Address32]
	Function() primitives.Call
	Extra() primitives.SignedExtra

	Validate(validator UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError)
	Apply(validator UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], primitives.TransactionValidityError)
}

// CheckedExtrinsic is the definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
//
// TODO: make it generic
// generic::CheckedExtrinsic<AccountId, RuntimeCall, SignedExtra>;
type checkedExtrinsic struct {
	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	signed   sc.Option[primitives.Address32]
	function primitives.Call
	extra    primitives.SignedExtra
}

func NewCheckedExtrinsic(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) checkedExtrinsic {
	return checkedExtrinsic{
		signed:   signed,
		function: function,
		extra:    extra,
	}
}

func (xt checkedExtrinsic) Signed() sc.Option[primitives.Address32] {
	return xt.signed
}

func (xt checkedExtrinsic) Function() primitives.Call {
	return xt.function
}

func (xt checkedExtrinsic) Extra() primitives.SignedExtra {
	return xt.extra
}

func (xt checkedExtrinsic) Validate(validator UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	if xt.Signed().HasValue {
		id := xt.Signed().Value
		return xt.Extra().Validate(id, xt.function, info, length)
	}

	valid, err := xt.Extra().ValidateUnsigned(xt.function, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	unsignedValidation, err := validator.ValidateUnsigned(source, xt.function)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	return valid.CombineWith(unsignedValidation), nil
}

func (xt checkedExtrinsic) Apply(validator UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], primitives.TransactionValidityError) {
	var (
		maybeWho sc.Option[primitives.Address32]
		maybePre sc.Option[sc.Sequence[primitives.Pre]]
	)

	if xt.Signed().HasValue {
		id := xt.Signed().Value
		pre, err := xt.Extra().PreDispatch(id, xt.function, info, length)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}
		maybeWho, maybePre = sc.NewOption[primitives.Address32](id), sc.NewOption[sc.Sequence[primitives.Pre]](pre)
	} else {
		// Do any pre-flight stuff for an unsigned transaction.
		//
		// Note this function by default delegates to `ValidateUnsigned`, so that
		// all checks performed for the transaction queue are also performed during
		// the dispatch phase (applying the extrinsic).
		//
		// If you ever override this function, you need to make sure to always
		// perform the same validation as in `ValidateUnsigned`.
		err := xt.Extra().PreDispatchUnsigned(xt.function, info, length)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		_, err = validator.PreDispatch(xt.function)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		maybeWho, maybePre = sc.NewOption[primitives.Address32](nil), sc.NewOption[sc.Sequence[primitives.Pre]](nil)
	}

	var resWithInfo primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]

	support.WithStorageLayer(
		func() (primitives.PostDispatchInfo, primitives.DispatchError) {
			resWithInfo = xt.Function().Dispatch(primitives.RawOriginFrom(maybeWho), xt.Function().Args())

			if resWithInfo.HasError {
				return primitives.PostDispatchInfo{}, resWithInfo.Err.Error
			}

			return resWithInfo.Ok, nil
		},
	)

	var postInfo primitives.PostDispatchInfo
	if resWithInfo.HasError {
		postInfo = resWithInfo.Err.PostInfo
	} else {
		postInfo = primitives.PostDispatchInfo{
			ActualWeight: sc.NewOption[primitives.Weight](info.Weight),
			PaysFee:      info.PaysFee[0].(sc.U8),
		}
	}

	dispatchResult := primitives.NewDispatchResult(resWithInfo.Err)
	err := xt.Extra().PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	return resWithInfo, err
}
