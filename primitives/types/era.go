package types

import sc "github.com/LimeChain/goscale"

type ExtrinsicEra struct {
	IsImmortalEra sc.Bool
	IsMortalEra   sc.Bool
	AsMortalEra   MortalEra
}

type MortalEra struct {
	First  sc.U8
	Second sc.U8
}
