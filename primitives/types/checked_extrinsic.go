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

type UnsignedValidatorForChecked struct{}

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
func (v UnsignedValidatorForChecked) PreDispatch(call *Call) (ok sc.Empty, err TransactionValidityError) {
	_, err = v.ValidateUnsigned(NewTransactionSource(InBlock), call) // .map(|_| ()).map_err(Into::into)
	return ok, err
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (v UnsignedValidatorForChecked) ValidateUnsigned(source TransactionSource, call *Call) (ok ValidTransaction, err TransactionValidityError) {
	// TODO:
	// implement it for a specific pallet
	return ok, err
}
