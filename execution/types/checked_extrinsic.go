package types

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// CheckedExtrinsic is the definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
//
// TODO: make it generic
// generic::CheckedExtrinsic<AccountId, RuntimeCall, SignedExtra>;
type CheckedExtrinsic struct {
	// Who this purports to be from and the number of extrinsics have come before
	// from the same signer, if anyone (note this is not a signature).
	Signed   sc.Option[primitives.Address32]
	Function primitives.Call
	Extra    primitives.SignedExtra
}

func NewCheckedExtrinsic(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) CheckedExtrinsic {
	return CheckedExtrinsic{
		Signed:   signed,
		Function: function,
		Extra:    extra,
	}
}
