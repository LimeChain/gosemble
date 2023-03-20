package extrinsic

import (
	"github.com/LimeChain/gosemble/execution/types"
	system "github.com/LimeChain/gosemble/frame/system/extensions"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// Create new `SignedPayload`.
//
// This function may fail if `additional_signed` of `Extra` is not available.
func NewSignedPayload(call types.Call, extra primitives.SignedExtra) (ok primitives.SignedPayload, err primitives.TransactionValidityError) {
	additionalSigned, err := system.Extra(extra).AdditionalSigned()
	if err != nil {
		return ok, err
	}

	return primitives.SignedPayload{
		Call:             call,
		Extra:            extra,
		AdditionalSigned: additionalSigned,
	}, err
}
