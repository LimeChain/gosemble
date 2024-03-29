package extensions

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_CheckTxVersion_AdditionalSigned(t *testing.T) {
	version := primitives.RuntimeVersion{
		TransactionVersion: 3,
	}
	target := setupCheckTxVersion()

	mockModule.On("Version").Return(version)

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, sc.NewVaryingData(version.TransactionVersion), result)
	mockModule.AssertCalled(t, "Version")
}

func Test_CheckTxVersion_Encode(t *testing.T) {
	target := setupCheckTxVersion()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, 0, buffer.Len())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_CheckTxVersion_Decode(t *testing.T) {
	target := setupCheckTxVersion()
	value := []byte{1, 2, 3}
	buffer := bytes.NewBuffer(value)

	target.Decode(buffer)

	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, bytes.NewBuffer(value), buffer)
}

func Test_CheckTxVersion_Bytes(t *testing.T) {
	target := setupCheckTxVersion()

	result := target.Bytes()

	assert.Equal(t, []byte(nil), result)
}

func Test_CheckTxVersion_DeepCopy(t *testing.T) {
	target := setupCheckTxVersion()

	result := target.DeepCopy()

	assert.Equal(t, &target, result)

	target.typesInfoAdditionalSignedData = nil
	assert.NotEqual(t, &target, result)
}

func Test_CheckTxVersion_Validate(t *testing.T) {
	target := setupCheckTxVersion()

	result, err := target.Validate(constants.OneAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckTxVersion_ValidateUnsigned(t *testing.T) {
	target := setupCheckTxVersion()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckTxVersion_PreDispatch(t *testing.T) {
	target := setupCheckTxVersion()

	result, err := target.PreDispatch(constants.OneAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)
}

func Test_CheckTxVersion_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckTxVersion()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckTxVersion_PostDispatch(t *testing.T) {
	target := setupCheckTxVersion()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckTxVersion_ModulePath(t *testing.T) {
	target := setupCheckTxVersion()

	expectedModulePath := "frame_system"
	actualModulePath := target.ModulePath()

	assert.Equal(t, expectedModulePath, actualModulePath)
}

func setupCheckTxVersion() CheckTxVersion {
	mockModule = new(mocks.SystemModule)
	extension, ok := NewCheckTxVersion(mockModule).(*CheckTxVersion)
	if !ok {
		panic("invalid type assert for *CheckTxVersion")
	}
	return *extension
}
