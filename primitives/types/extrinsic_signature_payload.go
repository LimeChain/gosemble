package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type SignedPayload interface {
	sc.Encodable

	AdditionalSigned() AdditionalSigned
	Call() Call
	Extra() SignedExtra
}

type AdditionalSigned = sc.VaryingData

// SignedPayload A payload that has been signed for an unchecked extrinsics.
//
// Note that the payload that we sign to produce unchecked extrinsic signature
// is going to be different than the `SignaturePayload` - so the thing the extrinsic
// actually contains.
//
// TODO: make it generic
// generic::SignedPayload<RuntimeCall, SignedExtra>;
type signedPayload struct {
	additionalSigned AdditionalSigned
	call             Call
	extra            SignedExtra
}

// NewSignedPayload creates a new `SignedPayload`.
// It may fail if `additional_signed` of `Extra` is not available.
func NewSignedPayload(call Call, extra SignedExtra) (SignedPayload, TransactionValidityError) {
	additionalSigned, err := extra.AdditionalSigned()
	if err != nil {
		return signedPayload{}, err
	}

	return signedPayload{
		call:             call,
		extra:            extra,
		additionalSigned: additionalSigned,
	}, nil
}

func (sp signedPayload) Encode(buffer *bytes.Buffer) {
	sp.call.Encode(buffer)
	sp.extra.Encode(buffer)
	sp.additionalSigned.Encode(buffer)
}

func (sp signedPayload) Bytes() []byte {
	return sc.EncodedBytes(sp)
}

func (sp signedPayload) AdditionalSigned() AdditionalSigned {
	return sp.additionalSigned
}

func (sp signedPayload) Call() Call {
	return sp.call
}

func (sp signedPayload) Extra() SignedExtra {
	return sp.extra
}
