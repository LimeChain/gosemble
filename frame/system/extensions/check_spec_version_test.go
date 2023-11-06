package extensions

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_CheckSpecVersion_AdditionalSigned(t *testing.T) {
	version := primitives.RuntimeVersion{
		SpecVersion: 4,
	}
	target := setupCheckSpecVersion()

	mockModule.On("Version").Return(version)

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, sc.NewVaryingData(version.SpecVersion), result)
	mockModule.AssertCalled(t, "Version")
}

func Test_CheckSpecVersion_Encode(t *testing.T) {
	target := setupCheckSpecVersion()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, 0, buffer.Len())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_CheckSpecVersion_Decode(t *testing.T) {
	target := setupCheckSpecVersion()
	value := []byte{1, 2, 3}
	buffer := bytes.NewBuffer(value)

	target.Decode(buffer)

	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, bytes.NewBuffer(value), buffer)
}

func Test_CheckSpecVersion_Bytes(t *testing.T) {
	target := setupCheckSpecVersion()

	result := target.Bytes()

	assert.Equal(t, []byte(nil), result)
}

func Test_CheckSpecVersion_Validate(t *testing.T) {
	target := setupCheckSpecVersion()

	result, err := target.Validate(constants.OneAddressAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckSpecVersion_ValidateUnsigned(t *testing.T) {
	target := setupCheckSpecVersion()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckSpecVersion_PreDispatch(t *testing.T) {
	target := setupCheckSpecVersion()

	result, err := target.PreDispatch(constants.OneAddressAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)
}

func Test_CheckSpecVersion_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckSpecVersion()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckSpecVersion_PostDispatch(t *testing.T) {
	target := setupCheckSpecVersion()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckSpecVersion_Metadata(t *testing.T) {
	expectType := primitives.NewMetadataTypeWithPath(
		metadata.CheckSpecVersion,
		"CheckSpecVersion",
		sc.Sequence[sc.Str]{"frame_system", "extensions", "check_spec_version", "CheckSpecVersion"},
		primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
	)
	expectSignedExtension := primitives.NewMetadataSignedExtension("CheckSpecVersion", metadata.CheckSpecVersion, metadata.PrimitiveTypesU32)

	resultType, resultSignedExtension := setupCheckSpecVersion().Metadata()

	assert.Equal(t, expectType, resultType)
	assert.Equal(t, expectSignedExtension, resultSignedExtension)
}

func setupCheckSpecVersion() CheckSpecVersion {
	mockModule = new(mocks.SystemModule)

	return NewCheckSpecVersion(mockModule)
}
