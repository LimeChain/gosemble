package extensions

import (
	"bytes"
	"errors"
	"math"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	invalidTransactionAncientBirthBlock = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
)

func Test_CheckMortality_Encode(t *testing.T) {
	era := primitives.NewImmortalEra()
	buffer := &bytes.Buffer{}

	target := setupCheckMortality()
	target.era = era

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, era.Bytes(), buffer.Bytes())
}

func Test_CheckMortality_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := setupCheckMortality().Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, primitives.Era{}.Bytes(), buffer.Bytes())
}

func Test_CheckMortality_Decode(t *testing.T) {
	era := primitives.NewImmortalEra()
	buffer := bytes.NewBuffer(era.Bytes())

	target := setupCheckMortality()

	err := target.Decode(buffer)
	assert.Nil(t, err)

	assert.Equal(t, era, target.era)
}

func Test_CheckMortality_Bytes(t *testing.T) {
	era := primitives.NewImmortalEra()
	target := setupCheckMortality()
	target.era = era

	result := target.Bytes()

	assert.Equal(t, era.Bytes(), result)
}

func Test_CheckMortality_AdditionalSigned_Success(t *testing.T) {
	hash := primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}

	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockModule.On("StorageBlockHashExists", sc.U64(0)).Return(true)
	mockModule.On("StorageBlockHash", sc.U64(0)).Return(hash, nil)

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, sc.NewVaryingData(primitives.H256(hash)), result)

	mockModule.AssertCalled(t, "StorageBlockNumber")
	mockModule.AssertCalled(t, "StorageBlockHashExists", sc.U64(0))
	mockModule.AssertCalled(t, "StorageBlockHash", sc.U64(0))
}

func Test_CheckMortality_AdditionalSigned_Failed(t *testing.T) {
	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockModule.On("StorageBlockHashExists", sc.U64(0)).Return(false)

	result, err := target.AdditionalSigned()

	assert.Equal(t, invalidTransactionAncientBirthBlock, err)
	assert.Nil(t, result)

	mockModule.AssertCalled(t, "StorageBlockNumber")
	mockModule.AssertCalled(t, "StorageBlockHashExists", sc.U64(0))
}

func Test_CheckMortality_AdditionalSigned_StorageBlockNumberError(t *testing.T) {
	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	expectedErr := errors.New("error")
	mockModule.On("StorageBlockNumber").Return(blockNumber, expectedErr)

	_, err := target.AdditionalSigned()

	assert.Equal(t, expectedErr, err)

	mockModule.AssertCalled(t, "StorageBlockNumber")
}

func Test_CheckMortality_AdditionalSigned_StorageBlockHashError(t *testing.T) {
	hash := primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}

	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	expectedErr := errors.New("error")
	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockModule.On("StorageBlockHashExists", sc.U64(0)).Return(true)
	mockModule.On("StorageBlockHash", sc.U64(0)).Return(hash, expectedErr)

	_, err := target.AdditionalSigned()

	assert.Equal(t, expectedErr, err)
}

func Test_CheckMortality_AdditionalSigned_NewH256Error(t *testing.T) {
	hash := primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			[]byte{},
		)}
	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	expectedErr := errors.New("H256 should be of size 32")
	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockModule.On("StorageBlockHashExists", sc.U64(0)).Return(true)
	mockModule.On("StorageBlockHash", sc.U64(0)).Return(hash, nil)

	_, err := target.AdditionalSigned()

	assert.Equal(t, expectedErr, err)
}

func Test_CheckMortality_Validate_Success(t *testing.T) {
	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	expect := primitives.DefaultValidTransaction()
	expect.Longevity = math.MaxUint64 - 1

	blockNumber := sc.U64(1)

	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)

	result, err := target.Validate(constants.OneAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
}

func Test_CheckMortality_ValidateUnsigned(t *testing.T) {
	target := setupCheckMortality()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckMortality_PreDispatch(t *testing.T) {
	target := setupCheckMortality()
	target.era = primitives.NewImmortalEra()

	blockNumber := sc.U64(1)

	mockModule.On("StorageBlockNumber").Return(blockNumber, nil)

	result, err := target.PreDispatch(constants.OneAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)
}

func Test_CheckMortality_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckMortality()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckMortality_PostDispatch(t *testing.T) {
	target := setupCheckMortality()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckMortality_ModulePath(t *testing.T) {
	target := setupCheckMortality()

	expectedModulePath := "frame_system"
	actualModulePath := target.ModulePath()

	assert.Equal(t, expectedModulePath, actualModulePath)
}

func setupCheckMortality() CheckMortality {
	mockModule = new(mocks.SystemModule)
	extension, ok := NewCheckMortality(mockModule).(*CheckMortality)
	if !ok {
		panic("invalid type assert for *CheckMortality")
	}
	return *extension
}
