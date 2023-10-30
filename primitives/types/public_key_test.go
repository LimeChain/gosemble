package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_DecodePublicKey(t *testing.T) {
	bytesPublicKey, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	expect := sc.BytesToFixedSequenceU8(bytesPublicKey)

	buffer := bytes.NewBuffer(bytesPublicKey)
	result := DecodePublicKey(buffer)

	assert.Equal(t, expect, result)
}