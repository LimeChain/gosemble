package extensions

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	invalidTransactionBadSigner, _ = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadSigner())
)

func Test_CheckNonZeroAddress_AdditionalSigned(t *testing.T) {
	target := setupCheckNonZeroSender()

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, primitives.AdditionalSigned{}, result)
}

func Test_CheckNonZeroAddress_Encode(t *testing.T) {
	target := setupCheckNonZeroSender()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, 0, buffer.Len())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_CheckNonZeroAddress_Decode(t *testing.T) {
	target := setupCheckNonZeroSender()
	value := []byte{1, 2, 3}
	buffer := bytes.NewBuffer(value)

	target.Decode(buffer)

	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, bytes.NewBuffer(value), buffer)
}

func Test_CheckNonZeroAddress_Bytes(t *testing.T) {
	target := setupCheckNonZeroSender()

	result := target.Bytes()

	assert.Equal(t, []byte(nil), result)
}

func Test_CheckNonZeroAddress_Validate_Success(t *testing.T) {
	target := setupCheckNonZeroSender()

	result, err := target.Validate(constants.OneAddressAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckNonZeroAddress_Validate_Fails(t *testing.T) {
	target := setupCheckNonZeroSender()

	result, err := target.Validate(constants.ZeroAddressAccountId, nil, nil, sc.Compact{})

	assert.Equal(t, invalidTransactionBadSigner, err)
	assert.Equal(t, primitives.ValidTransaction{}, result)
}

func Test_CheckNonZeroAddress_ValidateUnsigned(t *testing.T) {
	target := setupCheckNonZeroSender()

	result, err := target.ValidateUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)
}

func Test_CheckNonZeroAddress_PreDispatch(t *testing.T) {
	target := setupCheckNonZeroSender()

	result, err := target.PreDispatch(constants.OneAddressAccountId, nil, nil, sc.Compact{})

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)
}

func Test_CheckNonZeroAddress_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckNonZeroSender()

	err := target.PreDispatchUnsigned(nil, nil, sc.Compact{})

	assert.Nil(t, err)
}

func Test_CheckNonZeroAddress_PostDispatch(t *testing.T) {
	target := setupCheckNonZeroSender()

	err := target.PostDispatch(sc.NewOption[primitives.Pre](nil), nil, nil, sc.Compact{}, nil)

	assert.Nil(t, err)
}

func Test_CheckNonZeroAddress_Metadata(t *testing.T) {
	expectType := primitives.NewMetadataTypeWithPath(
		metadata.CheckNonZeroSender,
		"CheckNonZeroSender",
		sc.Sequence[sc.Str]{"frame_system", "extensions", "check_non_zero_sender", "CheckNonZeroSender"},
		primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
	)
	expectSignedExtension := primitives.NewMetadataSignedExtension("CheckNonZeroSender", metadata.CheckNonZeroSender, metadata.TypesEmptyTuple)

	resultType, resultSignedExtension := setupCheckNonZeroSender().Metadata()

	assert.Equal(t, expectType, resultType)
	assert.Equal(t, expectSignedExtension, resultSignedExtension)
}

func setupCheckNonZeroSender() CheckNonZeroAddress {
	return NewCheckNonZeroAddress()
}
