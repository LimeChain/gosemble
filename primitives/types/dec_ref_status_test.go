package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DecRefStatus_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}
	err := DecRefStatusExists.Encode(buffer)

	assert.Nil(t, err)
	assert.Equal(t, []byte{1}, buffer.Bytes())
}

func Test_DecRefStatus_Bytes(t *testing.T) {
	result := DecRefStatusReaped.Bytes()

	assert.Equal(t, []byte{0}, result)
}
