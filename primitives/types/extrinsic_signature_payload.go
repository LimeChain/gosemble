package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/hashing"
)

// A payload that has been signed for an unchecked extrinsics.
//
// Note that the payload that we sign to produce unchecked extrinsic signature
// is going to be different than the `SignaturePayload` - so the thing the extrinsic
// actually contains.
type SignedPayload struct { // TODO: make it generic [C Call, E SignedExtension]
	// Address   MultiAddress
	// Signature MultiSignature

	Call Call

	Extra Extra
	// Weight

	AdditionalSigned
}

// type SignedExtension struct {
// 	// Mi: the indicator of the Polkadot module
// 	// Fi(m):  the function indicator of the module
// 	Call Call

// 	// E: the extra data
// 	Extra Extra

// 	AdditionalSigned
// }

type AdditionalSigned struct {
	// Rv: a UINT32 containing the specification version of 14.
	SpecVersion sc.U32 // RuntimeVersion

	// Fv: a UINT32 containing the format version of 2.
	FormatVersion sc.U32 // Version uint8

	// Hh(G): a 32-byte array containing the genesis hash.
	GenesisHash H256 // size 32

	// Hh(B): a 32-byte array containing the hash of the block which starts the mortality period, as described in
	BlockHash H256 // size 32

	TransactionVersion sc.U32

	BlockNumber BlockNumber
}

// Create new `SignedPayload`.
//
// This function may fail if `additional_signed` of `Extra` is not available.
func NewSignedPayload(call Call, extra Extra) (sp SignedPayload, err TransactionValidityError) {
	additionalSigned, err := extra.AdditionalSigned()
	if err != nil {
		return sp, err
	}

	sp = SignedPayload{
		Call:             call,
		Extra:            extra,
		AdditionalSigned: additionalSigned,
	}

	return sp, err
}

func (sp SignedPayload) Encode(buffer *bytes.Buffer) {
	// sp.Signer.Encode(buffer)
	// sp.Signature.Encode(buffer)
	sp.Call.Encode(buffer)
	sp.Extra.Encode(buffer)
	sp.SpecVersion.Encode(buffer)
	sp.FormatVersion.Encode(buffer)
	sp.GenesisHash.Encode(buffer)
	sp.BlockHash.Encode(buffer)
	sp.TransactionVersion.Encode(buffer)
	sp.BlockNumber.Encode(buffer)
}

func DecodeSignedPayload(buffer *bytes.Buffer) SignedPayload {
	sp := SignedPayload{}
	// sp.Signer.Encode(buffer)
	// sp.Signature.Encode(buffer)
	sp.Call = DecodeCall(buffer)
	sp.Extra = DecodeExtra(buffer)
	sp.SpecVersion = sc.DecodeU32(buffer)
	sp.FormatVersion = sc.DecodeU32(buffer)
	sp.GenesisHash = DecodeH256(buffer)
	sp.BlockHash = DecodeH256(buffer)
	sp.TransactionVersion = sc.DecodeU32(buffer)
	sp.BlockNumber = sc.DecodeU32(buffer)
	return sp
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
