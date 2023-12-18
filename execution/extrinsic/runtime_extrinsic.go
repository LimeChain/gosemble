package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeExtrinsic interface {
	Module(index sc.U8) (module primitives.Module, isFound bool)
	CreateInherents(inherentData primitives.InherentData) ([]byte, error)
	CheckInherents(data primitives.InherentData, block primitives.Block) (primitives.CheckInherentsResult, error)
	EnsureInherentsAreFirst(block primitives.Block) int
	OnInitialize(n sc.U64) (primitives.Weight, error)
	OnRuntimeUpgrade() primitives.Weight
	OnFinalize(n sc.U64) error
	OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight
	OffchainWorker(n sc.U64)
	Metadata(metadataTypesIds map[string]int) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV14], primitives.MetadataExtrinsicV14)
	MetadataLatest(metadataTypesIds map[string]int) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV15], primitives.MetadataExtrinsicV15, primitives.OuterEnums, primitives.CustomMetadata)
}

type runtimeExtrinsic struct {
	modules []primitives.Module
	extra   primitives.SignedExtra
	logger  log.DebugLogger
}

func New(modules []primitives.Module, extra primitives.SignedExtra, logger log.DebugLogger) RuntimeExtrinsic {
	return runtimeExtrinsic{
		modules: modules,
		extra:   extra,
		logger:  logger,
	}
}

func (re runtimeExtrinsic) Module(index sc.U8) (primitives.Module, bool) {
	mod, err := primitives.GetModule(index, re.modules)
	if err != nil {
		return nil, false
	}
	return mod, true
}

func (re runtimeExtrinsic) CreateInherents(inherentData primitives.InherentData) ([]byte, error) {
	i := 0
	var result []byte

	for _, module := range re.modules {
		inherent, err := module.CreateInherent(inherentData)
		if err != nil {
			return []byte{}, err
		}

		if inherent.HasValue {
			i++
			extrinsic := types.NewUnsignedUncheckedExtrinsic(inherent.Value)
			result = append(result, extrinsic.Bytes()...)
		}
	}

	if i == 0 {
		return []byte{}, nil
	}

	return append(sc.ToCompact(i).Bytes(), result...), nil
}

func (re runtimeExtrinsic) CheckInherents(data primitives.InherentData, block primitives.Block) (primitives.CheckInherentsResult, error) {
	result := primitives.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics() {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false
		call := extrinsic.Function()

		for _, module := range re.modules {
			if !module.IsInherent(call) {
				continue
			}

			isInherent = true

			if err := module.CheckInherent(call, data); err != nil {
				fatalErr, ok := err.(primitives.FatalError)
				if !ok {
					return result, err
				}

				// todo add new tests for assertPanic in the api modules tests
				if err := result.PutError(module.InherentIdentifier(), fatalErr); err != nil {
					if inherentErr, ok := err.(primitives.InherentError); ok && inherentErr.VaryingData[0] == primitives.InherentErrorInherentDataExists {
						re.logger.Debug(inherentErr.Error())
					} else {
						return result, err
					}
				}

				if fatalErr.IsFatal() {
					return result, nil
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	return result, nil
}

// func handlePutErrorResult(err error) error {
// 	if err == nil {
// 		return err
// 	}

// 	inherentErr, ok := err.(primitives.InherentError)
// 	if !ok {
// 		return fmt.Errorf("Unexpected error from `put_error` operation: %v", err)
// 	}

// 	switch inherentErr.VaryingData[0] {
// 	case primitives.InherentErrorInherentDataExists:

// 	case primitives.InherentErrorFatalErrorReported:
// 		return inherentErr
// 	}
// }

// EnsureInherentsAreFirst checks if the inherents are before non-inherents.
func (re runtimeExtrinsic) EnsureInherentsAreFirst(block primitives.Block) int {
	signedExtrinsicFound := false

	for i, extrinsic := range block.Extrinsics() {
		isInherent := false

		if extrinsic.IsSigned() {
			signedExtrinsicFound = true
		} else {
			call := extrinsic.Function()

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

func (re runtimeExtrinsic) OnInitialize(n sc.U64) (primitives.Weight, error) {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		init, err := m.OnInitialize(n)
		if err != nil {
			return primitives.Weight{}, err
		}
		weight = weight.Add(init)
	}

	return weight, nil
}

func (re runtimeExtrinsic) OnRuntimeUpgrade() primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnRuntimeUpgrade())
	}

	return weight
}

func (re runtimeExtrinsic) OnFinalize(n sc.U64) error {
	for _, m := range re.modules {
		err := m.OnFinalize(n)
		if err != nil {
			return err
		}
	}
	return nil
}

func (re runtimeExtrinsic) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	weight := primitives.WeightZero()
	for _, m := range re.modules {
		adjustedRemainingWeight := remainingWeight.SaturatingSub(weight)
		weight = weight.SaturatingAdd(m.OnIdle(n, adjustedRemainingWeight))
	}

	return weight
}

func (re runtimeExtrinsic) OffchainWorker(n sc.U64) {
	for _, m := range re.modules {
		m.OffchainWorker(n)
	}
}

func (re runtimeExtrinsic) Metadata(metadataTypesIds map[string]int) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV14], primitives.MetadataExtrinsicV14) {
	metadataTypes := sc.Sequence[primitives.MetadataType]{}
	modules := sc.Sequence[primitives.MetadataModuleV14]{}

	callVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}
	eventVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}
	errorVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}

	// iterate all modules and append their types and modules
	for _, module := range re.modules {
		mTypes, mModule := module.Metadata()

		metadataTypes = append(metadataTypes, mTypes...)
		mModuleV14 := mModule.ModuleV14
		modules = append(modules, mModuleV14)

		callVariants = append(callVariants, mModuleV14.CallDef)
		eventVariants = append(eventVariants, mModuleV14.EventDef)
		errorVariants = append(errorVariants, mModuleV14.ErrorDef)
	}

	// append runtime event
	metadataTypes = append(metadataTypes, re.runtimeEvent(eventVariants))

	// get the signed extra types and extensions
	signedExtraTypes, signedExtensions := re.extra.Metadata(metadataTypesIds)
	// append to signed extra types to all types
	metadataTypes = append(metadataTypes, signedExtraTypes...)

	// create runtime call type
	runtimeCall := re.runtimeCall(callVariants)
	// append runtime call to all types
	metadataTypes = append(metadataTypes, runtimeCall)

	runtimeError := re.runtimeError(errorVariants)

	metadataTypes = append(metadataTypes, runtimeError)

	// create the unchecked extrinsic type using runtime call id
	uncheckedExtrinsicType := createUncheckedExtrinsicType(runtimeCall)

	// append it to all types
	metadataTypes = append(metadataTypes, uncheckedExtrinsicType)

	// create the metadata extrinsic, which uses the id of the unchecked extrinsic and signed extra extensions
	extrinsic := primitives.MetadataExtrinsicV14{
		Type:             uncheckedExtrinsicType.Id,
		Version:          types.ExtrinsicFormatVersion,
		SignedExtensions: signedExtensions,
	}

	return metadataTypes, modules, extrinsic
}

