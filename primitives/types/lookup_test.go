package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	multiAddressId = NewMultiAddressId(
		AccountId{
			Address32: NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...),
		},
	)
	multiAddress32 = NewMultiAddress32(
		Address32{
			sc.BytesToFixedSequenceU8(signerAddressBytes),
		},
	)
	multiAddressIndex   = NewMultiAddressIndex(accountIndex)
	expectedAccountId   = NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...)
	invalidMultiAddress = MultiAddress{sc.NewVaryingData(sc.U8(5), sc.ToCompact(accountIndex))}
)

func Test_Lookup_AccountId(t *testing.T) {
	result, err := Lookup(multiAddressId)
	assert.Nil(t, err)
	assert.Equal(t, result, expectedAccountId)
}

func Test_Lookup_Address32(t *testing.T) {
	result, err := Lookup(multiAddress32)
	assert.Nil(t, err)
	assert.Equal(t, result, expectedAccountId)
}

func Test_Lookup_AccountIndex(t *testing.T) {
	result, err := Lookup(multiAddressIndex)
	assert.Equal(t, err, NewTransactionValidityError(NewUnknownTransactionCannotLookup()))
	assert.Equal(t, result, Address32{})
}

func Test_Lookup_MultiAddress_NotValid(t *testing.T) {
	result, err := Lookup(invalidMultiAddress)
	assert.Equal(t, err, NewTransactionValidityError(NewUnknownTransactionCannotLookup()))
	assert.Equal(t, result, Address32{})
}
