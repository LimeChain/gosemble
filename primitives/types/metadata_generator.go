package types

import (
	"reflect"
	"strings"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	indexOptionNone sc.U8 = iota
	indexOptionSome
)

const (
	additionalSignedTypeName = "typesInfoAdditionalSignedData"
	moduleTypeName           = "Module"
	hookOnChargeTypeName     = "OnChargeTransaction"
	varyingDataTypeName      = "VaryingData"
	encodableTypeName        = "Encodable"
	primitivesPackagePath    = "github.com/LimeChain/gosemble/primitives/types."
	goscalePath              = "github.com/LimeChain/goscale."
)

type MetadataTypeGenerator struct {
	MetadataTypes sc.Sequence[MetadataType]
	metadataIds   map[string]int
}

func NewMetadataTypeGenerator() *MetadataTypeGenerator {
	return &MetadataTypeGenerator{
		metadataIds:   BuildMetadataTypesIdsMap(),
		MetadataTypes: sc.Sequence[MetadataType]{},
	}
}

func BuildMetadataTypesIdsMap() map[string]int {
	return map[string]int{
		"Bool":                       metadata.PrimitiveTypesBool,
		"Str":                        metadata.PrimitiveTypesString,
		"U8":                         metadata.PrimitiveTypesU8,
		"U16":                        metadata.PrimitiveTypesU16,
		"U32":                        metadata.PrimitiveTypesU32,
		"U64":                        metadata.PrimitiveTypesU64,
		"U128":                       metadata.PrimitiveTypesU128,
		"I8":                         metadata.PrimitiveTypesI8,
		"I16":                        metadata.PrimitiveTypesI16,
		"I32":                        metadata.PrimitiveTypesI32,
		"I64":                        metadata.PrimitiveTypesI64,
		"I128":                       metadata.PrimitiveTypesI128,
		"H256":                       metadata.TypesH256,
		"SequenceU8":                 metadata.TypesSequenceU8,
		"MultiAddress":               metadata.TypesMultiAddress,
		"Header":                     metadata.Header,
		"SequenceUncheckedExtrinsic": metadata.TypesSequenceUncheckedExtrinsics,
		"SequenceSequence[U8]":       metadata.TypesSequenceSequenceU8,
	}
}

func (g *MetadataTypeGenerator) GetIdsMap() map[string]int {
	return g.metadataIds
}

func (g *MetadataTypeGenerator) GetMetadataTypes() sc.Sequence[MetadataType] {
	return g.MetadataTypes
}

func (g *MetadataTypeGenerator) AppendMetadataTypes(types sc.Sequence[MetadataType]) {
	g.MetadataTypes = append(g.MetadataTypes, types...)
}

