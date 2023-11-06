package types

import (
	"bytes"
	"errors"
	"strconv"

	sc "github.com/LimeChain/goscale"
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

func (s MultiSignature) IsEd25519() bool {
	switch s.VaryingData[0] {
	case MultiSignatureEd25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEd25519() (SignatureEd25519, error) {
	if s.IsEd25519() {
		return s.VaryingData[1].(SignatureEd25519), nil
	} else {
		return SignatureEd25519{}, NewTypeError("SignatureEd25519")
	}
}

func (s MultiSignature) IsSr25519() bool {
	switch s.VaryingData[0] {
	case MultiSignatureSr25519:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsSr25519() (SignatureSr25519, error) {
	if s.IsSr25519() {
		return s.VaryingData[1].(SignatureSr25519), nil
	} else {
		return SignatureSr25519{}, NewTypeError("SignatureSr25519")
	}
}

func (s MultiSignature) IsEcdsa() bool {
	switch s.VaryingData[0] {
	case MultiSignatureEcdsa:
		return true
	default:
		return false
	}
}

func (s MultiSignature) AsEcdsa() (SignatureEcdsa, error) {
	if s.IsEcdsa() {
		return s.VaryingData[1].(SignatureEcdsa), nil
	} else {
		return SignatureEcdsa{}, NewTypeError("SignatureEcdsa")
	}
}

func DecodeMultiSignature(buffer *bytes.Buffer) (MultiSignature, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MultiSignature{}, err
	}

	switch b {
	case MultiSignatureEd25519:
		ed25519, err := DecodeSignatureEd25519(buffer)
		if err != nil {
			return MultiSignature{}, err
		}
		return NewMultiSignatureEd25519(ed25519), nil
	case MultiSignatureSr25519:
		sr25519, err := DecodeSignatureSr25519(buffer)
		if err != nil {
			return MultiSignature{}, err
		}
		return NewMultiSignatureSr25519(sr25519), nil
	case MultiSignatureEcdsa:
		ecdsa, err := DecodeSignatureEcdsa(buffer)
		if err != nil {
			return MultiSignature{}, err
		}
		return NewMultiSignatureEcdsa(ecdsa), nil
	default:
		return MultiSignature{}, errors.New("invalid MultiSignature type in Decode: " + strconv.Itoa(int(b)))
	}
}
