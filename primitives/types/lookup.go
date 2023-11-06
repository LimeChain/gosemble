package types

import sc "github.com/LimeChain/goscale"

func Lookup(a MultiAddress) (Address32, TransactionValidityError) {
	// TODO https://github.com/LimeChain/gosemble/issues/271
	address, _ := lookupAddress(a)
	if address.HasValue {
		return address.Value, nil
	}

	unknownTransactionCannotLookup, _ := NewTransactionValidityError(NewUnknownTransactionCannotLookup())
	return Address32{}, unknownTransactionCannotLookup
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func lookupAddress(a MultiAddress) (sc.Option[Address32], error) {
	if a.IsAccountId() {
		accountId, err := a.AsAccountId()
		if err != nil {
			return sc.NewOption[Address32](nil), err
		}
		return sc.NewOption[Address32](accountId.Address32), nil
	}

	if a.IsAddress32() {
		address32, err := a.AsAddress32()
		if err != nil {
			return sc.NewOption[Address32](nil), err
		}
		return sc.NewOption[Address32](address32), nil
	}

	if a.IsAccountIndex() {
		index, err := a.AsAccountIndex()
		if err != nil {
			return sc.NewOption[Address32](nil), err
		}
		return lookupIndex(index), nil
	}

	return sc.NewOption[Address32](nil), nil
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func lookupIndex(index AccountIndex) sc.Option[Address32] {
	// TODO:
	return sc.NewOption[Address32](nil)
}