// BuildMetadataTypeRecursively Builds the metadata type (recursively) if it does not exist
func (g *MetadataTypeGenerator) BuildMetadataTypeRecursively(v reflect.Value, path *sc.Sequence[sc.Str], def *MetadataTypeDefinition, params *sc.Sequence[MetadataTypeParameter]) int {
	valueType := v.Type()
	typeName := valueType.Name()
	log.NewLogger().Info("Building metadata type: " + typeName)
	var typeId int
	var ok bool
	switch valueType.Kind() {
	case reflect.Struct:
		typeId, ok = g.metadataIds[typeName]
		if !ok {
			typeId, ok = g.isCompactVariation(v)
			if ok {
				log.NewLogger().Info("IsCompactVariation")
				return typeId
			}
			typeId = g.assignNewMetadataId(typeName)
			typeNumFields := valueType.NumField()
			metadataFields := sc.Sequence[MetadataTypeDefinitionField]{}
			for i := 0; i < typeNumFields; i++ {
				fieldName := valueType.Field(i).Name
				fieldTypeName := valueType.Field(i).Type.Name()
				if isIgnoredName(fieldName) || isIgnoredType(fieldTypeName) {
					continue
				}
				log.NewLogger().Info("Current field type name: " + fieldTypeName)
				fieldId, ok := g.metadataIds[fieldTypeName]
				if !ok {
					fieldId = g.BuildMetadataTypeRecursively(v.Field(i), nil, nil, nil)
				}
				if strings.HasPrefix(fieldTypeName, "Sequence") {
					fieldName = "Vec<" + fieldName + ">"
				}
				metadataFields = append(metadataFields, NewMetadataTypeDefinitionFieldWithName(fieldId, sc.Str(fieldName)))
			}
			metadataTypeDef := NewMetadataTypeDefinitionComposite(metadataFields)
			metadataTypePath := sc.Sequence[sc.Str]{}
			metadataTypeParams := sc.Sequence[MetadataTypeParameter]{}

			metadataDocs := typeName
			metadataDocs = strings.Replace(metadataDocs, primitivesPackagePath, "", 1)
			metadataDocs = strings.Replace(metadataDocs, goscalePath, "", 1)
			if def != nil {
				metadataTypeDef = *def
			}
			if path != nil {
				metadataTypePath = *path
			}
			if params != nil {
				metadataTypeParams = *params
			}
			if strings.HasPrefix(typeName, "Option") {
				typeParameterId := g.GetIdsMap()[v.FieldByName("Value").Type().Name()]
				metadataTypeParams = append(metadataTypeParams, NewMetadataTypeParameter(typeParameterId, "T"))
				metadataTypeDef = optionTypeDefinition(v.FieldByName("Value").Type().Name(), typeParameterId)
				metadataDocs = "Option<" + v.FieldByName("Value").Type().Name() + ">"
				metadataTypePath = sc.Sequence[sc.Str]{"Option"}
			}
			newMetadataType := NewMetadataTypeWithParams(typeId, metadataDocs, metadataTypePath, metadataTypeDef, metadataTypeParams)
			g.MetadataTypes = append(g.MetadataTypes, newMetadataType)
			log.NewLogger().Info("Built new type: " + metadataDocs)
		}
	case reflect.Slice:
		sequenceType := valueType.Elem().Name()
		if sequenceType == encodableTypeName { // TransactionSource (alias for sc.VaryingData)
			typeId = g.assignNewMetadataId(typeName)
			newMetadataType := NewMetadataTypeWithParams(typeId, typeName, *path, *def, sc.Sequence[MetadataTypeParameter]{})
			g.MetadataTypes = append(g.MetadataTypes, newMetadataType)
		} else {
			sequenceName := "Sequence"
			log.NewLogger().Info("SequenceType: " + sequenceType)
			//sequenceSeq := valueType.Name()
			//log.NewLogger().Info("SequenceSeq: " + sequenceSeq)
			sequence := sequenceName + sequenceType
			if strings.HasPrefix(sequenceType, "Sequence") { // We are dealing with double sequence (e.g. SequenceSequenceU8)
				log.NewLogger().Info("Double Sequence")
				sequence = strings.Replace(sequence, goscalePath, "", 1)
			}
			log.NewLogger().Info("Sequence Result: " + sequence)
			//sequenceName = sequenceType
			//log.NewLogger().Info("SequenceType: " + sequenceName)
			sequenceTypeId, ok := g.metadataIds[sequence]
			if !ok {
				log.NewLogger().Info(sequence + " does not exist. Generating...")
				if v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
					sequenceTypeId = g.BuildMetadataTypeRecursively(v.Elem(), path, nil, nil)
				}
			}
			typeId = sequenceTypeId
		}

	case reflect.Array: // types U128 and U64
		typeId = g.metadataIds[typeName]
	}
	return typeId
}

