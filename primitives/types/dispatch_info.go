package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// A bundle of static information collected from the `#[pallet::weight]` attributes.
type DispatchInfo struct {
	// Weight of this transaction.
	Weight Weight

	// Class of this transaction.
	Class DispatchClass

	// Does this transaction pay fees.
	PaysFee Pays
}

func (di DispatchInfo) Encode(buffer *bytes.Buffer) {
	di.Weight.Encode(buffer)
	di.Class.Encode(buffer)
	di.PaysFee.Encode(buffer)
}

func DecodeDispatchInfo(buffer *bytes.Buffer) DispatchInfo {
	di := DispatchInfo{}
	di.Weight = DecodeWeight(buffer)
	di.Class = DecodeDispatchClass(buffer)
	di.PaysFee = DecodePays(buffer)
	return di
}

func (di DispatchInfo) Bytes() []byte {
	return sc.EncodedBytes(di)
}

func (di DispatchInfo) Validate() (ok ValidTransaction, err TransactionValidityError) {
	// TODO:
	return ok, err
}

func (di DispatchInfo) ValidateUnsigned() (ok ValidTransaction, err TransactionValidityError) {
	// TODO:
	return ok, err
}

func (di DispatchInfo) PreDispatch() (ok Pre, err TransactionValidityError) {
	// TODO:
	ok = Pre{}
	return ok, err
}

func (di DispatchInfo) PreDispatchUnsigned() (ok Pre, err TransactionValidityError) {
	// TODO:
	ok = Pre{}
	return ok, err
}

func (di DispatchInfo) PostDispatch() (ok Pre, err TransactionValidityError) {
	// TODO:
	ok = Pre{}
	return ok, err
}
