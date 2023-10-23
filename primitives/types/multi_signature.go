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

func NewMultiSignatureEd25519(signature SignatureEd25519) MultiSignature {
	return MultiSignature{sc.NewVaryingData(MultiSignatureEd25519, signature)}
}

func NewMultiSignatureSr25519(signature SignatureSr25519) MultiSignature {
	return MultiSignature{sc.NewVaryingData(MultiSignatureSr25519, signature)}
}

func NewMultiSignatureEcdsa(signature SignatureEcdsa) MultiSignature {
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

func (s MultiSignature) AsEd25519() SignatureEd25519 {
	if s.IsEd25519() {
		return s.VaryingData[1].(SignatureEd25519)
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

func (s MultiSignature) AsSr25519() SignatureSr25519 {
	if s.IsSr25519() {
		return s.VaryingData[1].(SignatureSr25519)
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

func (s MultiSignature) AsEcdsa() SignatureEcdsa {
	if s.IsEcdsa() {
		return s.VaryingData[0].(SignatureEcdsa)
	} else {
		log.Critical("not a Ecdsa signature type")
	}

	panic("unreachable")
}

func DecodeMultiSignature(buffer *bytes.Buffer) MultiSignature {
	b := sc.DecodeU8(buffer)

	switch b {
	case MultiSignatureEd25519:
		return NewMultiSignatureEd25519(DecodeSignatureEd25519(buffer))
	case MultiSignatureSr25519:
		return NewMultiSignatureSr25519(DecodeSignatureSr25519(buffer))
	case MultiSignatureEcdsa:
		return NewMultiSignatureEcdsa(DecodeSignatureEcdsa(buffer))
	default:
		log.Critical("invalid MultiSignature type in Decode: " + string(b))
	}

	panic("unreachable")
}
