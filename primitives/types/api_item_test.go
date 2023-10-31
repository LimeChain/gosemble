package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesApiItem, _ = hex.DecodeString("74657374696e673005000000")
)

var (
	apiName    = [8]byte{'t', 'e', 's', 't', 'i', 'n', 'g', '0'}
	apiVersion = sc.U32(5)

	apiItem = NewApiItem(apiName, apiVersion)
)

func Test_NewApiItem(t *testing.T) {
	expect := ApiItem{
		Name:    sc.BytesToFixedSequenceU8(apiName[:]),
		Version: apiVersion,
	}

	result := NewApiItem(apiName, apiVersion)

	assert.Equal(t, expect, result)
}

func Test_ApiItem_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	apiItem.Encode(buffer)

	assert.Equal(t, expectBytesApiItem, buffer.Bytes())
}

func Test_ApiItem_Decode(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesApiItem)

	result, err := DecodeApiItem(buffer)
	assert.NoError(t, err)

	assert.Equal(t, apiItem, result)
}

func Test_ApiItem_Bytes(t *testing.T) {

	assert.Equal(t, expectBytesApiItem, apiItem.Bytes())
}
