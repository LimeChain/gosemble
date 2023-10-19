package types

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedWeightBytes, _ = hex.DecodeString(
		"0408",
	)
)

var (
	targetWeight = Weight{
		RefTime:   1,
		ProofSize: 2,
	}
)

func Test_Weight_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetWeight.Encode(buffer)

	assert.Equal(t, expectedWeightBytes, buffer.Bytes())
}

func Test_DecodeWeight(t *testing.T) {
	buffer := bytes.NewBuffer(expectedWeightBytes)

	result := DecodeWeight(buffer)

	assert.Equal(t, targetWeight, result)
}

func Test_Weight_Bytes(t *testing.T) {
	assert.Equal(t, expectedWeightBytes, targetWeight.Bytes())
}

func Test_Weight_Add(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(1,2).Add(weight(3,4))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(3, 4)},
			expectation: WeightFromParts(4, 6),
		},
		{
			label:       "weight(1,1).Add(weight(MaxU64,MaxU64))",
			input:       [2]Weight{WeightFromParts(1, 1), WeightFromParts(math.MaxUint64, math.MaxUint64)},
			expectation: WeightFromParts(0, 0),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].Add(testExample.input[1]))
		})
	}
}

func Test_Weight_SaturatingAdd(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(1,2).SaturatingAdd(weight(3,4))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(3, 4)},
			expectation: WeightFromParts(4, 6),
		},
		{
			label:       "weight(1,1).SaturatingAdd(weight(MaxU64,MaxU64))",
			input:       [2]Weight{WeightFromParts(1, 1), WeightFromParts(math.MaxUint64, math.MaxUint64)},
			expectation: WeightFromParts(math.MaxUint64, math.MaxUint64),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].SaturatingAdd(testExample.input[1]))
		})
	}
}

func Test_Weight_CheckedAdd(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation sc.Option[Weight]
	}{
		{
			label:       "weight(0,1).CheckedAdd(weight(0,MaxU64))",
			input:       [2]Weight{WeightFromParts(0, 1), WeightFromParts(0, math.MaxUint64)},
			expectation: sc.NewOption[Weight](nil),
		},
		{
			label:       "weight(1,0).CheckedAdd(weight(MaxU64,0))",
			input:       [2]Weight{WeightFromParts(1, 0), WeightFromParts(math.MaxUint64, 0)},
			expectation: sc.NewOption[Weight](nil),
		},
		{
			label:       "weight(1,1).CheckedAdd(weight(MaxU64-1,MaxU64-1))",
			input:       [2]Weight{WeightFromParts(1, 1), WeightFromParts(math.MaxUint64-1, math.MaxUint64-1)},
			expectation: sc.NewOption[Weight](WeightFromParts(math.MaxUint64, math.MaxUint64)),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].CheckedAdd(testExample.input[1]))
		})
	}
}

func Test_Weight_Sub(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).Sub(weight(3,4))",
			input:       [2]Weight{WeightFromParts(4, 5), WeightFromParts(3, 4)},
			expectation: WeightFromParts(1, 1),
		},
		{
			label:       "weight(0,0).Sub(weight(1,1))",
			input:       [2]Weight{WeightFromParts(0, 0), WeightFromParts(1, 1)},
			expectation: WeightFromParts(math.MaxUint64, math.MaxUint64),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].Sub(testExample.input[1]))
		})
	}
}

func Test_Weight_SaturatingSub(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).SaturatingSub(weight(3,4))",
			input:       [2]Weight{WeightFromParts(4, 5), WeightFromParts(3, 4)},
			expectation: WeightFromParts(1, 1),
		},
		{
			label:       "weight(0,0).SaturatingSub(weight(1,MaxU64))",
			input:       [2]Weight{WeightFromParts(0, 0), WeightFromParts(1, math.MaxUint64)},
			expectation: WeightFromParts(0, 0),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].SaturatingSub(testExample.input[1]))
		})
	}
}

func Test_Weight_Mul(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).Mul(2)",
			input:       WeightFromParts(4, 5),
			expectation: WeightFromParts(8, 10),
		},
		{
			label:       "weight(MaxU64,0).Mul(2)",
			input:       WeightFromParts(math.MaxUint64, 0),
			expectation: WeightFromParts(math.MaxUint64-1, 0),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input.Mul(2))
		})
	}
}