func (re runtimeExtrinsic) MetadataLatest(metadataTypesIds map[string]int) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV15], primitives.MetadataExtrinsicV15, primitives.OuterEnums, primitives.CustomMetadata) {
	metadataTypes := sc.Sequence[primitives.MetadataType]{}
	modules := sc.Sequence[primitives.MetadataModuleV15]{}

	callVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}
	eventVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}
	errorVariants := sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]{}

	outerEnums := primitives.OuterEnums{
		CallEnumType:  sc.ToCompact(metadata.RuntimeCall),
		EventEnumType: sc.ToCompact(metadata.TypesRuntimeEvent),
		ErrorEnumType: sc.ToCompact(metadata.TypesRuntimeError),
	}

	custom := primitives.CustomMetadata{
		Map: sc.Dictionary[sc.Str, primitives.CustomValueMetadata]{},
	}

	// iterate all modules and append their types and modules
	for _, module := range re.modules {
		mTypes, mModule := module.Metadata()

		moduleV15 := mModule.ModuleV15

		metadataTypes = append(metadataTypes, mTypes...)
		modules = append(modules, moduleV15)

		callVariants = append(callVariants, moduleV15.CallDef)
		eventVariants = append(eventVariants, moduleV15.EventDef)
		errorVariants = append(errorVariants, moduleV15.ErrorDef)
	}

	// append runtime event
	metadataTypes = append(metadataTypes, re.runtimeEvent(eventVariants))

	// get the signed extra types and extensions
	signedExtraTypes, signedExtensions := re.extra.Metadata(metadataTypesIds)
	// append to signed extra types to all types
	metadataTypes = append(metadataTypes, signedExtraTypes...)

	// create runtime call type
	runtimeCall := re.runtimeCall(callVariants)
	// append runtime call to all types
	metadataTypes = append(metadataTypes, runtimeCall)

	runtimeError := re.runtimeError(errorVariants)

	metadataTypes = append(metadataTypes, runtimeError)

	// create the unchecked extrinsic type using runtime call id
	uncheckedExtrinsicType := createUncheckedExtrinsicType(runtimeCall)

	// append it to all types
	metadataTypes = append(metadataTypes, uncheckedExtrinsicType)

	extrinsicV15 := primitives.MetadataExtrinsicV15{
		Version:          types.ExtrinsicFormatVersion,
		Address:          sc.ToCompact(metadata.TypesMultiAddress),
		Call:             runtimeCall.Id,
		Signature:        sc.ToCompact(metadata.TypesMultiSignature),
		Extra:            sc.ToCompact(metadata.SignedExtra),
		SignedExtensions: signedExtensions,
	}

	return metadataTypes, modules, extrinsicV15, outerEnums, custom
}

func createUncheckedExtrinsicType(runtimeCall primitives.MetadataType) primitives.MetadataType {
	return primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
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
}

func (re runtimeExtrinsic) runtimeCall(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]) primitives.MetadataType {
	return re.runtimeType(
		variants,
		metadata.RuntimeCall,
		"RuntimeCall",
		sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
	)
}

func (re runtimeExtrinsic) runtimeEvent(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]) primitives.MetadataType {
	return re.runtimeType(
		variants,
		metadata.TypesRuntimeEvent,
		"node_template_runtime RuntimeEvent",
		sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"},
	)
}

func (re runtimeExtrinsic) runtimeError(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]]) primitives.MetadataType {
	return re.runtimeType(
		variants,
		metadata.TypesRuntimeError,
		"node_template_runtime RuntimeError",
		sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeError"},
	)
}

func (re runtimeExtrinsic) runtimeType(variants sc.Sequence[sc.Option[primitives.MetadataDefinitionVariant]], id int, docs string, path sc.Sequence[sc.Str]) primitives.MetadataType {
	subTypes := sc.Sequence[primitives.MetadataDefinitionVariant]{}

	for _, v := range variants {
		if v.HasValue {
			subTypes = append(subTypes, v.Value)
		}
	}

	return primitives.NewMetadataTypeWithPath(
		id,
		docs,
		path,
		primitives.NewMetadataTypeDefinitionVariant(subTypes),
	)
}
