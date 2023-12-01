package types

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pays_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := PaysYes.Encode(buffer)

	assert.Nil(t, err)
	assert.Equal(t, []byte{0}, buffer.Bytes())
}

func Test_Pays_Bytes(t *testing.T) {
	result := PaysNo.Bytes()

	assert.Equal(t, []byte{1}, result)
}

func Test_DecodePays_Yes(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result, err := DecodePays(buffer)
	assert.NoError(t, err)

	assert.Equal(t, PaysYes, result)
}

func Test_DecodePays_No(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodePays(buffer)
	assert.NoError(t, err)

	assert.Equal(t, PaysNo, result)
}

func Test_DecodePays_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	res, err := DecodePays(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'Pays' type", err.Error())
	assert.Equal(t, Pays(0), res)
}

func Test_DecodePays_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	res, err := DecodePays(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, Pays(0), res)
}
