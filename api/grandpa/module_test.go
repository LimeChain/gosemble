package grandpa

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	target          Module
	mockGrandpa     *mocks.GrandpaModule
	mockMemoryUtils *mocks.MemoryTranslator
)

func setup() {
	mockGrandpa = new(mocks.GrandpaModule)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target = New(mockGrandpa)
	target.memUtils = mockMemoryUtils
}

func Test_Name(t *testing.T) {
	setup()

	assert.Equal(t, "GrandpaApi", target.Name())
}

func Test_Item(t *testing.T) {
	setup()

	hash := common.MustBlake2b8([]byte("GrandpaApi"))

	expected := types.ApiItem{
		Name:    sc.BytesToFixedSequenceU8(hash[:]),
		Version: 3,
	}

	assert.Equal(t, expected, target.Item())
}

func Test_Authorities_None(t *testing.T) {
	setup()

	authorities := sc.Sequence[types.Authority]{
		{
			Id:     constants.ZeroAddress.FixedSequence,
			Weight: sc.U64(64),
		},
	}

	mockGrandpa.On("Authorities").Return(authorities)
	mockMemoryUtils.On("BytesToOffsetAndSize", authorities.Bytes()).Return(int64(13))

	target.Authorities()

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", authorities.Bytes())
	mockMemoryUtils.AssertNumberOfCalls(t, "BytesToOffsetAndSize", 1)
}
