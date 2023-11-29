package types

import sc "github.com/LimeChain/goscale"

type PublicKeyType = sc.U8

const (
	PublicKeyEd25519 PublicKeyType = iota
	PublicKeySr25519
	PublicKeyEcdsa
)
