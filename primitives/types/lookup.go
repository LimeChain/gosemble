package types

import sc "github.com/LimeChain/goscale"

func Lookup(a MultiAddress) (AccountId[PublicKey], error) {
	address, err := lookupAddress(a)
	if err != nil {
		return AccountId[PublicKey]{}, err
	}
	if address.HasValue {
		return address.Value, nil
	}

	return AccountId[PublicKey]{}, NewTransactionValidityError(NewUnknownTransactionCannotLookup())
}

// LookupAddress Lookup an address to get an Id, if there's one there.
func lookupAddress(a MultiAddress) (sc.Option[AccountId[PublicKey]], error) {
	if a.IsAccountId() {
		accountId, err := a.AsAccountId()
		if err != nil {
			return sc.NewOption[AccountId[PublicKey]](nil), err
		}
		return sc.NewOption[AccountId[PublicKey]](accountId), nil
	}

	if a.IsAccountIndex() {
		index, err := a.AsAccountIndex()
		if err != nil {
			return sc.NewOption[AccountId[PublicKey]](nil), err
		}
		return lookupIndex(index), nil
	}

	return sc.NewOption[AccountId[PublicKey]](nil), nil
}

// LookupIndex Lookup an T::AccountIndex to get an Id, if there's one there.
func lookupIndex(index AccountIndex) sc.Option[AccountId[PublicKey]] {
	// TODO:
	return sc.NewOption[AccountId[PublicKey]](nil)
}
