package extensions

import (
	"bytes"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	invalidTransactionStale, _  = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	invalidTransactionFuture, _ = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionFuture())
)

var (
	oneAddress = constants.OneAddressAccountId
)

func Test_CheckNonce_Encode(t *testing.T) {
	nonce := sc.U32(1)
	buffer := &bytes.Buffer{}

	target := setupCheckNonce()
	target.nonce = nonce

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, sc.ToCompact(nonce).Bytes(), buffer.Bytes())
}

func Test_CheckNonce_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := setupCheckNonce().Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, sc.ToCompact(sc.U32(0)).Bytes(), buffer.Bytes())
}

func Test_CheckNonce_Decode(t *testing.T) {
	nonce := sc.U32(1)
	buffer := bytes.NewBuffer(sc.ToCompact(nonce).Bytes())

	target := setupCheckNonce()

	target.Decode(buffer)

	assert.Equal(t, nonce, target.nonce)
}

func Test_CheckNonce_Bytes(t *testing.T) {
	nonce := sc.U32(1)
	target := setupCheckNonce()
	target.nonce = nonce

	result := target.Bytes()

	assert.Equal(t, sc.ToCompact(nonce).Bytes(), result)
}

func Test_CheckNonce_AdditionalSigned(t *testing.T) {
	target := setupCheckNonce()

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, primitives.AdditionalSigned{}, result)
}

func Test_CheckNonce_Validate_WithRequires_Success(t *testing.T) {
	nonce := sc.U32(1)
	accountInfo := primitives.AccountInfo{
		Nonce: 0,
	}
	expect := primitives.ValidTransaction{
		Priority: 0,
		Provides: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.BytesToSequenceU8(append(oneAddress.Bytes(), sc.ToCompact(nonce).Bytes()...)),
		},
		Requires: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.BytesToSequenceU8(append(oneAddress.Bytes(), sc.ToCompact(nonce-1).Bytes()...)),
		},
		Longevity: math.MaxUint64,
		Propagate: true,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)

	result, err := target.Validate(oneAddress, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
}

func Test_CheckNonce_Validate_NoRequires_Success(t *testing.T) {
	nonce := sc.U32(1)
	accountInfo := primitives.AccountInfo{
		Nonce: nonce,
	}
	expect := primitives.ValidTransaction{
		Priority: 0,
		Provides: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.BytesToSequenceU8(append(oneAddress.Bytes(), sc.ToCompact(nonce).Bytes()...)),
		},
		Requires:  sc.Sequence[sc.Sequence[sc.U8]]{},
		Longevity: math.MaxUint64,
		Propagate: true,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)

	result, err := target.Validate(oneAddress, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
}

func Test_CheckNonce_Validate_Fails(t *testing.T) {
	nonce := sc.U32(0)
	accountInfo := primitives.AccountInfo{
		Nonce: 1,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)

	result, err := target.Validate(oneAddress, nil, nil, sc.Compact{})

	assert.Equal(t, invalidTransactionStale, err)
	assert.Equal(t, primitives.ValidTransaction{}, result)
	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
}

func Test_CheckNonce_ValidateUnsigned(t *testing.T) {
	target := setupCheckNonce()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckNonce_PreDispatch_Success(t *testing.T) {
	nonce := sc.U32(1)
	accountInfo := primitives.AccountInfo{
		Nonce: 1,
	}
	expectAccountInfo := primitives.AccountInfo{
		Nonce: accountInfo.Nonce + 1,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)
	mockModule.On("StorageAccountSet", oneAddress, expectAccountInfo).Return()

	result, err := target.PreDispatch(oneAddress, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)

	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
	mockModule.AssertCalled(t, "StorageAccountSet", oneAddress, expectAccountInfo)
}

func Test_CheckNonce_PreDispatch_Fails_Stale(t *testing.T) {
	nonce := sc.U32(0)
	accountInfo := primitives.AccountInfo{
		Nonce: 1,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)

	result, err := target.PreDispatch(oneAddress, nil, nil, sc.Compact{})

	assert.Equal(t, invalidTransactionStale, err)
	assert.Equal(t, primitives.Pre{}, result)

	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
	mockModule.AssertNotCalled(t, "StorageAccountSet", oneAddress, mock.Anything)
}

func Test_CheckNonce_PreDispatch_Fails_Future(t *testing.T) {
	nonce := sc.U32(2)
	accountInfo := primitives.AccountInfo{
		Nonce: 1,
	}

	target := setupCheckNonce()
	target.nonce = nonce

	mockModule.On("StorageAccount", oneAddress).Return(accountInfo, nil)

	result, err := target.PreDispatch(oneAddress, nil, nil, sc.Compact{})

	assert.Equal(t, invalidTransactionFuture, err)
	assert.Equal(t, primitives.Pre{}, result)

	mockModule.AssertCalled(t, "StorageAccount", oneAddress)
	mockModule.AssertNotCalled(t, "StorageAccountSet", oneAddress, mock.Anything)
}

func Test_CheckNonce_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckNonce()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckNonce_PostDispatch(t *testing.T) {
	target := setupCheckNonce()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckNonce_Metadata(t *testing.T) {
	expectType := primitives.NewMetadataTypeWithPath(
		metadata.CheckNonce,
		"CheckNonce",
		sc.Sequence[sc.Str]{"frame_system", "extensions", "check_nonce", "CheckNonce"},
		primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32)),
	)
	expectSignedExtension := primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple)

	resultType, resultSignedExtension := setupCheckNonce().Metadata()

	assert.Equal(t, expectType, resultType)
	assert.Equal(t, expectSignedExtension, resultSignedExtension)
}

func setupCheckNonce() CheckNonce {
	mockModule = new(mocks.SystemModule)
	return CheckNonce{systemModule: mockModule}
}
