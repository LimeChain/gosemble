package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Unchecked types.UncheckedExtrinsic

func (uxt Unchecked) Check(lookup types.AccountIdLookup) (ok types.CheckedExtrinsic, err types.TransactionValidityError) {
	switch uxt.Signature.HasValue {
	case true:
		signer, signature, extra := uxt.Signature.Value.Signer, uxt.Signature.Value.Signature, uxt.Signature.Value.Extra

		signedAddress, err := lookup.Lookup(signer)
		if err != nil {
			return ok, err
		}

		rawPayload, err := NewSignedPayload(uxt.Function, extra)
		if err != nil {
			err = types.NewTransactionValidityError(types.NewUnknownTransaction(err))
			return ok, err
		}

		if !signature.Verify(rawPayload.UsingEncoded(), signedAddress) {
			err := types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadProofError))
			return ok, err
		}

		function, extra, _ := rawPayload.Call, rawPayload.Extra, rawPayload.AdditionalSigned

		ok = types.CheckedExtrinsic{
			Signed:   sc.NewOption[types.AccountIdExtra](types.AccountIdExtra{Address32: signedAddress, SignedExtra: extra}),
			Function: function,
		}
	case false:
		ok = types.CheckedExtrinsic{
			Signed:   sc.NewOption[types.AccountIdExtra](nil),
			Function: uxt.Function,
		}
	}

	return ok, err
}
