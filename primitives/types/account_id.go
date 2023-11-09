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

type SignerAddress interface {
	sc.Encodable
	SignatureType() sc.U8
}

// AccountId It's an account ID (pubkey).
type AccountId[S SignerAddress] struct {
	signerAddress S
}

func New[S SignerAddress](signer S) AccountId[S] {
	return AccountId[S]{signerAddress: signer}
}

func (a AccountId[S]) Encode(buffer *bytes.Buffer) {
	a.signerAddress.Encode(buffer)
}

func (a AccountId[S]) Bytes() []byte {
	return sc.EncodedBytes(a)
}

func DecodeAccountId[S SignerAddress](buffer *bytes.Buffer) (AccountId[SignerAddress], error) {
	switch reflect.Zero(reflect.TypeOf(*new(S))).Interface().(type) {
	case Ed25519Signer:
		pkEd25519, err := DecodeEd25519Signer(buffer)
		if err != nil {
			return AccountId[SignerAddress]{}, err
		}
		return AccountId[SignerAddress]{signerAddress: pkEd25519}, nil
	case Sr25519Signer:
		pkSr25519, err := DecodeSr25519Signer(buffer)
		if err != nil {
			return AccountId[SignerAddress]{}, err
		}
		return AccountId[SignerAddress]{signerAddress: pkSr25519}, nil
	case EcdsaSigner:
		pkEcdsa, err := DecodeEcdsaSigner(buffer)
		if err != nil {
			return AccountId[SignerAddress]{}, err
		}
		return AccountId[SignerAddress]{signerAddress: pkEcdsa}, nil
	}
	return AccountId[SignerAddress]{}, nil
}
