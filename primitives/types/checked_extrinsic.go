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
	ae := AccountIdExtra{}
	ae.Address32 = DecodeAddress32(buffer)
	ae.Extra = DecodeExtra(buffer)
	return ae
}

func (ae AccountIdExtra) Bytes() []byte {
	return sc.EncodedBytes(ae)
}

// Implementation for checked extrinsic.
func (xt CheckedExtrinsic) GetDispatchInfo() DispatchInfo {
	// TODO:
	// return xt.Function.GetDispatchInfo()
	return DispatchInfo{
		Weight:  WeightFromRefTime(sc.U64(len(xt.Bytes()))),
		Class:   NormalDispatch,
		PaysFee: PaysYes,
	}
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
