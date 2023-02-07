package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
type CheckedExtrinsic struct {
	Version sc.U8

	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	Signed   sc.Option[AccountIdExtra]
	Function Call
}

func (ex CheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	// TODO:
}

func DecodeCheckedExtrinsic(buffer *bytes.Buffer) CheckedExtrinsic {
	// TODO:
	return CheckedExtrinsic{}
}

func (ex CheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(ex)
}

type AccountIdExtra struct {
	Address32
	Extra
}

func (ae AccountIdExtra) Encode(buffer *bytes.Buffer) {
	ae.Address32.Encode(buffer)
	ae.Extra.Encode(buffer)
}

func DecodeAccountIdExtra(buffer *bytes.Buffer) AccountIdExtra {
	// TODO:
	return AccountIdExtra{}
}

func (ae AccountIdExtra) Bytes() []byte {
	return sc.EncodedBytes(ae)
}

// Implementation for checked extrinsic.
func (xt CheckedExtrinsic) GetDispatchInfo() DispatchInfo {
	// TODO: return xt.Function.GetDispatchInfo()
	return DispatchInfo{
		Weight:  WeightFromRefTime(sc.U64(len(xt.Bytes()))),
		Class:   NormalDispatch,
		PaysFee: PaysYes,
	}
}

// info *DispatchInfoOfCall
// ApplyExtrinsicResultWithInfo<PostDispatchInfoOf<Self::Call>> { // <U: ValidateUnsigned<Call = Self::Call>>
func (xt CheckedExtrinsic) ApplyUnsignedValidator(info *DispatchInfo, length sc.Compact) (ok PostDispatchInfo, err DispatchErrorWithPostInfo) {
	var (
		maybeWho interface{}
		maybePre sc.Option[Pre]
	)

	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.Extra
		pre, err := extra.PreDispatch(&id, &xt.Function, info, length)
		if err != nil {
			return ok, err
		}
		maybeWho, maybePre = id, sc.NewOption[Pre](pre)
	} else {
		// Do any pre-flight stuff for a unsigned transaction.
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

		maybeWho, maybePre = nil, sc.NewOption[Pre](nil)
	}

	var postInfo PostDispatchInfo
	res, err2 := xt.Function.Dispatch(maybeWho) // RuntimeOrigin::from()
	_ = res
	if err2 != nil {
		// postInfo = err.PostInfo
	}
	// postInfo = res.Info

	dispatchResult := NewDispatchResult(err2)
	_, err = PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	return ok, err
}

// Do any pre-flight stuff for a signed transaction.
//
// Make sure to perform the same checks as in [`Self::validate`].
func (e Extra) PreDispatch(who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError) {
	ok, err = who.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = call.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = info.PreDispatch()
	if err != nil {
		return ok, err
	}

	ok, err = Length(length).PreDispatch()
	if err != nil {
		return ok, err
	}

	return ok, err
}

func PreDispatchUnsigned(call *Call, info *DispatchInfo, length sc.Compact) (ok Pre, err TransactionValidityError) {
	ok, err = call.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = info.PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	ok, err = Length(length).PreDispatchUnsigned()
	if err != nil {
		return ok, err
	}

	return ok, err
}

func PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) (ok Pre, err TransactionValidityError) {
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

		ok, err = Length(length).PostDispatch()
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

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (e Extra) Validate(who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ok ValidTransaction, err TransactionValidityError) {
	valid := DefaultValidTransaction()

	ok, err = who.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = call.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = info.Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = Length(length).Validate()
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	return valid, err
}

// Validate an unsigned transaction for the transaction queue.
//
// This function can be called frequently by the transaction queue
// to obtain transaction validity against current state.
// It should perform all checks that determine a valid unsigned transaction,
// and quickly eliminate ones that are stale or incorrect.
//
// Make sure to perform the same checks in `pre_dispatch_unsigned` function.
func ValidateUnsigned(_call *Call, _info *DispatchInfo, _length sc.Compact) (ok ValidTransaction, err TransactionValidityError) {
	ok = DefaultValidTransaction()
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
func UPreDispatch(call *Call) (ok sc.Empty, err TransactionValidityError) {
	_, err = UValidateUnsigned(NewTransactionSource(InBlock), call) // .map(|_| ()).map_err(Into::into)
	return ok, err
}

// / Information on a transaction's validity and, if valid, on how it relates to other transactions.
func UValidateUnsigned(source TransactionSource, call *Call) (ok ValidTransaction, err TransactionValidityError) {
	// TODO:
	return ok, err
}
