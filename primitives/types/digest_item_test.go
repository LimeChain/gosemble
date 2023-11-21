package types

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	consensusEngineId = sc.BytesToFixedSequenceU8([]byte{'t', 'e', 's', 't'})
	message           = sc.BytesToSequenceU8([]byte{'m', 'e', 's', 's', 'a', 'g', 'e'})
)

var (
	expectBytesDigestItemOther, _            = hex.DecodeString("001c6d657373616765")
	expectBytesDigestItemConsensusMessage, _ = hex.DecodeString("04746573741c6d657373616765")
	expectBytesDigestItemSeal, _             = hex.DecodeString("05746573741c6d657373616765")
	expectBytesDigestItemPreRuntime, _       = hex.DecodeString("06746573741c6d657373616765")
)

func Test_NewDigestItemOther(t *testing.T) {
	assert.Equal(t, DigestItem{sc.NewVaryingData(DigestItemOther, message)}, NewDigestItemOther(message))
}

func Test_NewDigestItemConsensusMessage(t *testing.T) {
	assert.Equal(t, DigestItem{sc.NewVaryingData(DigestItemConsensusMessage, consensusEngineId, message)}, NewDigestItemConsensusMessage(consensusEngineId, message))
}

func Test_NewDigestItemSeal(t *testing.T) {
	assert.Equal(t, DigestItem{sc.NewVaryingData(DigestItemSeal, consensusEngineId, message)}, NewDigestItemSeal(consensusEngineId, message))
}

func Test_NewDigestItemPreRuntime(t *testing.T) {
	assert.Equal(t, DigestItem{sc.NewVaryingData(DigestItemPreRuntime, consensusEngineId, message)}, NewDigestItemPreRuntime(consensusEngineId, message))
}

func Test_NewDigestItemRuntimeEnvironmentUpgrade(t *testing.T) {
	assert.Equal(t, DigestItem{sc.NewVaryingData(DigestItemRuntimeEnvironmentUpgraded)}, NewDigestItemRuntimeEnvironmentUpgrade())
}

func Test_DecodeDigestItem_Fails(t *testing.T) {
	buffer := &bytes.Buffer{}

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_Other(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesDigestItemOther)

	result, err := DecodeDigestItem(buffer)

	assert.NoError(t, err)
	assert.Equal(t, NewDigestItemOther(message), result)
}

func Test_decodeDigestItem_Other_Fails_Message(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_ConsensusMessage(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesDigestItemConsensusMessage)

	result, err := DecodeDigestItem(buffer)

	assert.NoError(t, err)
	assert.Equal(t, NewDigestItemConsensusMessage(consensusEngineId, message), result)
}

func Test_DecodeDigestItem_ConsensusMessage_Fails_Engine(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_ConsensusMessage_Fails_Message(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)
	buffer.Write(consensusEngineId.Bytes())

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_Seal(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesDigestItemSeal)

	result, err := DecodeDigestItem(buffer)

	assert.NoError(t, err)
	assert.Equal(t, NewDigestItemSeal(consensusEngineId, message), result)
}

func Test_DecodeDigestItem_Seal_Fails_Engine(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_Seal_Fails_Message(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)
	buffer.Write(consensusEngineId.Bytes())

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_PreRuntime(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesDigestItemPreRuntime)

	result, err := DecodeDigestItem(buffer)

	assert.NoError(t, err)
	assert.Equal(t, NewDigestItemPreRuntime(consensusEngineId, message), result)
}

func Test_DecodeDigestItem_PreRuntime_Fails_Engine(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_PreRuntime_Fails_Message(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)
	buffer.Write(consensusEngineId.Bytes())

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DecodeDigestItem_RuntimeEnvironmentUpgraded(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(8)

	result, err := DecodeDigestItem(buffer)

	assert.NoError(t, err)
	assert.Equal(t, NewDigestItemRuntimeEnvironmentUpgrade(), result)
}

func Test_DecodeDigestItem_InvalidType(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(9)

	result, err := DecodeDigestItem(buffer)

	assert.Equal(t, newTypeError("DigestItem"), err)
	assert.Equal(t, DigestItem{}, result)
}

func Test_DigestItem_IsSeal(t *testing.T) {
	assert.Equal(t, true, NewDigestItemSeal(consensusEngineId, message).IsSeal())
}

func Test_DigestItem_IsSeal_False(t *testing.T) {
	assert.Equal(t, false, NewDigestItemRuntimeEnvironmentUpgrade().IsSeal())
}

func Test_DigestItem_AsPreRuntime(t *testing.T) {
	expect := NewDigestPreRuntime(consensusEngineId, message)

	target := NewDigestItemPreRuntime(consensusEngineId, message)

	result, err := target.AsPreRuntime()

	assert.NoError(t, err)
	assert.Equal(t, expect, result)
}

func Test_DigestItem_AsPreRuntime_Fails_MismatchingType(t *testing.T) {
	target := NewDigestItemOther(message)

	result, err := target.AsPreRuntime()

	assert.Equal(t, newTypeError("DigestPreRuntime"), err)
	assert.Equal(t, DigestPreRuntime{}, result)
}
