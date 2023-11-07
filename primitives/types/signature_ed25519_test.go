package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesSignatureEd25519, _ = hex.DecodeString("86b2e08855e33d00e67a516d389224ed97e63de8e660a31649993eda5fa4297bd5e8c3faec06eb84918773b766a31d03a72df88d83a2f75ac82802fc7efd7c8d")
)

var (
	signatureEd25519 = NewSignatureEd25519(sc.BytesToSequenceU8(bytesSignatureEd25519)...)
)

func Test_NewSignatureEd25519(t *testing.T) {
	expect := SignatureEd25519{sc.BytesToFixedSequenceU8(bytesSignatureEd25519)}

	assert.Equal(t, expect, signatureEd25519)
}

func Test_SignatureEd25519_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := signatureEd25519.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, bytesSignatureEd25519, buffer.Bytes())
}

func Test_DecodeSignatureEd25519(t *testing.T) {
	buffer := bytes.NewBuffer(bytesSignatureEd25519)

	result, err := DecodeSignatureEd25519(buffer)
	assert.NoError(t, err)

	assert.Equal(t, signatureEd25519, result)
}

func Test_SignatureEd25519_Bytes(t *testing.T) {
	assert.Equal(t, bytesSignatureEd25519, signatureEd25519.Bytes())
}
