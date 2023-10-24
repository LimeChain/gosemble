package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectSessionKeyBytes, _ = hex.DecodeString("0c01020374657374")
)

var (
	key    = []byte{1, 2, 3}
	typeId = [4]byte{'t', 'e', 's', 't'}

	sessionKey = NewSessionKey(key, typeId)
)

func Test_NewSessionKey(t *testing.T) {
	expect := SessionKey{
		Key:    sc.BytesToSequenceU8(key),
		TypeId: sc.BytesToFixedSequenceU8(typeId[:]),
	}

	target := NewSessionKey(key, typeId)

	assert.Equal(t, expect, target)
}

func Test_SessionKey_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	sessionKey.Encode(buffer)

	assert.Equal(t, expectSessionKeyBytes, buffer.Bytes())
}

func Test_DecodeSessionKey(t *testing.T) {
	buffer := bytes.NewBuffer(expectSessionKeyBytes)

	result := DecodeSessionKey(buffer)

	assert.Equal(t, sessionKey, result)
}

func Test_SessionKey_Bytes(t *testing.T) {
	assert.Equal(t, expectSessionKeyBytes, sessionKey.Bytes())
}
