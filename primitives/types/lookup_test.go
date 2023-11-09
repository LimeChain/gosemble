package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	multiAddressId = NewMultiAddressId(
		New[PublicKey](ed25519SignerOnesAddress),
	)
	multiAddressIndex   = NewMultiAddressIndex(accountIndex)
	expectedAccountId   = New[PublicKey](ed25519SignerOnesAddress)
	invalidMultiAddress = MultiAddress{sc.NewVaryingData(sc.U8(5), sc.ToCompact(accountIndex))}

	expectedTransactionCannotLookupErr, _ = NewTransactionValidityError(NewUnknownTransactionCannotLookup())
)

func Test_Lookup_AccountId(t *testing.T) {
	result, err := Lookup(multiAddressId)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountId, result)
}

func Test_Lookup_AccountIndex(t *testing.T) {
	result, err := Lookup(multiAddressIndex)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, AccountId[PublicKey]{}, result)
}

func Test_Lookup_MultiAddress_NotValid(t *testing.T) {
	result, err := Lookup(invalidMultiAddress)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, AccountId[PublicKey]{}, result)
}
