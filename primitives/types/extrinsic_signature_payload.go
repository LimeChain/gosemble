package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
)

// SignedPayload A payload that has been signed for an unchecked extrinsics.
//
// Note that the payload that we sign to produce unchecked extrinsic signature
// is going to be different than the `SignaturePayload` - so the thing the extrinsic
// actually contains.
//
// TODO: make it generic
// generic::SignedPayload<RuntimeCall, SignedExtra>;
type SignedPayload struct {
	Call  Call
	Extra SignedExtra
	AdditionalSigned
}

type AdditionalSigned = sc.VaryingData

func (sp SignedPayload) Encode(buffer *bytes.Buffer) {
	sp.Call.Encode(buffer)
	sp.Extra.Encode(buffer)
	sp.AdditionalSigned.Encode(buffer)
}

func (sp SignedPayload) Bytes() []byte {
	return sc.EncodedBytes(sp)
}

func (sp SignedPayload) UsingEncoded() sc.Sequence[sc.U8] {
	enc := sp.Bytes()

	if len(enc) > 256 {
		return sc.BytesToSequenceU8(hashing.Blake256(enc))
	} else {
		return sc.BytesToSequenceU8(enc)
	}
}
