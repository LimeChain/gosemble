package extrinsic

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	cts "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/types"
)

func ApplyUnsignedValidator(xt types.CheckedExtrinsic, info *types.DispatchInfo, length sc.Compact) (ok types.DispatchResultWithPostInfo[types.PostDispatchInfo], err types.TransactionValidityError) {
	var (
		maybeWho sc.Option[types.Address32]
		maybePre sc.Option[types.Pre]
	)

	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.Extra
		pre, err := extra.PreDispatch(&id, &xt.Function, info, length)
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
		_, err := PreDispatchUnsigned(&xt.Function, info, length)
		if err != nil {
			return ok, err
		}

		_, err = UPreDispatch(&xt.Function)
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
	_, err = PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	dispatchResultWithPostInfo := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	if resWithInfo.HasError {
		dispatchResultWithPostInfo.HasError = true
		dispatchResultWithPostInfo.Err = resWithInfo.Err
	} else {
		dispatchResultWithPostInfo.Ok = resWithInfo.Ok
	}

	return dispatchResultWithPostInfo, err
}

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) (ok types.PostDispatchInfo, err types.DispatchResultWithPostInfo[types.PostDispatchInfo]) {
	switch call.CallIndex.ModuleIndex {
	case system.ModuleIndex:
		switch call.CallIndex.FunctionIndex {
		case system.FunctionIndex:
			// TODO:
		}

	case cts.ModuleIndex:
		switch call.CallIndex.FunctionIndex {
		case cts.FunctionIndex:
			buffer := &bytes.Buffer{}
			buffer.Write(sc.SequenceU8ToBytes(call.Args))

			ts := sc.DecodeU64(buffer)
			timestamp.Set(ts)
		}
	}

	return ok, err
}

func PreDispatchUnsigned(call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	ok, err = call.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = info.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = types.Length(length).PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	return ok, err
}

func PostDispatch(pre sc.Option[types.Pre], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, length sc.Compact, result *types.DispatchResult) (ok types.Pre, err types.TransactionValidityError) {
	switch pre.HasValue {
	case true:
		// ok, err = pre.Value.PostDispatch()
		// if err != nil {
		// 	return ok, err
		// }

		ok, err = info.PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = postInfo.PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = types.Length(length).PostDispatch()
		if err != nil {
			return ok, err
		}

		ok, err = result.PostDispatch()
		if err != nil {
			return ok, err
		}

	case false:
		// sc.Empty
		// info
		// postInfo
		// length
		// result
	}

	return ok, err
}

// Validate the call right before dispatch.
//
// This method should be used to prevent transactions already in the pool
// (i.e. passing [`validate_unsigned`](Self::validate_unsigned)) from being included in blocks
// in case they became invalid since being added to the pool.
//
// By default it's a good idea to call [`validate_unsigned`](Self::validate_unsigned) from
// within this function again to make sure we never include an invalid transaction. Otherwise
// the implementation of the call or this method will need to provide proper validation to
// ensure that the transaction is valid.
//
// Changes made to storage *WILL* be persisted if the call returns `Ok`.
func UPreDispatch(call *types.Call) (ok sc.Empty, err types.TransactionValidityError) {
	_, err = UValidateUnsigned(types.NewTransactionSource(types.InBlock), call) // .map(|_| ()).map_err(Into::into)
	return ok, err
}

// / Information on a transaction's validity and, if valid, on how it relates to other transactions.
func UValidateUnsigned(source types.TransactionSource, call *types.Call) (ok types.ValidTransaction, err types.TransactionValidityError) {
	// TODO:
	return ok, err
}
