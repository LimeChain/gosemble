package grandpa

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GrandpaApi"
	apiVersion    = 3
)

type Module struct {
	grandpa  grandpa.GrandpaModule
	memUtils utils.WasmMemoryTranslator
}

func New(grandpa grandpa.GrandpaModule) Module {
	return Module{
		grandpa:  grandpa,
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (m Module) Authorities() int64 {
	authorities, err := m.grandpa.Authorities()
	if err != nil {
		log.Critical(err.Error())
	}
	return m.memUtils.BytesToOffsetAndSize(authorities.Bytes())
}

func (m Module) Metadata() primitives.RuntimeApiMetadata {
	methods := sc.Sequence[primitives.RuntimeApiMethodMetadata]{
		primitives.RuntimeApiMethodMetadata{
			Name:   "grandpa_authorities",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
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
	}

	return primitives.RuntimeApiMetadata{
		Name:    ApiModuleName,
		Methods: methods,
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
}
