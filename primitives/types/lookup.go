package types

func Lookup(a MultiAddress) (AccountId, error) {
	if !a.IsAccountId() {
		return AccountId{}, NewTransactionValidityError(NewUnknownTransactionCannotLookup())
	}

	return a.AsAccountId()
}
