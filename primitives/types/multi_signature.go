package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type MultiSignature sc.VaryingData

func NewMultiSignature(value sc.Encodable) MultiSignature {
	switch value.(type) {
	case Ed25519, Sr25519, Ecdsa:
		return MultiSignature(sc.NewVaryingData(value))
	default:
		panic("invalid Signature type")
	}
}

func (s MultiSignature) IsEd25519() sc.Bool {
	switch s[0].(type) {
	case Ed25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEd25519() Ed25519 {
	if s.IsEd25519() {
		return s[0].(Ed25519)
	} else {
		panic("not a Ed25519 signature type")
	}
}

func (s MultiSignature) IsSr25519() sc.Bool {
	switch s[0].(type) {
	case Sr25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsSr25519() Sr25519 {
	if s.IsSr25519() {
		return s[0].(Sr25519)
	} else {
		panic("not a Sr25519 signature type")
	}
}

func (s MultiSignature) IsEcdsa() sc.Bool {
	switch s[0].(type) {
	case Ecdsa:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEcdsa() Ecdsa {
	if s.IsEcdsa() {
		return s[0].(Ecdsa)
	} else {
		panic("not a Ecdsa signature type")
	}
}

func (s MultiSignature) Encode(buffer *bytes.Buffer) {
	if s.IsEd25519() {
		sc.U8(0).Encode(buffer)
		s.AsEd25519().Encode(buffer)
	} else if s.IsSr25519() {
		sc.U8(1).Encode(buffer)
		s.AsSr25519().Encode(buffer)
	} else if s.IsEcdsa() {
		sc.U8(2).Encode(buffer)
		s.AsEcdsa().Encode(buffer)
	} else {
		panic("invalid MultiSignature type in Encode")
	}
}

func DecodeMultiSignature(buffer *bytes.Buffer) MultiSignature {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return MultiSignature{DecodeEd25519(buffer)}
	case 1:
		return MultiSignature{DecodeSr25519(buffer)}
	case 2:
		return MultiSignature{DecodeEcdsa(buffer)}
	default:
		panic("invalid MultiSignature type in Decode: " + string(b))
	}
}

func (s MultiSignature) Verify(msg sc.Sequence[sc.U8], signer Address32) sc.Bool {
	if s.IsEd25519() {
		return s.AsEd25519().Verify(msg, signer)
	} else if s.IsSr25519() {
		return s.AsSr25519().Verify(msg, signer)
	} else if s.IsEcdsa() {
		// TODO:
		return true
		// let m = sp_io::hashing::blake2_256(msg.get());
		// match sp_io::crypto::secp256k1_ecdsa_recover_compressed(sig.as_ref(), &m) {
		// 	Ok(pubkey) =>
		// 		&sp_io::hashing::blake2_256(pubkey.as_ref()) ==
		// 			<dyn AsRef<[u8; 32]>>::as_ref(who),
		// 	_ => false,
		// }
	} else {
		panic("invalid MultiSignature type in Verify")
	}
}
