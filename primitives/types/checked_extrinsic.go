package types

import (
	sc "github.com/LimeChain/goscale"
)

// Definition of something that the external world might want to say; its
// existence implies that it has been checked and is good, particularly with
// regards to the signature.
type CheckedExtrinsic struct{}

// Implementation for checked extrinsic.
func (xt CheckedExtrinsic) GetDispatchInfo() DispatchInfo {
	return DispatchInfo{}
}

func (xt CheckedExtrinsic) ApplyUnsignedValidator(info *DispatchInfo, length sc.Compact) (ok PostDispatchInfo, err DispatchErrorWithPostInfo) {
	return ok, err
}

// The type that encodes information that can be passed from pre_dispatch to post-dispatch.
type Pre = sc.Empty
