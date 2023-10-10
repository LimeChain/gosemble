package grandpa

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

const moduleId = sc.U8(3)

var (
	mockStorageAuthorities *mocks.StorageValue[primitives.VersionedAuthorityList]
	target                 Module
)

func Test_Module_New(t *testing.T) {
	setup()

	assert.Equal(t, Module{
		DefaultInherentProvider: primitives.DefaultInherentProvider{},
		DefaultDispatchModule:   hooks.DefaultDispatchModule{},
		Index:                   moduleId,
		storage: &storage{
			mockStorageAuthorities,
		},
	}, target)
}

func Test_Module_KeyType(t *testing.T) {
	setup()
	assert.Equal(t, primitives.PublicKeyEd25519, target.KeyType())
}

func Test_Module_KeyTypeId(t *testing.T) {
	setup()
	assert.Equal(t, KeyTypeId, target.KeyTypeId())
}

func Test_Module_GetIndex(t *testing.T) {
	setup()
	assert.Equal(t, moduleId, target.GetIndex())
}

func Test_Module_Functions(t *testing.T) {
	setup()

	assert.Equal(t, map[sc.U8]primitives.Call{}, target.Functions())
}

func Test_Module_PreDispatch(t *testing.T) {
	setup()

	result, err := target.PreDispatch(new(mocks.Call))

	assert.Nil(t, err)
	assert.Equal(t, sc.Empty{}, result)
}

func Test_Module_ValidateUnsigned(t *testing.T) {
	setup()

	result, err := target.ValidateUnsigned(primitives.NewTransactionSourceLocal(), new(mocks.Call))

	assert.Equal(t, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator()), err)
	assert.Equal(t, primitives.ValidTransaction{}, result)
}

func Test_Module_Authorities_Success(t *testing.T) {
	setup()
	expectAuthorites := sc.Sequence[primitives.Authority]{
		{
			Id:     constants.ZeroAddress.FixedSequence,
			Weight: 5,
		},
	}
	storageAuthorites := primitives.VersionedAuthorityList{
		Version:       AuthorityVersion,
		AuthorityList: expectAuthorites,
	}

	mockStorageAuthorities.On("Get").Return(storageAuthorites)

	result := target.Authorities()

	assert.Equal(t, expectAuthorites, result)
	mockStorageAuthorities.AssertCalled(t, "Get")
}

func Test_Module_Authorities_DifferentVersion(t *testing.T) {
	setup()
	storageAuthorites := primitives.VersionedAuthorityList{
		Version: sc.U8(255),
		AuthorityList: sc.Sequence[primitives.Authority]{
			{
				Id:     constants.ZeroAddress.FixedSequence,
				Weight: sc.U64(64),
			},
		},
	}

	mockStorageAuthorities.On("Get").Return(storageAuthorites)

	result := target.Authorities()

	assert.Equal(t, sc.Sequence[primitives.Authority]{}, result)
	mockStorageAuthorities.AssertCalled(t, "Get")
}

func Test_Module_Metadata(t *testing.T) {
	setup()

	expectMetadataTypes := sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParams(metadata.GrandpaCalls, "Grandpa calls", sc.Sequence[sc.Str]{"pallet_grandpa", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			})}
	expectMetadataModule := primitives.MetadataModule{
		Name:      name,
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](nil),
		CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     moduleId,
	}

	metadataTypes, metadataModule := target.Metadata()

	assert.Equal(t, expectMetadataTypes, metadataTypes)
	assert.Equal(t, expectMetadataModule, metadataModule)
}

func setup() {
	mockStorageAuthorities = new(mocks.StorageValue[primitives.VersionedAuthorityList])
	target = New(moduleId)

	target.storage.Authorities = mockStorageAuthorities
}
