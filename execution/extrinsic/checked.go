package extrinsic

import (
	sc "github.com/LimeChain/goscale"

	system "github.com/LimeChain/gosemble/frame/system/extensions"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Checked types.CheckedExtrinsic

func (xt Checked) Validate(validator types.UnsignedValidator, source types.TransactionSource, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.SignedExtra
		ok, err = system.Extra(extra).Validate(&id, &xt.Function, info, length)
	} else {
		valid, err := system.Extra(types.SignedExtra{}).ValidateUnsigned(&xt.Function, info, length)
		if err != nil {
			return ok, err
		}

		unsignedValidation, err := validator.ValidateUnsigned(source, &xt.Function)
		if err != nil {
			return ok, err
		}

		ok = valid.CombineWith(unsignedValidation)
	}

	return ok, err
}

func (xt Checked) Apply(validator types.UnsignedValidator, info *types.DispatchInfo, length sc.Compact) (ok types.DispatchResultWithPostInfo[types.PostDispatchInfo], err types.TransactionValidityError) {
	var (
		maybeWho sc.Option[types.Address32]
		maybePre sc.Option[types.Pre]
	)

	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.SignedExtra
		pre, err := system.Extra(extra).PreDispatch(&id, &xt.Function, info, length)
		if err != nil {
			return ok, err
		}
		maybeWho, maybePre = sc.NewOption[types.Address32](id), sc.NewOption[types.Pre](pre)
	} else {
		// Do any pre-flight stuff for an unsigned transaction.
		//
		// Note this function by default delegates to `ValidateUnsigned`, so that
		// all checks performed for the transaction queue are also performed during
		// the dispatch phase (applying the extrinsic).
		//
		// If you ever override this function, you need to make sure to always
		// perform the same validation as in `ValidateUnsigned`.
		_, err := system.Extra{}.PreDispatchUnsigned(&xt.Function, info, length)
		if err != nil {
			return ok, err
		}

		_, err = validator.PreDispatch(&xt.Function)
		if err != nil {
			return ok, err
		}

		maybeWho, maybePre = sc.NewOption[types.Address32](nil), sc.NewOption[types.Pre](nil)
	}

	postDispatchInfo, resWithInfo := Dispatch(xt.Function, types.RawOriginFrom(maybeWho))

	var postInfo types.PostDispatchInfo
	if resWithInfo.HasError {
		postInfo = resWithInfo.Err.PostInfo
	}
	postInfo = postDispatchInfo

	dispatchResult := types.NewDispatchResult(resWithInfo.Err)
	_, err = system.Extra{}.PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	dispatchResultWithPostInfo := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	if resWithInfo.HasError {
		dispatchResultWithPostInfo.HasError = true
		dispatchResultWithPostInfo.Err = resWithInfo.Err
	} else {
		dispatchResultWithPostInfo.Ok = resWithInfo.Ok
	}

	return dispatchResultWithPostInfo, err
}
