package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
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
	Signed   sc.Option[types.AccountIdExtra]
	Function Call
}

func (xt CheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	xt.Version.Encode(buffer)
	xt.Signed.Encode(buffer)
	xt.Function.Encode(buffer)
}

func (ex CheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(ex)
}
