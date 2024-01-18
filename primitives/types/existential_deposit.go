package types

import sc "github.com/LimeChain/goscale"

type ExistentialDeposit struct {
	sc.U128
}

func (ed ExistentialDeposit) Docs() string {
	return "The minimum amount required to keep an account open. MUST BE GREATER THAN ZERO!"
}
