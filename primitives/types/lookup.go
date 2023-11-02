package types

import sc "github.com/LimeChain/goscale"

func Lookup(a MultiAddress) (Address32, TransactionValidityError) {
	address := lookupAddress(a)
	if address.HasValue {
		return address.Value, nil
	}

	return Address32{}, NewTransactionValidityError(NewUnknownTransactionCannotLookup())
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func lookupAddress(a MultiAddress) sc.Option[Address32] {
	if a.IsAccountId() {
		return sc.NewOption[Address32](a.AsAccountId().Address32)
	}

	if a.IsAddress32() {
		return sc.NewOption[Address32](a.AsAddress32())
	}

	if a.IsAccountIndex() {
		return lookupIndex(a.AsAccountIndex())
	}

	return sc.NewOption[Address32](nil)
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func lookupIndex(index AccountIndex) sc.Option[Address32] {
	// TODO:
	return sc.NewOption[Address32](nil)
}
