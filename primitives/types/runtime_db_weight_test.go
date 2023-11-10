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
	expectBytesRuntimeDbWeight, _ = hex.DecodeString("02000000000000000300000000000000")
)

var (
	runtimeDbWeight = RuntimeDbWeight{
		Read:  2,
		Write: 3,
	}
)

func Test_RuntimeDbWeight_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := runtimeDbWeight.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesRuntimeDbWeight, buffer.Bytes())
}

func Test_RuntimeDbWeight_Bytes(t *testing.T) {
	result := runtimeDbWeight.Bytes()

	assert.Equal(t, expectBytesRuntimeDbWeight, result)
}

func Test_RuntimeDbWeight_Reads(t *testing.T) {
	readMultiplier := sc.U64(5)
	expect := WeightFromParts(runtimeDbWeight.Read*readMultiplier, 0)

	result := runtimeDbWeight.Reads(readMultiplier)

	assert.Equal(t, expect, result)
}

func Test_RuntimeDbWeight_Reads_Max(t *testing.T) {
	expect := WeightFromParts(math.MaxUint64, 0)

	result := runtimeDbWeight.Reads(math.MaxUint64)

	assert.Equal(t, expect, result)
}

func Test_RuntimeDbWeight_Writes(t *testing.T) {
	writeMultiplier := sc.U64(5)
	expect := WeightFromParts(runtimeDbWeight.Write*writeMultiplier, 0)

	result := runtimeDbWeight.Writes(writeMultiplier)

	assert.Equal(t, expect, result)
}

func Test_RuntimeDbWeight_Writes_Max(t *testing.T) {
	expect := WeightFromParts(math.MaxUint64, 0)

	result := runtimeDbWeight.Writes(math.MaxUint64)

	assert.Equal(t, expect, result)
}

func Test_RuntimeDbWeight_ReadWrites(t *testing.T) {
	readMultiplier := sc.U64(5)
	writeMultiplier := sc.U64(3)
	expect := WeightFromParts(readMultiplier*runtimeDbWeight.Read+writeMultiplier*runtimeDbWeight.Write, 0)

	result := runtimeDbWeight.ReadsWrites(readMultiplier, writeMultiplier)

	assert.Equal(t, expect, result)
}

func Test_RuntimeDbWeight_ReadWrites_Max(t *testing.T) {
	expect := WeightFromParts(math.MaxUint64, 0)

	result := runtimeDbWeight.ReadsWrites(math.MaxUint64, math.MaxUint64)

	assert.Equal(t, expect, result)
}
