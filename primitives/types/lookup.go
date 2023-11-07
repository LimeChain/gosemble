package types

import sc "github.com/LimeChain/goscale"

func Lookup(a MultiAddress) (AccountId, TransactionValidityError) {
	// TODO https://github.com/LimeChain/gosemble/issues/271
	address, _ := lookupAddress(a)
	if address.HasValue {
		return address.Value, nil
	}

	unknownTransactionCannotLookup, _ := NewTransactionValidityError(NewUnknownTransactionCannotLookup())
	return AccountId{}, unknownTransactionCannotLookup
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func lookupAddress(a MultiAddress) (sc.Option[AccountId], error) {
	if a.IsAccountId() {
		accountId, err := a.AsAccountId()
		if err != nil {
			return sc.NewOption[AccountId](nil), err
		}
		return sc.NewOption[AccountId](accountId), nil
	}

	if a.IsAccountIndex() {
		index, err := a.AsAccountIndex()
		if err != nil {
			return sc.NewOption[AccountId](nil), err
		}
		return lookupIndex(index), nil
	}

	return sc.NewOption[AccountId](nil), nil
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func lookupIndex(index AccountIndex) sc.Option[AccountId] {
	// TODO:
	return sc.NewOption[AccountId](nil)
}
