package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
//
// TODO: make it generic
// generic::CheckedExtrinsic<AccountId, RuntimeCall, SignedExtra>;
type CheckedExtrinsic struct {
	Version sc.U8

	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	Signed   sc.Option[AccountIdExtra]
	Function Call
}

func (xt CheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	xt.Version.Encode(buffer)
	xt.Signed.Encode(buffer)
	xt.Function.Encode(buffer)
}

func DecodeCheckedExtrinsic(buffer *bytes.Buffer) CheckedExtrinsic {
	xt := CheckedExtrinsic{}
	xt.Version = sc.DecodeU8(buffer)
	xt.Signed = sc.DecodeOptionWith(buffer, DecodeAccountIdExtra)
	xt.Function = DecodeCall(buffer)
	return xt
}

func (ex CheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(ex)
}

type AccountIdExtra struct {
	Address32
	SignedExtra
}

func (ae AccountIdExtra) Encode(buffer *bytes.Buffer) {
	ae.Address32.Encode(buffer)
	ae.SignedExtra.Encode(buffer)
}

func DecodeAccountIdExtra(buffer *bytes.Buffer) AccountIdExtra {
	ae := AccountIdExtra{}
	ae.Address32 = DecodeAddress32(buffer)
	ae.SignedExtra = DecodeExtra(buffer)
	return ae
}

func (ae AccountIdExtra) Bytes() []byte {
	return sc.EncodedBytes(ae)
}
