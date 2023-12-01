package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IncRefStatus_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := IncRefStatusCreated.Encode(buffer)

	assert.Nil(t, err)
	assert.Equal(t, []byte{0}, buffer.Bytes())
}

func Test_IncRefStatus_Bytes(t *testing.T) {
	result := IncRefStatusExisted.Bytes()

	assert.Equal(t, []byte{1}, result)
}
