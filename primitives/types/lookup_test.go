package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	signerAccountId = NewMultiAddressId(
		AccountId{
			Address32: NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...),
		},
	)
	signerAddress32 = NewMultiAddress32(
		Address32{
			sc.BytesToFixedSequenceU8(signerAddressBytes),
		},
	)
	signerAccountIndex  = NewMultiAddressIndex(accountIndex)
	expectedAccountId   = NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...)
	invalidMultiAddress = MultiAddress{sc.NewVaryingData(sc.U8(5), sc.ToCompact(accountIndex))}
)

func Test_Lookup_AccountId(t *testing.T) {
	result, err := Lookup(signerAccountId)
	assert.Nil(t, err)
	assert.Equal(t, result, expectedAccountId)
}

func Test_Lookup_Address32(t *testing.T) {
	result, err := Lookup(signerAddress32)
	assert.Nil(t, err)
	assert.Equal(t, result, expectedAccountId)
}

func Test_Lookup_AccountIndex(t *testing.T) {
	result, err := Lookup(signerAccountIndex)
	assert.Equal(t, err, NewTransactionValidityError(NewUnknownTransactionCannotLookup()))
	assert.Equal(t, result, Address32{})
}

func Test_Lookup_MultiAddress_NotValid(t *testing.T) {
	result, err := Lookup(invalidMultiAddress)
	assert.Equal(t, err, NewTransactionValidityError(NewUnknownTransactionCannotLookup()))
	assert.Equal(t, result, Address32{})
}
