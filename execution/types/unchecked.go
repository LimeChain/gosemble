package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	// ExtrinsicFormatVersion is the current version of the [`UncheckedExtrinsic`] encoded format.
	//
	// This version needs to be bumped if the encoded representation changes.
	// It ensures that if the representation is changed and the format is not known,
	// the decoding fails.
	ExtrinsicFormatVersion = 4
	ExtrinsicBitSigned     = 0b1000_0000
	ExtrinsicUnmaskVersion = 0b0111_1111
)

type UncheckedExtrinsic interface {
	sc.Encodable

	Signature() sc.Option[primitives.ExtrinsicSignature]
	Function() primitives.Call
	Extra() primitives.SignedExtra

	IsSigned() sc.Bool
	Check(lookup primitives.AccountIdLookup) (sc.Option[primitives.Address32], primitives.TransactionValidityError)
}

type uncheckedExtrinsic struct {
	version sc.U8
	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	signature sc.Option[primitives.ExtrinsicSignature]
	function  primitives.Call
	extra     primitives.SignedExtra
	crypto    io.Crypto
}

// NewUncheckedExtrinsic returns a new instance of an unchecked extrinsic.
func NewUncheckedExtrinsic(version sc.U8, signature sc.Option[primitives.ExtrinsicSignature], function primitives.Call, extra primitives.SignedExtra) uncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:   version,
		signature: signature,
		function:  function,
		extra:     extra,
		crypto:    io.NewCrypto(),
	}
}

// NewUnsignedUncheckedExtrinsic returns a new instance of an unsigned extrinsic.
func NewUnsignedUncheckedExtrinsic(function primitives.Call) UncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:   ExtrinsicFormatVersion,
		signature: sc.NewOption[primitives.ExtrinsicSignature](nil),
		function:  function,
		crypto:    io.NewCrypto(),
	}
}
func (uxt uncheckedExtrinsic) Encode(buffer *bytes.Buffer) {
	tempBuffer := &bytes.Buffer{}

	uxt.version.Encode(tempBuffer)
	if uxt.signature.HasValue {
		uxt.signature.Value.Encode(tempBuffer)
	}
	uxt.function.Encode(tempBuffer)

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	encodedLen.Encode(buffer)
	buffer.Write(tempBuffer.Bytes())
}

func (uxt uncheckedExtrinsic) Bytes() []byte {
	return sc.EncodedBytes(uxt)
}

func (uxt uncheckedExtrinsic) Signature() sc.Option[primitives.ExtrinsicSignature] {
	return uxt.signature
}

func (uxt uncheckedExtrinsic) Function() primitives.Call {
	return uxt.function
}

func (uxt uncheckedExtrinsic) Extra() primitives.SignedExtra {
	return uxt.extra
}

func (uxt uncheckedExtrinsic) IsSigned() sc.Bool {
	return uxt.signature.HasValue
}

func (uxt uncheckedExtrinsic) Check(lookup primitives.AccountIdLookup) (sc.Option[primitives.Address32], primitives.TransactionValidityError) {
	if uxt.signature.HasValue {
		signer, signature, extra := uxt.signature.Value.Signer, uxt.signature.Value.Signature, uxt.signature.Value.Extra

		signedAddress, err := lookup.Lookup(signer)
		if err != nil {
			return sc.NewOption[primitives.Address32](nil), err
		}

		rawPayload, err := primitives.NewSignedPayload(uxt.function, extra)
		if err != nil {
			return sc.NewOption[primitives.Address32](nil), err
		}

		if !uxt.verify(signature, rawPayload.UsingEncoded(), signedAddress) {
			return sc.NewOption[primitives.Address32](nil), primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
		}

		return sc.NewOption[primitives.Address32](signedAddress), nil
	}

	return sc.NewOption[primitives.Address32](nil), nil
}

func (uxt uncheckedExtrinsic) verify(signature primitives.MultiSignature, msg sc.Sequence[sc.U8], signer primitives.Address32) bool {
	msgBytes := sc.SequenceU8ToBytes(msg)
	signerBytes := sc.FixedSequenceU8ToBytes(signer.FixedSequence)

	if signature.IsEd25519() {
		sigBytes := sc.FixedSequenceU8ToBytes(signature.AsEd25519().H512.FixedSequence)
		return uxt.crypto.Ed25519Verify(sigBytes, msgBytes, signerBytes)
	} else if signature.IsSr25519() {
		sigBytes := sc.FixedSequenceU8ToBytes(signature.AsSr25519().H512.FixedSequence)
		return uxt.crypto.Sr25519Verify(sigBytes, msgBytes, signerBytes)
	} else if signature.IsEcdsa() {
		// TODO:
		return true
		// let m = sp_io::hashing::blake2_256(msg.get());
		// match sp_io::crypto::secp256k1_ecdsa_recover_compressed(sig.as_ref(), &m) {
		// 	Ok(pubkey) =>
		// 		&sp_io::hashing::blake2_256(pubkey.as_ref()) ==
		// 			<dyn AsRef<[u8; 32]>>::as_ref(who),
		// 	_ => false,
		// }
	}
	log.Critical("invalid MultiSignature type in Verify")
	panic("unreachable")
}
