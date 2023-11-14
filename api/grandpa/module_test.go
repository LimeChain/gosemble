package grandpa

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
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
			Id:     constants.ZeroAddressAccountId,
			Weight: sc.U64(64),
		},
	}

	mockGrandpa.On("Authorities").Return(authorities, nil)
	mockMemoryUtils.On("BytesToOffsetAndSize", authorities.Bytes()).Return(int64(13))

	target.Authorities()

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", authorities.Bytes())
	mockMemoryUtils.AssertNumberOfCalls(t, "BytesToOffsetAndSize", 1)
}

func Test_Module_Metadata(t *testing.T) {
	setup()

	expect := types.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[types.RuntimeApiMethodMetadata]{
			types.RuntimeApiMethodMetadata{
				Name:   "grandpa_authorities",
				Inputs: sc.Sequence[types.RuntimeApiMethodParamMetadata]{},
				Output: sc.ToCompact(metadata.TypesSequenceTupleGrandpaAppPublic),
				Docs: sc.Sequence[sc.Str]{
					" Get the current GRANDPA authorities and weights. This should not change except",
					" for when changes are scheduled and the corresponding delay has passed.",
					"",
					" When called at block B, it will return the set of authorities that should be",
					" used to finalize descendants of this block (B+1, B+2, ...). The block B itself",
					" is finalized by the authorities from block B-1.",
				},
			},
		},
		Docs: sc.Sequence[sc.Str]{
			" APIs for integrating the GRANDPA finality gadget into runtimes.",
			" This should be implemented on the runtime side.",
			"",
			" This is primarily used for negotiating authority-set changes for the",
			" gadget. GRANDPA uses a signaling model of changing authority sets:",
			" changes should be signaled with a delay of N blocks, and then automatically",
			" applied in the runtime after those N blocks have passed.",
			"",
			" The consensus protocol will coordinate the handoff externally.",
		},
	}

	assert.Equal(t, expect, target.Metadata())
}
