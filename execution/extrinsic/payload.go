package extrinsic

import (
	system "github.com/LimeChain/gosemble/frame/system/extensions"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Create new `SignedPayload`.
//
// This function may fail if `additional_signed` of `Extra` is not available.
func NewSignedPayload(call types.Call, extra types.SignedExtra) (ok types.SignedPayload, err types.TransactionValidityError) {
	additionalSigned, err := system.Extra(extra).AdditionalSigned()
	if err != nil {
		return ok, err
	}

	return types.SignedPayload{
		Call:             call,
		Extra:            extra,
		AdditionalSigned: additionalSigned,
	}, err
}
