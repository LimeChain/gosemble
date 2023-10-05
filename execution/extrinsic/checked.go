package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Checked types.CheckedExtrinsic

func (xt Checked) Validate(validator UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	if xt.Signed.HasValue {
		id := xt.Signed.Value
		return xt.Extra.Validate(id, xt.Function, info, length)
	}

	valid, err := xt.Extra.ValidateUnsigned(xt.Function, info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	unsignedValidation, err := validator.ValidateUnsigned(source, xt.Function)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	return valid.CombineWith(unsignedValidation), nil
}

func (xt Checked) Apply(validator UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], primitives.TransactionValidityError) {
	var (
		maybeWho sc.Option[primitives.Address32]
		maybePre sc.Option[sc.Sequence[primitives.Pre]]
	)

	if xt.Signed.HasValue {
		id := xt.Signed.Value
		pre, err := xt.Extra.PreDispatch(id, xt.Function, info, length)
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
		err := xt.Extra.PreDispatchUnsigned(xt.Function, info, length)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		_, err = validator.PreDispatch(xt.Function)
		if err != nil {
			return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, err
		}

		maybeWho, maybePre = sc.NewOption[primitives.Address32](nil), sc.NewOption[sc.Sequence[primitives.Pre]](nil)
	}

	var resWithInfo primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]

	support.WithStorageLayer(
		func() (primitives.PostDispatchInfo, primitives.DispatchError) {
			resWithInfo = xt.Function.Dispatch(primitives.RawOriginFrom(maybeWho), xt.Function.Args())

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
	err := xt.Extra.PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	return resWithInfo, err
}
