package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedInclusionFeeBytes, _ = hex.DecodeString(
		"010000000000000000000000000000000200000000000000000000000000000003000000000000000000000000000000",
	)
)

var (
	targetInclusionFee = InclusionFee{
		BaseFee:           sc.NewU128(1),
		LenFee:            sc.NewU128(2),
		AdjustedWeightFee: sc.NewU128(3),
	}
)

func Test_NewInclusionFee(t *testing.T) {
	result := NewInclusionFee(sc.NewU128(1), sc.NewU128(2), sc.NewU128(3))

	assert.Equal(t, targetInclusionFee, result)
}

func Test_InclusionFee_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetInclusionFee.Encode(buffer)

	assert.Equal(t, expectedInclusionFeeBytes, buffer.Bytes())
}

func Test_InclusionFee_Bytes(t *testing.T) {
	assert.Equal(t, expectedInclusionFeeBytes, targetInclusionFee.Bytes())
}

func Test_DecodeInclusionFee(t *testing.T) {
	result, err := DecodeInclusionFee(bytes.NewBuffer(expectedInclusionFeeBytes))
	assert.NoError(t, err)

	assert.Equal(t, targetInclusionFee, result)
}

func Test_InclusionFee_InclusionFee(t *testing.T) {
	result := targetInclusionFee.InclusionFee()

	assert.Equal(t, sc.NewU128(6), result)
}
