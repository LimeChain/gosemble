package types

import sc "github.com/LimeChain/goscale"

type MaxReserves struct {
	sc.U32
}

func (mr MaxReserves) Docs() string {
	return "The maximum number of named reserves that can exist on an account."
}
