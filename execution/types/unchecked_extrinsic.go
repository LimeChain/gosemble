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

type PayloadInitializer = func(call primitives.Call, extra primitives.SignedExtra) (
	primitives.SignedPayload, primitives.TransactionValidityError,
)

type uncheckedExtrinsic struct {
	version sc.U8
	// The signature, address, number of extrinsics have come before from
	// the same signer and an era describing the longevity of this transaction,
	// if this is a signed extrinsic.
	signature         sc.Option[primitives.ExtrinsicSignature]
	function          primitives.Call
	extra             primitives.SignedExtra
	initializePayload PayloadInitializer
	crypto            io.Crypto
	hashing           io.Hashing
}

// NewUncheckedExtrinsic returns a new instance of an unchecked extrinsic.
func NewUncheckedExtrinsic(version sc.U8, signature sc.Option[primitives.ExtrinsicSignature], function primitives.Call, extra primitives.SignedExtra) primitives.UncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:           version,
		signature:         signature,
		function:          function,
		extra:             extra,
		initializePayload: primitives.NewSignedPayload,
		crypto:            io.NewCrypto(),
		hashing:           io.NewHashing(),
	}
}

// NewUnsignedUncheckedExtrinsic returns a new instance of an unsigned extrinsic.
func NewUnsignedUncheckedExtrinsic(function primitives.Call) primitives.UncheckedExtrinsic {
	return uncheckedExtrinsic{
		version:   ExtrinsicFormatVersion,
		signature: sc.NewOption[primitives.ExtrinsicSignature](nil),
		function:  function,
		crypto:    io.NewCrypto(),
	}
}

func (uxt uncheckedExtrinsic) Encode(buffer *bytes.Buffer) error {
	tempBuffer := &bytes.Buffer{}

	err := uxt.version.Encode(tempBuffer)
	if err != nil {
		return err
	}
	if uxt.signature.HasValue {
		err := uxt.signature.Value.Encode(tempBuffer)
		if err != nil {
			return err
		}
	}
	err = uxt.function.Encode(tempBuffer)
	if err != nil {
		return err
	}

	encodedLen := sc.ToCompact(uint64(tempBuffer.Len()))
	err = encodedLen.Encode(buffer)
	if err != nil {
		return err
	}
	_, err = buffer.Write(tempBuffer.Bytes())
	return err
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

func (uxt uncheckedExtrinsic) IsSigned() bool {
	return bool(uxt.signature.HasValue)
}

func (uxt uncheckedExtrinsic) Check() (primitives.CheckedExtrinsic, primitives.TransactionValidityError) {
	if uxt.signature.HasValue {
		signer, signature, extra := uxt.signature.Value.Signer, uxt.signature.Value.Signature, uxt.signature.Value.Extra

		signerAddress, err := primitives.Lookup(signer)
		if err != nil {
			return nil, err
		}

		rawPayload, err := uxt.initializePayload(uxt.function, extra)
		if err != nil {
			return nil, err
		}

		if !uxt.verify(signature, uxt.usingEncoded(rawPayload), signerAddress) {
			// https://github.com/LimeChain/gosemble/issues/271
			invalidTransactionBadProof, _ := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadProof())
			return nil, invalidTransactionBadProof
		}

		return NewCheckedExtrinsic(sc.NewOption[primitives.AccountId[primitives.SignerAddress]](signerAddress), uxt.function, extra), nil
	}

	return NewCheckedExtrinsic(sc.NewOption[primitives.AccountId[primitives.SignerAddress]](nil), uxt.function, uxt.extra), nil
}

func (uxt uncheckedExtrinsic) verify(signature primitives.MultiSignature, msg sc.Sequence[sc.U8], signer primitives.AccountId[primitives.SignerAddress]) bool {
	msgBytes := sc.SequenceU8ToBytes(msg)
	signerBytes := signer.Bytes()

	if signature.IsEd25519() {
		sigEd25519, err := signature.AsEd25519()
		if err != nil {
			log.Critical(err.Error())
		}
		sigBytes := sc.FixedSequenceU8ToBytes(sigEd25519.FixedSequence)
		return uxt.crypto.Ed25519Verify(sigBytes, msgBytes, signerBytes)
	} else if signature.IsSr25519() {

		sigSr25519, err := signature.AsSr25519()
		if err != nil {
			log.Critical(err.Error())
		}
		sigBytes := sc.FixedSequenceU8ToBytes(sigSr25519.FixedSequence)
		return uxt.crypto.Sr25519Verify(sigBytes, msgBytes, signerBytes)
	} else if signature.IsEcdsa() {
		return true
		// TODO:
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

func (uxt uncheckedExtrinsic) usingEncoded(sp primitives.SignedPayload) sc.Sequence[sc.U8] {
	enc := sp.Bytes()

	if len(enc) > 256 {
		hash := uxt.hashing.Blake256(enc)
		return sc.BytesToSequenceU8(hash)
	} else {
		return sc.BytesToSequenceU8(enc)
	}
}
