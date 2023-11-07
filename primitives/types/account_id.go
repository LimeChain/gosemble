package types

import (
	"bytes"
	"errors"
	"reflect"

	sc "github.com/LimeChain/goscale"
)

var (
	errorPubKeyNotSupported = errors.New("public key type not supported")
)

type ISigner interface {
	sc.Encodable
}

// AccountId It's an account ID (pubkey).
type AccountId struct {
	Ed25519Signer
	Sr25519Signer
	EcdsaSigner
}

func (a AccountId) Encode(buffer *bytes.Buffer) {
	a.Ed25519Signer.Encode(buffer)
	a.Sr25519Signer.Encode(buffer)
	a.EcdsaSigner.Encode(buffer)
}

func (a AccountId) Bytes() []byte {
	return sc.EncodedBytes(a)
}

func DecodeAccountId[T ISigner](buffer *bytes.Buffer) (AccountId, error) {
	switch reflect.Zero(reflect.TypeOf(*new(T))).Interface().(type) {
	case Ed25519Signer:
		pkEd25519, err := DecodeEd25519Signer(buffer)
		if err != nil {
			return AccountId{}, err
		}
		return AccountId{Ed25519Signer: pkEd25519}, nil
	case Sr25519Signer:
		pkSr25519, err := DecodeSr25519Signer(buffer)
		if err != nil {
			return AccountId{}, err
		}
		return AccountId{Sr25519Signer: pkSr25519}, nil
	case EcdsaSigner:
		pkEcdsa, err := DecodeEcdsaSigner(buffer)
		if err != nil {
			return AccountId{}, err
		}
		return AccountId{EcdsaSigner: pkEcdsa}, nil
	}
	return AccountId{}, errorPubKeyNotSupported
}
