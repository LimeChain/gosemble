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
	bytesAddress, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	address, _      = NewAccountId(sc.BytesToSequenceU8(bytesAddress)...)

	rootOrigin                 = NewRawOriginRoot()
	signedOrigin               = NewRawOriginSigned(address)
	signedOriginInvalidAddress = NewRawOriginSigned(AccountId{})
	noneOrigin                 = NewRawOriginNone()
)

func Test_NewRawOriginRoot(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginRoot)}

	assert.Equal(t, expect, rootOrigin)
}

func Test_NewRawOriginSigned(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginSigned, address)}

	assert.Equal(t, expect, NewRawOriginSigned(address))
}

func Test_NewRawOriginNone(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginNone)}

	assert.Equal(t, expect, noneOrigin)
}

func Test_RawOriginFrom(t *testing.T) {
	option := sc.NewOption[AccountId](address)

	result := RawOriginFrom(option)

	assert.Equal(t, signedOrigin, result)
}

func Test_RawOriginFrom_Empty(t *testing.T) {
	option := sc.NewOption[AccountId](nil)

	result := RawOriginFrom(option)

	assert.Equal(t, noneOrigin, result)
}

func Test_RawOrigin_IsRootOrigin(t *testing.T) {
	assert.Equal(t, true, rootOrigin.IsRootOrigin())
	assert.Equal(t, false, signedOrigin.IsRootOrigin())
	assert.Equal(t, false, noneOrigin.IsRootOrigin())
}

func Test_RawOrigin_IsSignedOrigin(t *testing.T) {
	assert.Equal(t, false, rootOrigin.IsSignedOrigin())
	assert.Equal(t, true, signedOrigin.IsSignedOrigin())
	assert.Equal(t, false, noneOrigin.IsSignedOrigin())
}

func Test_RawOrigin_IsNoneOrigin(t *testing.T) {
	assert.Equal(t, false, rootOrigin.IsNoneOrigin())
	assert.Equal(t, false, signedOrigin.IsNoneOrigin())
	assert.Equal(t, true, noneOrigin.IsNoneOrigin())
}

func Test_RawOrigin_AsSigned(t *testing.T) {
	result, err := signedOrigin.AsSigned()

	assert.NoError(t, err)
	assert.Equal(t, address, result)
}

func Test_RawOriginRoot_AsSigned_TypeError(t *testing.T) {
	address, err := rootOrigin.AsSigned()

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'RawOrigin' type", err.Error())
	assert.Equal(t, AccountId{}, address)
}

func Test_RawOriginNone_AsSigned_TypeError(t *testing.T) {
	address, err := noneOrigin.AsSigned()

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'RawOrigin' type", err.Error())
	assert.Equal(t, AccountId{}, address)
}

func Test_DecodeRawOrigin_Root(t *testing.T) {
	buffer := bytes.NewBuffer(rootOrigin.Bytes())

	result, err := DecodeRawOrigin(buffer)

	assert.NoError(t, err)
	assert.Equal(t, rootOrigin, result)
}

func Test_DecodeRawOrigin_Signed(t *testing.T) {
	buffer := bytes.NewBuffer(signedOrigin.Bytes())

	result, err := DecodeRawOrigin(buffer)

	assert.NoError(t, err)
	assert.Equal(t, signedOrigin, result)
}

func Test_DecodeRawOrigin_Signed_InvalidAddress(t *testing.T) {
	buffer := bytes.NewBuffer(signedOriginInvalidAddress.Bytes())

	result, err := DecodeRawOrigin(buffer)

	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, RawOrigin{}, result)
}

func Test_DecodeRawOrigin_None(t *testing.T) {
	buffer := bytes.NewBuffer(noneOrigin.Bytes())

	result, err := DecodeRawOrigin(buffer)

	assert.NoError(t, err)
	assert.Equal(t, noneOrigin, result)
}

func Test_DecodeRawOrigin_InvalidType(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{0x03})

	result, err := DecodeRawOrigin(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'RawOrigin' type", err.Error())
	assert.Equal(t, RawOrigin{}, result)
}
