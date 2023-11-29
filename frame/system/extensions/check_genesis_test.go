package extensions

import (
	"bytes"
	"errors"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	mockModule *mocks.SystemModule
)

func Test_CheckGenesis_AdditionalSigned(t *testing.T) {
	hash := primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}
	target := setupCheckGenesis()

	mockModule.On("StorageBlockHash", sc.U64(0)).Return(hash, nil)

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, sc.NewVaryingData(primitives.H256(hash)), result)
	mockModule.AssertCalled(t, "StorageBlockHash", sc.U64(0))
}

func Test_CheckGenesis_AdditionalSigned_Error(t *testing.T) {
	hash := primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}
	target := setupCheckGenesis()

	expectedErr := errors.New("error")

	mockModule.On("StorageBlockHash", sc.U64(0)).Return(hash, expectedErr)

	_, err := target.AdditionalSigned()
	assert.Equal(t, expectedErr, err)

	mockModule.AssertCalled(t, "StorageBlockHash", sc.U64(0))
}

func Test_CheckGenesis_Encode(t *testing.T) {
	target := setupCheckGenesis()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, 0, buffer.Len())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_CheckGenesis_Decode(t *testing.T) {
	target := setupCheckGenesis()
	value := []byte{1, 2, 3}
	buffer := bytes.NewBuffer(value)

	target.Decode(buffer)

	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, bytes.NewBuffer(value), buffer)
}

func Test_CheckGenesis_Bytes(t *testing.T) {
	target := setupCheckGenesis()

	result := target.Bytes()

	assert.Equal(t, []byte(nil), result)
}

func Test_CheckGenesis_Validate(t *testing.T) {
	target := setupCheckGenesis()

	result, err := target.Validate(constants.ZeroAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckGenesis_ValidateUnsigned(t *testing.T) {
	target := setupCheckGenesis()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckGenesis_PreDispatch(t *testing.T) {
	target := setupCheckGenesis()

	result, err := target.PreDispatch(constants.ZeroAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)
}

func Test_CheckGenesis_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckGenesis()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckGenesis_PostDispatch(t *testing.T) {
	target := setupCheckGenesis()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckGenesis_Metadata(t *testing.T) {
	expectType := primitives.NewMetadataTypeWithPath(
		metadata.CheckGenesis,
		"CheckGenesis",
		sc.Sequence[sc.Str]{"frame_system", "extensions", "check_genesis", "CheckGenesis"},
		primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
	)
	expectSignedExtension := primitives.NewMetadataSignedExtension("CheckGenesis", metadata.CheckGenesis, metadata.TypesH256)

	resultType, resultSignedExtension := setupCheckGenesis().Metadata()

	assert.Equal(t, expectType, resultType)
	assert.Equal(t, expectSignedExtension, resultSignedExtension)
}

func setupCheckGenesis() CheckGenesis {
	mockModule = new(mocks.SystemModule)
	extension, ok := NewCheckGenesis(mockModule).(*CheckGenesis)
	if !ok {
		panic("invalid type assert for *CheckGenesis")
	}
	return *extension
}
