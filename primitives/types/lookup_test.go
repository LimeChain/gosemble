package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	multiAddressId       = NewMultiAddressId(accId)
	multiAddressIndex    = NewMultiAddressIndex(accountIndex)
	expectedAccountId, _ = NewAccountId(sc.BytesToSequenceU8(signerAddressBytes)...)
	invalidMultiAddress  = MultiAddress{sc.NewVaryingData(sc.U8(5), sc.ToCompact(accountIndex))}

	expectedTransactionCannotLookupErr = NewTransactionValidityError(NewUnknownTransactionCannotLookup())
)

func Test_Lookup_AccountId(t *testing.T) {
	result, err := Lookup(multiAddressId)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountId, result)
}

func Test_Lookup_AccountIndex(t *testing.T) {
	result, err := Lookup(multiAddressIndex)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, AccountId{}, result)
}

func Test_Lookup_MultiAddress_NotValid(t *testing.T) {
	result, err := Lookup(invalidMultiAddress)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, AccountId{}, result)
}