// BuildCallsMetadata returns metadata calls type of a module
func (g *MetadataTypeGenerator) BuildCallsMetadata(moduleName string, moduleFunctions map[sc.U8]Call, params *sc.Sequence[MetadataTypeParameter]) int {
	balancesCallsMetadataId := g.assignNewMetadataId(moduleName + "Calls")

	functionVariants := sc.Sequence[MetadataDefinitionVariant]{}

	lenFunctions := len(moduleFunctions)
	for i := 0; i < lenFunctions; i++ {
		f := moduleFunctions[sc.U8(i)]

		functionValue := reflect.ValueOf(f)
		functionType := functionValue.Type()

		functionName := functionType.Name()

		args := functionValue.FieldByName("Arguments")

		fields := sc.Sequence[MetadataTypeDefinitionField]{}

		if args.IsValid() {
			argsLen := args.Len()
			for j := 0; j < argsLen; j++ {
				currentArg := args.Index(j).Elem()
				currentArgId := g.BuildMetadataTypeRecursively(currentArg, nil, nil, nil)
				fields = append(fields, NewMetadataTypeDefinitionField(currentArgId))
			}
		}

		functionVariant := NewMetadataDefinitionVariant(
			g.constructFunctionName(functionName),
			fields,
			sc.U8(i),
			f.Docs())
		functionVariants = append(functionVariants, functionVariant)
	}

	variant := NewMetadataTypeDefinitionVariant(functionVariants)

	g.MetadataTypes = append(g.MetadataTypes, NewMetadataTypeWithParams(balancesCallsMetadataId, moduleName+" calls", sc.Sequence[sc.Str]{sc.Str("pallet_" + strings.ToLower(moduleName)), "pallet", "Call"}, variant, *params))

	return balancesCallsMetadataId
}

// BuildErrorsMetadata returns metadata errors type of a module
func (g *MetadataTypeGenerator) BuildErrorsMetadata(moduleName string, definition *MetadataTypeDefinition) int {
	var errorsTypeId = -1
	var ok bool
	switch moduleName {
	case "System":
		errorsTypeId, ok = g.metadataIds[moduleName+"Errors"]
		if !ok {
			errorsTypeId = g.assignNewMetadataId(moduleName + "Errors")
			g.MetadataTypes = append(g.MetadataTypes, NewMetadataTypeWithPath(errorsTypeId,
				"frame_system pallet Error",
				sc.Sequence[sc.Str]{"frame_system", "pallet", "Error"}, *definition))
		}
	}
	return errorsTypeId
}

// constructFunctionName constructs the formal name of a function call for the module metadata type given its struct name as an input (e.g. callTransferAll -> transfer_all)
func (g *MetadataTypeGenerator) constructFunctionName(input string) string {
	input, _ = strings.CutPrefix(input, "call")
	var result strings.Builder

	for i, char := range input {
		if i > 0 && 'A' <= char && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}

func (g *MetadataTypeGenerator) assignNewMetadataId(name string) int {
	newId := len(g.metadataIds) + 1
	g.metadataIds[name] = newId
	return newId
}

func (g *MetadataTypeGenerator) isCompactVariation(v reflect.Value) (int, bool) {
	log.NewLogger().Info("IsCompactVariation(): " + v.String())
	field := v.FieldByName("Number")
	if field.IsValid() {
		log.NewLogger().Info("Number field present")
		if v.Type() == reflect.TypeOf(sc.Compact{}) {
			switch field.Elem().Type() {
			case reflect.TypeOf(*new(sc.U128)):
				typeId, ok := g.metadataIds["CompactU128"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU128")
					g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU128", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))))
				}
				return typeId, true
			case reflect.TypeOf(*new(sc.U64)):
				typeId, ok := g.metadataIds["CompactU64"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU64")
					g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU64", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))))
				}
				return typeId, true
			case reflect.TypeOf(*new(sc.U32)):
				typeId, ok := g.metadataIds["CompactU32"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU32")
					g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU32", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))))
				}
				return typeId, true
			}
		}
	}
	return -1, false
}

func optionTypeDefinition(typeName string, typeParameterId int) MetadataTypeDefinition {
	return NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"None",
				sc.Sequence[MetadataTypeDefinitionField]{},
				indexOptionNone,
				"Option<"+typeName+">(nil)"),
			NewMetadataDefinitionVariant(
				"Some",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(typeParameterId),
				},
				indexOptionSome,
				"Option<"+typeName+">(value)"),
		})
}

func isIgnoredType(t string) bool {
	return t == moduleTypeName || t == hookOnChargeTypeName || t == varyingDataTypeName
}

func isIgnoredName(name string) bool {
	return name == additionalSignedTypeName
}
