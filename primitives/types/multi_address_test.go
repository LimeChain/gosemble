package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedaccountRawBytes, _ = hex.DecodeString("84010000000000000000010000000000000000000100000000000000000001000000")
)

var (
	addr20Bytes = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	addr32Bytes = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	addr33Bytes = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0}

	accountId    = AccountId{NewAddress32(sc.BytesToFixedSequenceU8(addr32Bytes)...)}
	accountIndex = sc.U32(2)
	accountRaw   = AccountRaw{sc.BytesToSequenceU8(addr33Bytes)}
	address32    = Address32{sc.BytesToFixedSequenceU8(addr32Bytes)}
	address20    = Address20{sc.BytesToFixedSequenceU8(addr20Bytes)}
)

func Test_AccountRaw_Encode(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	accountRaw.Encode(buffer)

	assert.Equal(t, expectedaccountRawBytes, buffer.Bytes())
}

func Test_DecodeAccountRaw(t *testing.T) {
	buffer := bytes.NewBuffer(expectedaccountRawBytes)

	result := DecodeAccountRaw(buffer)

	assert.Equal(t, accountRaw, result)
}

func Test_NewAddress32(t *testing.T) {
	expect := Address32{sc.NewFixedSequence(32, sc.BytesToSequenceU8(addr32Bytes)...)}

	assert.Equal(t, expect, NewAddress32(sc.BytesToSequenceU8(addr32Bytes)...))
}

func Test_NewAddress32_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "Address32 should be of size 32", func() {
		NewAddress32(sc.BytesToSequenceU8(addr20Bytes)...)
	})
}

func Test_DecodeAddress32(t *testing.T) {
	buffer := bytes.NewBuffer(addr32Bytes)

	result := DecodeAddress32(buffer)

	assert.Equal(t, address32, result)
}

func Test_NewAddress20(t *testing.T) {
	expect := Address20{sc.NewFixedSequence(20, sc.BytesToSequenceU8(addr20Bytes)...)}

	assert.Equal(t, expect, NewAddress20(sc.BytesToSequenceU8(addr20Bytes)...))
}

func Test_DecodeAddress20(t *testing.T) {
	buffer := bytes.NewBuffer(addr20Bytes)

	result := DecodeAddress20(buffer)

	assert.Equal(t, address20, result)
}

func Test_NewMultiAddressId(t *testing.T) {
	expect := MultiAddress{sc.NewVaryingData(sc.U8(0), accountId)}

	assert.Equal(t, expect, NewMultiAddressId(accountId))
}

func Test_NewMultiAddressIndex(t *testing.T) {
	expect := MultiAddress{sc.NewVaryingData(sc.U8(1), sc.ToCompact(accountIndex))}

	assert.Equal(t, expect, NewMultiAddressIndex(accountIndex))
}

func Test_NewMultiAddressRaw(t *testing.T) {
	expect := MultiAddress{sc.NewVaryingData(sc.U8(2), accountRaw)}

	assert.Equal(t, expect, NewMultiAddressRaw(accountRaw))
}

func Test_NewMultiAddress32(t *testing.T) {
	expect := MultiAddress{sc.NewVaryingData(sc.U8(3), address32)}

	assert.Equal(t, expect, NewMultiAddress32(address32))
}

func Test_NewMultiAddress20(t *testing.T) {
	expect := MultiAddress{sc.NewVaryingData(sc.U8(4), address20)}

	assert.Equal(t, expect, NewMultiAddress20(address20))
}

func Test_DecodeMultiAddress(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation MultiAddress
	}{
		{
			label:       "AccountId",
			input:       []byte{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			expectation: NewMultiAddressId(accountId),
		},
		{
			label:       "AccountIndex",
			input:       []byte{1, 8},
			expectation: NewMultiAddressIndex(accountIndex),
		},
		{
			label:       "AccountRaw",
			input:       []byte{2, 132, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
			expectation: NewMultiAddressRaw(accountRaw),
		},
		{
			label:       "Address32",
			input:       []byte{3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			expectation: NewMultiAddress32(address32),
		},
		{
			label:       "Address20",
			input:       []byte{4, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
			expectation: NewMultiAddress20(address20),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := bytes.NewBuffer(testExample.input)

			result := DecodeMultiAddress(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeMultiAddress_Panic(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{5, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0})

	assert.PanicsWithValue(t, "invalid MultiAddress type in Decode", func() {
		DecodeMultiAddress(buffer)
	})
}

func Test_IsAccountId(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       MultiAddress
		expectation bool
	}{
		{
			label:       "AccountId",
			input:       NewMultiAddressId(accountId),
			expectation: true,
		},
		{
			label:       "AddressIndex",
			input:       NewMultiAddressIndex(accountIndex),
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.IsAccountId())
		})
	}
}

func Test_AsAccountId(t *testing.T) {
	assert.Equal(t, accountId, NewMultiAddressId(accountId).AsAccountId())
}

func Test_AsAccountId_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "not an AccountId type", func() {
		NewMultiAddressIndex(accountIndex).AsAccountId()
	})
}

func Test_IsAccountIndex(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       MultiAddress
		expectation bool
	}{
		{
			label:       "AddressIndex",
			input:       NewMultiAddressIndex(accountIndex),
			expectation: true,
		},
		{
			label:       "AccountId",
			input:       NewMultiAddressId(accountId),
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.IsAccountIndex())
		})
	}
}

func Test_AsAccountIndex(t *testing.T) {
	assert.Equal(t, accountIndex, NewMultiAddressIndex(accountIndex).AsAccountIndex())
}

func Test_AsAccountIndex_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "not an AccountIndex type", func() {
		NewMultiAddressId(accountId).AsAccountIndex()
	})
}

func Test_IsRaw(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       MultiAddress
		expectation bool
	}{
		{
			label:       "AddressRaw",
			input:       NewMultiAddressRaw(accountRaw),
			expectation: true,
		},
		{
			label:       "AccountId",
			input:       NewMultiAddressId(accountId),
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.IsRaw())
		})
	}
}

func Test_AsRaw(t *testing.T) {
	assert.Equal(t, accountRaw, NewMultiAddressRaw(accountRaw).AsRaw())
}

func Test_AsRaw_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "not an AccountRaw type", func() {
		NewMultiAddressId(accountId).AsRaw()
	})
}

func Test_IsAddress32(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       MultiAddress
		expectation bool
	}{
		{
			label:       "Address32",
			input:       NewMultiAddress32(address32),
			expectation: true,
		},
		{
			label:       "AccountId",
			input:       NewMultiAddressId(accountId),
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.IsAddress32())
		})
	}
}

func Test_AsAddress32(t *testing.T) {
	assert.Equal(t, address32, NewMultiAddress32(address32).AsAddress32())
}

func Test_AsAddress32_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "not an Address32 type", func() {
		NewMultiAddressId(accountId).AsAddress32()
	})
}

func Test_IsAddress20(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       MultiAddress
		expectation bool
	}{
		{
			label:       "Address20",
			input:       NewMultiAddress20(address20),
			expectation: true,
		},
		{
			label:       "AccountId",
			input:       NewMultiAddressId(accountId),
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.IsAddress20())
		})
	}
}

func Test_AsAddress20(t *testing.T) {
	assert.Equal(t, address20, NewMultiAddress20(address20).AsAddress20())
}

func Test_AsAddress20_Panic(t *testing.T) {
	assert.PanicsWithValue(t, "not an Address20 type", func() {
		NewMultiAddressId(accountId).AsAddress20()
	})
}
