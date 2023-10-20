package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesWeightPerClass, _ = hex.DecodeString("0408010c10011418011c20")
)

var (
	targetWeightPerClass = WeightsPerClass{
		BaseExtrinsic: WeightFromParts(1, 2),
		MaxExtrinsic:  sc.NewOption[Weight](WeightFromParts(3, 4)),
		MaxTotal:      sc.NewOption[Weight](WeightFromParts(5, 6)),
		Reserved:      sc.NewOption[Weight](WeightFromParts(7, 8)),
	}
)

func Test_WeightPerClass_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetWeightPerClass.Encode(buffer)

	assert.Equal(t, expectBytesWeightPerClass, buffer.Bytes())
}

func Test_DecodeWeightPerClass(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesWeightPerClass)

	result := DecodeWeightsPerClass(buffer)

	assert.Equal(t, targetWeightPerClass, result)
}

func Test_WeightPerClass_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesWeightPerClass, targetWeightPerClass.Bytes())
}
