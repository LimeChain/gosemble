package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeExtrinsic[N sc.Numeric] struct {
	modules map[sc.U8]types.Module[N]
	extra   primitives.SignedExtra
}

func New[N sc.Numeric](modules map[sc.U8]types.Module[N], extra primitives.SignedExtra) RuntimeExtrinsic[N] {
	return RuntimeExtrinsic[N]{modules: modules, extra: extra}
}

func (re RuntimeExtrinsic[N]) Module(index sc.U8) (module types.Module[N], isFound bool) {
	m, ok := re.modules[index]
	return m, ok
}

func (re RuntimeExtrinsic[N]) CreateInherents(inherentData primitives.InherentData) []byte {
	i := 0
	var result []byte

	for _, module := range re.modules {
		inherent := module.CreateInherent(inherentData)

		if inherent.HasValue {
			i++
			extrinsic := types.NewUnsignedUncheckedExtrinsic(inherent.Value)
			result = append(result, extrinsic.Bytes()...)
		}
	}

	if i == 0 {
		return []byte{}
	}

	return append(sc.ToCompact(i).Bytes(), result...)
}

func (re RuntimeExtrinsic[N]) CheckInherents(data primitives.InherentData, block types.Block[N]) primitives.CheckInherentsResult {
	result := primitives.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false
		call := extrinsic.Function

		for _, module := range re.modules {
			if module.IsInherent(call) {
				isInherent = true

				err := module.CheckInherent(extrinsic.Function, data)
				if err != nil {
					e := err.(primitives.IsFatalError)
					err := result.PutError(module.InherentIdentifier(), e)
					// TODO: log depending on error type - handle_put_error_result
					if err != nil {
						log.Critical(err.Error())
					}

					if e.IsFatal() {
						return result
					}
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	return result
}

// EnsureInherentsAreFirst checks if the inherents are before non-inherents.
func (re RuntimeExtrinsic[N]) EnsureInherentsAreFirst(block types.Block[N]) int {
	signedExtrinsicFound := false

	for i, extrinsic := range block.Extrinsics {
		isInherent := false

		if extrinsic.IsSigned() {
			signedExtrinsicFound = true
		} else {
			call := extrinsic.Function

			for _, module := range re.modules {
				if module.IsInherent(call) {
					isInherent = true
				}
			}
		}

		if signedExtrinsicFound && isInherent {
			return i
		}
	}

	return -1
}

func (re RuntimeExtrinsic[N]) OnInitialize(n N) primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnInitialize(n))
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OnRuntimeUpgrade() primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnRuntimeUpgrade())
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OnFinalize(n N) {
	for _, m := range re.modules {
		m.OnFinalize(n)
	}
}

func (re RuntimeExtrinsic[N]) OnIdle(n N, remainingWeight primitives.Weight) primitives.Weight {
	weight := primitives.WeightZero()
	for _, m := range re.modules {
		adjustedRemainingWeight := remainingWeight.SaturatingSub(weight)
		weight = weight.SaturatingAdd(m.OnIdle(n, adjustedRemainingWeight))
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OffchainWorker(n N) {
	for _, m := range re.modules {
		m.OffchainWorker(n)
	}
}

func (re RuntimeExtrinsic[N]) Metadata() (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModule], primitives.MetadataExtrinsic) {
	metadataTypes := sc.Sequence[primitives.MetadataType]{}
	modules := sc.Sequence[primitives.MetadataModule]{}

	callVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}
	eventVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}

	// iterate all modules and append their types and modules
	for _, module := range re.modules {
		mTypes, mModule := module.Metadata()

		metadataTypes = append(metadataTypes, mTypes...)
		modules = append(modules, mModule)

		callVariants = append(callVariants, mModule.CallDef)
		eventVariants = append(eventVariants, mModule.EventDef)
	}

	// append runtime event
	metadataTypes = append(metadataTypes, re.runtimeEvent(eventVariants))

	// get the signed extra types and extensions
	signedExtraTypes, signedExtensions := re.extra.Metadata()
	// append to signed extra types to all types
	metadataTypes = append(metadataTypes, signedExtraTypes...)

	// create runtime call type
	runtimeCall := re.runtimeCall(callVariants)
	// append runtime call to all types
	metadataTypes = append(metadataTypes, runtimeCall)

	// create the unchecked extrinsic type using runtime call id
	uncheckedExtrinsicType := primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
		sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
		primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
			}),
		sc.Sequence[primitives.MetadataTypeParameter]{
			primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress, "Address"),
			primitives.NewMetadataTypeParameterCompactId(runtimeCall.Id, "Call"),
			primitives.NewMetadataTypeParameter(metadata.TypesMultiSignature, "Signature"),
			primitives.NewMetadataTypeParameter(metadata.SignedExtra, "Extra"),
		},
	)

	// append it to all types
	metadataTypes = append(metadataTypes, uncheckedExtrinsicType)

	// create the metadata extrinsic, which uses the id of the unchecked extrinsic and signed extra extensions
	extrinsic := primitives.MetadataExtrinsic{
		Type:             uncheckedExtrinsicType.Id,
		Version:          types.ExtrinsicFormatVersion,
		SignedExtensions: signedExtensions,
	}

	return metadataTypes, modules, extrinsic
}

func (re RuntimeExtrinsic[N]) runtimeCall(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]) primitives.MetadataType {
	runtimeCallSubTypes := sc.Sequence[primitives.MetadataDefinitionVariant]{}

	for _, v := range variants {
		if v.HasValue {
			runtimeCallSubTypes = append(runtimeCallSubTypes, v.Value)
		}
	}

	return primitives.NewMetadataTypeWithPath(
		metadata.RuntimeCall,
		"RuntimeCall",
		sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
		primitives.NewMetadataTypeDefinitionVariant(runtimeCallSubTypes),
	)
}

func (re RuntimeExtrinsic[N]) runtimeEvent(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]) primitives.MetadataType {
	runtimeEventSubTypes := sc.Sequence[primitives.MetadataDefinitionVariant]{}

	for _, v := range variants {
		if v.HasValue {
			runtimeEventSubTypes = append(runtimeEventSubTypes, v.Value)
		}
	}

	return primitives.NewMetadataTypeWithPath(
		metadata.TypesRuntimeEvent,
		"node_template_runtime RuntimeEvent",
		sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"},
		primitives.NewMetadataTypeDefinitionVariant(runtimeEventSubTypes),
	)
}
