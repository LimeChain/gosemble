package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	MultiSignatureEd25519 sc.U8 = iota
	MultiSignatureSr25519
	MultiSignatureEcdsa
)

type MultiSignature struct {
	sc.VaryingData
}

func NewMultiSignatureEd25519(signature Ed25519) MultiSignature {
	return MultiSignature{sc.NewVaryingData(MultiSignatureEd25519, signature)}
}

func NewMultiSignatureSr25519(signature Sr25519) MultiSignature {
	return MultiSignature{sc.NewVaryingData(MultiSignatureSr25519, signature)}
}

func NewMultiSignatureEcdsa(signature Ecdsa) MultiSignature {
	return MultiSignature{sc.NewVaryingData(MultiSignatureEcdsa, signature)}
}

func (s MultiSignature) IsEd25519() sc.Bool {
	switch s.VaryingData[0] {
	case MultiSignatureEd25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEd25519() Ed25519 {
	if s.IsEd25519() {
		return s.VaryingData[1].(Ed25519)
	} else {
		log.Critical("not a Ed25519 signature type")
	}

	panic("unreachable")
}

func (s MultiSignature) IsSr25519() sc.Bool {
	switch s.VaryingData[0] {
	case MultiSignatureSr25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsSr25519() Sr25519 {
	if s.IsSr25519() {
		return s.VaryingData[1].(Sr25519)
	} else {
		log.Critical("not a Sr25519 signature type")
	}

	panic("unreachable")
}

func (s MultiSignature) IsEcdsa() sc.Bool {
	switch s.VaryingData[0] {
	case MultiSignatureEcdsa:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEcdsa() Ecdsa {
	if s.IsEcdsa() {
		return s.VaryingData[0].(Ecdsa)
	} else {
		log.Critical("not a Ecdsa signature type")
	}

	panic("unreachable")
}

func DecodeMultiSignature(buffer *bytes.Buffer) MultiSignature {
	b := sc.DecodeU8(buffer)

	switch b {
	case MultiSignatureEd25519:
		return NewMultiSignatureEd25519(DecodeEd25519(buffer))
	case MultiSignatureSr25519:
		return NewMultiSignatureSr25519(DecodeSr25519(buffer))
	case MultiSignatureEcdsa:
		return NewMultiSignatureEcdsa(DecodeEcdsa(buffer))
	default:
		log.Critical("invalid MultiSignature type in Decode: " + string(b))
	}

	panic("unreachable")
}
