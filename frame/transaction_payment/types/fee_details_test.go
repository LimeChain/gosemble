package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedFeeDetailsBytes, _ = hex.DecodeString(
		"01010000000000000000000000000000000200000000000000000000000000000003000000000000000000000000000000",
	)
)

var (
	targetFeeDetails = FeeDetails{
		InclusionFee: sc.NewOption[InclusionFee](InclusionFee{
			BaseFee:           sc.NewU128(1),
			LenFee:            sc.NewU128(2),
			AdjustedWeightFee: sc.NewU128(3),
		}),
		Tip: sc.NewU128(4),
	}
)

func Test_FeeDetails_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetFeeDetails.Encode(buffer)

	assert.Equal(t, expectedFeeDetailsBytes, buffer.Bytes())

}

func Test_FeeDetails_Bytes(t *testing.T) {
	assert.Equal(t, expectedFeeDetailsBytes, targetFeeDetails.Bytes())
}

func Test_FeeDetails_DecodeFeeDetails(t *testing.T) {
	result, err := DecodeFeeDetails(bytes.NewBuffer(expectedFeeDetailsBytes))
	assert.NoError(t, err)

	expectedFeeDetails := FeeDetails{
		InclusionFee: sc.NewOption[InclusionFee](InclusionFee{
			BaseFee:           sc.NewU128(1),
			LenFee:            sc.NewU128(2),
			AdjustedWeightFee: sc.NewU128(3),
		}),
	}

	assert.Equal(t, expectedFeeDetails, result)
}

func Test_FeeDetails_FinalFee_WithInclusionFee(t *testing.T) {
	result := targetFeeDetails.FinalFee()

	assert.Equal(t, sc.NewU128(10), result)
}

func Test_FeeDetails_FinalFee_WithoutInclusionFee(t *testing.T) {
	targetFeeDetails := FeeDetails{
		InclusionFee: sc.NewOption[InclusionFee](nil),
		Tip:          sc.NewU128(4),
	}

	result := targetFeeDetails.FinalFee()

	assert.Equal(t, sc.NewU128(4), result)
}
