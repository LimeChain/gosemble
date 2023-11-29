package types

import sc "github.com/LimeChain/goscale"

type PublicKey interface {
	sc.Encodable
	SignatureType() sc.U8
}
