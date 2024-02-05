package types

import sc "github.com/LimeChain/goscale"

type MaxLocks struct {
	sc.U32
}

func (ml MaxLocks) Docs() string {
	return "The maximum number of locks that should exist on an account.  Not strictly enforced, but used for weight estimation."
}
