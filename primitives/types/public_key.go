package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type PublicKeyType = sc.U8

const (
	PublicKeyEd25519 PublicKeyType = iota
	PublicKeySr25519
	PublicKeyEcdsa
)

// TODO: Extend for different types (ecdsa, ed25519, sr25519)
type PublicKey = sc.FixedSequence[sc.U8]

func DecodePublicKey(buffer *bytes.Buffer) (PublicKey, error) {
	return sc.DecodeFixedSequence[sc.U8](32, buffer)
}
