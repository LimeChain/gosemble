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

type PublicKeyType = sc.U8

const (
	PublicKeyEd25519 PublicKeyType = iota
	PublicKeySr25519
	PublicKeyEcdsa
)

type PublicKey interface {
	sc.Encodable
	SignatureType() sc.U8
}

// AccountId It's an account ID (pubkey).
type AccountId[T PublicKey] struct {
	publicKeyType T
}

func NewAccountId[T PublicKey](pkType T) AccountId[T] {
	return AccountId[T]{publicKeyType: pkType}
}

func (a AccountId[S]) Encode(buffer *bytes.Buffer) error {
	return a.publicKeyType.Encode(buffer)
}

func (a AccountId[S]) Bytes() []byte {
	return sc.EncodedBytes(a)
}

func DecodeAccountId[T PublicKey](buffer *bytes.Buffer) (AccountId[PublicKey], error) {
	switch reflect.Zero(reflect.TypeOf(*new(T))).Interface().(type) {
	case Ed25519PublicKey:
		pkEd25519, err := DecodeEd25519PublicKey(buffer)
		if err != nil {
			return AccountId[PublicKey]{}, err
		}
		return AccountId[PublicKey]{publicKeyType: pkEd25519}, nil
	case Sr25519PublicKey:
		pkSr25519, err := DecodeSr25519PublicKey(buffer)
		if err != nil {
			return AccountId[PublicKey]{}, err
		}
		return AccountId[PublicKey]{publicKeyType: pkSr25519}, nil
	case EcdsaPublicKey:
		pkEcdsa, err := DecodeEcdsaPublicKey(buffer)
		if err != nil {
			return AccountId[PublicKey]{}, err
		}
		return AccountId[PublicKey]{publicKeyType: pkEcdsa}, nil
	}
	return AccountId[PublicKey]{}, errorPubKeyNotSupported
}

func DecodeSequenceSr25519PublicKey(buffer *bytes.Buffer) (sc.Sequence[Sr25519PublicKey], error) {
	// decode length
	// for each decode Sr25519PublicKey
	// return the slice
	// TODO:
	DecodeSr25519PublicKey(buffer)
	return sc.Sequence[Sr25519PublicKey]{}, nil
}
