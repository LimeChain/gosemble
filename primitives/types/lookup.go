package types

import sc "github.com/LimeChain/goscale"

type AccountIdLookup interface {
	Lookup(a MultiAddress) (Address32, TransactionValidityError)
}

//// AccountIdLookup A lookup implementation returning the `AccountId` from a `MultiAddress`.
//type accountIdLookup struct { // TODO: make it generic [AccountId, AccountIndex]
//	// TODO: PhantomData[(AccountId, AccountIndex)]
//}

//func DefaultAccountIdLookup() accountIdLookup {
//	return accountIdLookup{}
//}

// TODO: MultiAddress[AccountId, AccountIndex]
func Lookup(a MultiAddress) (Address32, TransactionValidityError) {
	address := lookupAddress(a)
	if address.HasValue {
		return address.Value, nil
	}

	return Address32{}, NewTransactionValidityError(NewUnknownTransactionCannotLookup())
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func lookupAddress(a MultiAddress) sc.Option[Address32] { // TODO: MultiAddress[AccountId, AccountIndex]
	if a.IsAccountId() {
		return sc.NewOption[Address32](a.AsAccountId().Address32)
	}

	if a.IsAddress32() {
		return sc.NewOption[Address32](a.AsAddress32())
	}

	if a.IsAccountIndex() {
		return sc.NewOption[Address32](lookupIndex(a.AsAccountIndex()))
	}

	return sc.NewOption[Address32](nil)
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func lookupIndex(index AccountIndex) sc.Option[Address32] {
	// TODO:
	return sc.NewOption[Address32](nil)
}
