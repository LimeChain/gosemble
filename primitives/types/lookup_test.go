package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedAccountId, _ = NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...)

	multiAddressId = NewMultiAddressId(AccountId{Address32: expectedAccountId})
	multiAddress32 = NewMultiAddress32(
		Address32{
			sc.BytesToFixedSequenceU8(signerAddressBytes),
		},
	)
	multiAddressIndex   = NewMultiAddressIndex(accountIndex)
	invalidMultiAddress = MultiAddress{sc.NewVaryingData(sc.U8(5), sc.ToCompact(accountIndex))}

	expectedTransactionCannotLookupErr, _ = NewTransactionValidityError(NewUnknownTransactionCannotLookup())
)

func Test_Lookup_AccountId(t *testing.T) {
	result, err := Lookup(multiAddressId)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountId, result)
}

func Test_Lookup_Address32(t *testing.T) {
	result, err := Lookup(multiAddress32)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccountId, result)
}

func Test_Lookup_AccountIndex(t *testing.T) {
	result, err := Lookup(multiAddressIndex)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, Address32{}, result)
}

func Test_Lookup_MultiAddress_NotValid(t *testing.T) {
	result, err := Lookup(invalidMultiAddress)
	assert.Equal(t, expectedTransactionCannotLookupErr, err)
	assert.Equal(t, Address32{}, result)
}
