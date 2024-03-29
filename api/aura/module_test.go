package aura

import (
	"errors"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	target          Module
	mockAura        *mocks.AuraModule
	mockMemoryUtils *mocks.MemoryTranslator
)

func setup() {
	mockAura = new(mocks.AuraModule)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target = New(mockAura, log.NewLogger())
	target.memUtils = mockMemoryUtils
}

func Test_Name(t *testing.T) {
	setup()

	assert.Equal(t, "AuraApi", target.Name())
}

func Test_Item(t *testing.T) {
	setup()

	hash := common.MustBlake2b8([]byte("AuraApi"))

	expected := types.ApiItem{
		Name:    sc.BytesToFixedSequenceU8(hash[:]),
		Version: 1,
	}

	assert.Equal(t, expected, target.Item())
}

func Test_Authorities_None(t *testing.T) {
	setup()

	mockAura.On("StorageAuthoritiesBytes").Return(sc.NewOption[sc.Sequence[sc.U8]](nil), nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", []byte{0}).Return(int64(0))

	target.Authorities()

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", []byte{0})
	mockMemoryUtils.AssertNumberOfCalls(t, "BytesToOffsetAndSize", 1)
}

func Test_Authorities_Some(t *testing.T) {
	setup()

	mockAura.On("StorageAuthoritiesBytes").Return(sc.NewOption[sc.Sequence[sc.U8]](
		sc.Sequence[sc.U8]{sc.U8(1), sc.U8(2), sc.U8(3)},
	), nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", []byte{1, 2, 3}).Return(int64(13))

	target.Authorities()

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", []byte{1, 2, 3})
	mockMemoryUtils.AssertNumberOfCalls(t, "BytesToOffsetAndSize", 1)
}

func Test_Authorities_Panics(t *testing.T) {
	setup()

	expectedErr := errors.New("panic")
	mockAura.On("StorageAuthoritiesBytes").Return(sc.NewOption[sc.Sequence[sc.U8]](
		sc.Sequence[sc.U8]{},
	), expectedErr)
	mockMemoryUtils.On("BytesToOffsetAndSize", []byte{1, 2, 3}).Return(int64(13))

	assert.PanicsWithValue(t,
		expectedErr.Error(),
		func() { target.Authorities() },
	)

	mockAura.AssertCalled(t, "StorageAuthoritiesBytes")
	mockMemoryUtils.AssertNotCalled(t, "BytesToOffsetAndSize", []byte{1, 2, 3})
}

func Test_SlotDuration(t *testing.T) {
	setup()

	duration := sc.U64(123456)
	durationBytes := duration.Bytes()

	mockAura.On("SlotDuration").Return(duration)
	mockMemoryUtils.On("BytesToOffsetAndSize", durationBytes).Return(int64(13))

	target.SlotDuration()

	mockAura.AssertCalled(t, "SlotDuration")
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", durationBytes)
}

func Test_Module_Metadata(t *testing.T) {
	setup()

	expect := types.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[types.RuntimeApiMethodMetadata]{
			types.RuntimeApiMethodMetadata{
				Name:   "authorities",
				Inputs: sc.Sequence[types.RuntimeApiMethodParamMetadata]{},
				Output: sc.ToCompact(metadata.TypesSequenceU8),
				Docs:   sc.Sequence[sc.Str]{""},
			},
			types.RuntimeApiMethodMetadata{
				Name:   "slot_duration",
				Inputs: sc.Sequence[types.RuntimeApiMethodParamMetadata]{},
				Output: sc.ToCompact(metadata.PrimitiveTypesU64),
				Docs:   sc.Sequence[sc.Str]{},
			},
		},
		Docs: sc.Sequence[sc.Str]{" The API to query account nonce."},
	}

	assert.Equal(t, expect, target.Metadata())
}
