package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// CheckedExtrinsic is the definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
//
// TODO: make it generic
// generic::CheckedExtrinsic<AccountId, RuntimeCall, SignedExtra>;
type checkedExtrinsic struct {
	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	signed        sc.Option[primitives.Address32]
	function      primitives.Call
	extra         primitives.SignedExtra
	transactional support.Transactional[primitives.PostDispatchInfo, primitives.DispatchError]
}

func (c checkedExtrinsic) Apply(validator primitives.UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], primitives.TransactionValidityError) {
	var (
		maybeWho sc.Option[primitives.Address32]
		maybePre sc.Option[sc.Sequence[primitives.Pre]]
	)

	if c.signed.HasValue {
		id := c.signed.Value
		pre, err := c.extra.PreDispatch(id, c.function, info, length)
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
		err := c.extra.PreDispatchUnsigned(c.function, info, length)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		_, err = validator.PreDispatch(c.function)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		maybeWho, maybePre = sc.NewOption[primitives.Address32](nil), sc.NewOption[sc.Sequence[primitives.Pre]](nil)
	}

	var resWithInfo primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]

	c.transactional.WithStorageLayer(
		func() (primitives.PostDispatchInfo, primitives.DispatchError) {
			resWithInfo = c.function.Dispatch(primitives.RawOriginFrom(maybeWho), c.function.Args())

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
	err := c.extra.PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	return resWithInfo, err
}

func (c checkedExtrinsic) Extra() primitives.SignedExtra {
	return c.extra
}

func (c checkedExtrinsic) Function() primitives.Call {
	return c.function
}

func (c checkedExtrinsic) Signed() sc.Option[primitives.Address32] {
	return c.signed
}

func (c checkedExtrinsic) Validate(validator primitives.UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	if c.signed.HasValue {
		id := c.signed.Value
		return c.extra.Validate(id, c.function, info, length)
	}

	valid, err := c.extra.ValidateUnsigned(c.function, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	unsignedValidation, err := validator.ValidateUnsigned(source, c.function)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	return valid.CombineWith(unsignedValidation), nil
}
