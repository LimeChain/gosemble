package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesSignatureSr25519, _ = hex.DecodeString("86b2e08855e33d00e67a516d389224ed97e63de8e660a31649993eda5fa4297bd5e8c3faec06eb84918773b766a31d03a72df88d83a2f75ac82802fc7efd7c8d")
)

var (
	signatureSr25519 = NewSignatureSr25519(sc.BytesToSequenceU8(bytesSignatureSr25519)...)
)

func Test_Signature_NewSignatureSr25519(t *testing.T) {
	expect := SignatureSr25519{sc.BytesToFixedSequenceU8(bytesSignatureSr25519)}

	assert.Equal(t, expect, signatureSr25519)
}

func Test_SignatureSr25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := signatureSr25519.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, bytesSignatureSr25519, buffer.Bytes())
}

func Test_DecodeSignatureSr25519(t *testing.T) {
	buffer := bytes.NewBuffer(bytesSignatureSr25519)

	result, err := DecodeSignatureSr25519(buffer)
	assert.NoError(t, err)

	assert.Equal(t, signatureSr25519, result)
}

func Test_SignatureSr25519_Bytes(t *testing.T) {
	assert.Equal(t, bytesSignatureSr25519, signatureSr25519.Bytes())
}
