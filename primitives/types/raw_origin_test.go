package types

import (
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	bytesAddress, _            = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	ed25519SignerZeroAddres, _ = NewEd25519PublicKey(sc.BytesToSequenceU8(bytesAddress)...)
	address                    = New[PublicKey](ed25519SignerZeroAddres)

	signedOrigin = NewRawOriginSigned(address)
)

func Test_NewRawOriginRoot(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginRoot)}

	assert.Equal(t, expect, NewRawOriginRoot())
}

func Test_NewRawOriginSigned(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginSigned, address)}

	assert.Equal(t, expect, NewRawOriginSigned(address))
}

func Test_NewRawOriginNone(t *testing.T) {
	expect := RawOrigin{sc.NewVaryingData(RawOriginNone)}

	assert.Equal(t, expect, NewRawOriginNone())
}

func Test_RawOriginFrom(t *testing.T) {
	option := sc.NewOption[AccountId[PublicKey]](address)

	result := RawOriginFrom(option)

	assert.Equal(t, signedOrigin, result)
}

func Test_RawOriginFrom_Empty(t *testing.T) {
	option := sc.NewOption[AccountId[PublicKey]](nil)

	result := RawOriginFrom(option)

	assert.Equal(t, NewRawOriginNone(), result)
}

func Test_RawOrigin_IsRootOrigin(t *testing.T) {
	assert.Equal(t, true, NewRawOriginRoot().IsRootOrigin())
	assert.Equal(t, false, signedOrigin.IsRootOrigin())
	assert.Equal(t, false, NewRawOriginNone().IsRootOrigin())
}

func Test_RawOrigin_IsSignedOrigin(t *testing.T) {
	assert.Equal(t, false, NewRawOriginRoot().IsSignedOrigin())
	assert.Equal(t, true, signedOrigin.IsSignedOrigin())
	assert.Equal(t, false, NewRawOriginNone().IsSignedOrigin())
}

func Test_RawOrigin_IsNoneOrigin(t *testing.T) {
	assert.Equal(t, false, NewRawOriginRoot().IsNoneOrigin())
	assert.Equal(t, false, signedOrigin.IsNoneOrigin())
	assert.Equal(t, true, NewRawOriginNone().IsNoneOrigin())
}

func Test_RawOrigin_AsSigned(t *testing.T) {
	result, err := signedOrigin.AsSigned()

	assert.NoError(t, err)
	assert.Equal(t, address, result)
}

func Test_RawOriginRoot_AsSigned_TypeError(t *testing.T) {
	address, err := NewRawOriginRoot().AsSigned()

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'RawOrigin' type", err.Error())
	assert.Equal(t, AccountId[PublicKey]{}, address)
}

func Test_RawOriginNone_AsSigned_TypeError(t *testing.T) {
	address, err := NewRawOriginNone().AsSigned()

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'RawOrigin' type", err.Error())
	assert.Equal(t, AccountId[PublicKey]{}, address)
}
