package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedAccountDataBytes, _ = hex.DecodeString(
		"01000000000000000000000000000000020000000000000000000000000000000300000000000000000000000000000004000000000000000000000000000000",
	)
)

var (
	targetAccountData = AccountData{
		Free:       sc.NewU128(1),
		Reserved:   sc.NewU128(2),
		MiscFrozen: sc.NewU128(3),
		FeeFrozen:  sc.NewU128(4),
	}
)

func Test_AccountData_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetAccountData.Encode(buffer)

	assert.Equal(t, expectedAccountDataBytes, buffer.Bytes())
}

func Test_AccountData_Bytes(t *testing.T) {
	assert.Equal(t, expectedAccountDataBytes, targetAccountData.Bytes())
}

func Test_DecodeAccountData(t *testing.T) {
	buffer := bytes.NewBuffer(expectedAccountDataBytes)

	result := DecodeAccountData(buffer)

	assert.Equal(t, targetAccountData, result)
}

func Test_Total(t *testing.T) {
	assert.Equal(t, sc.NewU128(3), targetAccountData.Total())
}