func Test_Weight_SaturatingMul(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).SaturatingMul(2)",
			input:       WeightFromParts(4, 5),
			expectation: WeightFromParts(8, 10),
		},
		{
			label:       "weight(MaxU64,0).SaturatingMul(2)",
			input:       WeightFromParts(math.MaxUint64, 0),
			expectation: WeightFromParts(math.MaxUint64, 0),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input.SaturatingMul(2))
		})
	}
}

func Test_Weight_Min(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).Min(weight(7,1))",
			input:       [2]Weight{WeightFromParts(4, 5), WeightFromParts(7, 1)},
			expectation: WeightFromParts(4, 1),
		},
		{
			label:       "weight(2,3).Min(weight(5,6))",
			input:       [2]Weight{WeightFromParts(2, 3), WeightFromParts(5, 6)},
			expectation: WeightFromParts(2, 3),
		},
		{
			label:       "weight(7,8).Min(weight(1,2))",
			input:       [2]Weight{WeightFromParts(7, 8), WeightFromParts(1, 2)},
			expectation: WeightFromParts(1, 2),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].Min(testExample.input[1]))
		})
	}
}

func Test_Weight_Max(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation Weight
	}{
		{
			label:       "weight(4,5).Max(weight(7,1))",
			input:       [2]Weight{WeightFromParts(4, 5), WeightFromParts(7, 1)},
			expectation: WeightFromParts(7, 5),
		},
		{
			label:       "weight(2,3).Max(weight(5,6))",
			input:       [2]Weight{WeightFromParts(2, 3), WeightFromParts(5, 6)},
			expectation: WeightFromParts(5, 6),
		},
		{
			label:       "weight(7,8).Max(weight(1,2))",
			input:       [2]Weight{WeightFromParts(7, 8), WeightFromParts(1, 2)},
			expectation: WeightFromParts(7, 8),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].Max(testExample.input[1]))
		})
	}
}

func Test_Weight_AllGt(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation bool
	}{
		{
			label:       "weight(1,2).AllGt(weight(0,0))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(0, 0)},
			expectation: true,
		},
		{
			label:       "weight(1,2).AllGt(weight(0,0))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(1, 0)},
			expectation: false,
		},
		{
			label:       "weight(1,2).AllGt(weight(2,3))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(2, 3)},
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].AllGt(testExample.input[1]))
		})
	}
}

func Test_Weight_AnyGt(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       [2]Weight
		expectation bool
	}{
		{
			label:       "weight(1,2).AnyGt(weight(0,0))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(0, 0)},
			expectation: true,
		},
		{
			label:       "weight(1,2).AnyGt(weight(0,0))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(1, 0)},
			expectation: true,
		},
		{
			label:       "weight(1,2).AnyGt(weight(2,3))",
			input:       [2]Weight{WeightFromParts(1, 2), WeightFromParts(2, 3)},
			expectation: false,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {

			assert.Equal(t, testExample.expectation, testExample.input[0].AnyGt(testExample.input[1]))
		})
	}
}

func Test_Weight_WeightFromParts(t *testing.T) {
	assert.Equal(t, Weight{RefTime: 1, ProofSize: 2}, WeightFromParts(1, 2))
}

func Test_Weight_WeightZero(t *testing.T) {
	assert.Equal(t, Weight{RefTime: 0, ProofSize: 0}, WeightZero())
}

func Test_Weight_SaturatingAccrue(t *testing.T) {
	w1 := WeightFromParts(1, 2)
	w2 := WeightFromParts(3, 4)
	w3 := WeightFromParts(1, math.MaxUint64)
	w4 := WeightFromParts(math.MaxUint64, 1)

	w1.SaturatingAccrue(w2)
	w3.SaturatingAccrue(w4)

	assert.Equal(t, WeightFromParts(4, 6), w1)
	assert.Equal(t, WeightFromParts(3, 4), w2)
	assert.Equal(t, WeightFromParts(math.MaxUint64, math.MaxUint64), w3)
	assert.Equal(t, WeightFromParts(math.MaxUint64, 1), w4)
}

func Test_Weight_SaturatingReduce(t *testing.T) {
	w1 := WeightFromParts(3, 4)
	w2 := WeightFromParts(1, 2)
	w3 := WeightFromParts(0, 0)
	w4 := WeightFromParts(1, 1)

	w1.SaturatingReduce(w2)
	w3.SaturatingReduce(w4)

	assert.Equal(t, WeightFromParts(2, 2), w1)
	assert.Equal(t, WeightFromParts(1, 2), w2)
	assert.Equal(t, WeightFromParts(0, 0), w3)
	assert.Equal(t, WeightFromParts(1, 1), w4)
}
