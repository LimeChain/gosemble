package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Unchecked types.UncheckedExtrinsic

func (uxt Unchecked) Check(lookup primitives.AccountIdLookup) (types.CheckedExtrinsic, primitives.TransactionValidityError) {
	if uxt.Signature.HasValue {
		signer, signature, extra := uxt.Signature.Value.Signer, uxt.Signature.Value.Signature, uxt.Signature.Value.Extra

		signedAddress, err := lookup.Lookup(signer)
		if err != nil {
			return types.CheckedExtrinsic{}, err
		}

		rawPayload, err := NewSignedPayload(uxt.Function, extra)
		if err != nil {
			return types.CheckedExtrinsic{}, err
		}

		if !signature.Verify(rawPayload.UsingEncoded(), signedAddress) {
			err := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
			return types.CheckedExtrinsic{}, err
		}

		function, extra, _ := rawPayload.Call, rawPayload.Extra, rawPayload.AdditionalSigned

		return types.NewCheckedExtrinsic(sc.NewOption[primitives.Address32](signedAddress), function, extra), nil
	}

	return types.NewCheckedExtrinsic(sc.NewOption[primitives.Address32](nil), uxt.Function, uxt.Extra), nil
}
