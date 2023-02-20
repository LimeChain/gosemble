package types

import sc "github.com/LimeChain/goscale"

type Length sc.Compact

func (l Length) GreaterThan(otherL Length) sc.Bool {
	return sc.U128(l).ToBigInt().Uint64() > sc.U128(l).ToBigInt().Uint64()
}

func (l Length) SaturatingAdd(otherL Length) Length {
	// TODO:
	return Length(sc.Compact{0})
}

func (l Length) Validate() (ok ValidTransaction, err TransactionValidityError) {
	// TODO
	ok = DefaultValidTransaction()
	return ok, err
}

func (l Length) ValidateUnsigned() (ok ValidTransaction, err TransactionValidityError) {
	// TODO
	ok = DefaultValidTransaction()
	return ok, err
}

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
