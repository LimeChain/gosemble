package types

import sc "github.com/LimeChain/goscale"

type Length sc.Compact

func (l Length) PreDispatch() (ok Pre, err TransactionValidityError) {
	// TODO
	ok = Pre{}
	return ok, err
}

func (l Length) PreDispatchUnsigned() (ok Pre, err TransactionValidityError) {
	// TODO
	ok = Pre{}
	return ok, err
}

func (l Length) PostDispatch() (ok Pre, err TransactionValidityError) {
	// TODO
	ok = Pre{}
	return ok, err
}

func (l Length) Validate() (ok ValidTransaction, err TransactionValidityError) {
	// TODO
	ok = DefaultValidTransaction()
	return ok, err
}
