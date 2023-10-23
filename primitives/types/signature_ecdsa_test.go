package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesSignatureEcdsa, _ = hex.DecodeString("86b2e08855e33d00e67a516d389224ed97e63de8e660a31649993eda5fa4297bd5e8c3faec06eb84918773b766a31d03a72df88d83a2f75ac82802fc7efd7c8dff")
)

var (
	signatureEcdsa = NewSignatureEcdsa(sc.BytesToSequenceU8(bytesSignatureEcdsa)...)
)

func Test_NewSignatureEcdsa(t *testing.T) {
	expect := SignatureEcdsa{sc.BytesToFixedSequenceU8(bytesSignatureEcdsa)}

	assert.Equal(t, expect, signatureEcdsa)
}

func Test_SignatureEcdsa_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	signatureEcdsa.Encode(buffer)

	assert.Equal(t, bytesSignatureEcdsa, buffer.Bytes())
}

func Test_DecodeSignatureEcdsa(t *testing.T) {
	buffer := bytes.NewBuffer(bytesSignatureEcdsa)

	result := DecodeSignatureEcdsa(buffer)

	assert.Equal(t, signatureEcdsa, result)
}

func Test_SignatureEcdsa_Bytes(t *testing.T) {
	assert.Equal(t, bytesSignatureEcdsa, signatureEcdsa.Bytes())
}
