package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_IdentityFee_Mul(t *testing.T) {
	refTime := sc.U64(2)
	proofSize := sc.U64(3)
	expect := sc.NewU128(refTime)

	weight := WeightFromParts(refTime, proofSize)
	identityFee := IdentityFee{}
	result := identityFee.WeightToFee(weight)

	assert.Equal(t, expect, result)
}
