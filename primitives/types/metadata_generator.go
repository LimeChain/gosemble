package types

import (
	"reflect"
	"strings"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
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
	goscalePathTrim          = "github.com/LimeChain/goscale."
	goscalePath              = "github.com/LimeChain/goscale"
)

type MetadataTypeGenerator struct {
	metadataTypes sc.Sequence[MetadataType]
	metadataIds   map[string]int
}

func NewMetadataTypeGenerator() *MetadataTypeGenerator {
	return &MetadataTypeGenerator{
		metadataIds:   BuildMetadataTypesIdsMap(),
		metadataTypes: sc.Sequence[MetadataType]{},
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

func (g *MetadataTypeGenerator) ClearMetadata() {
	g.metadataTypes = sc.Sequence[MetadataType]{}
	g.metadataIds = BuildMetadataTypesIdsMap()
}

func (g *MetadataTypeGenerator) GetIdsMap() map[string]int {
	return g.metadataIds
}

func (g *MetadataTypeGenerator) GetId(typeName string) (int, bool) {
	id, ok := g.metadataIds[typeName]
	return id, ok
}

func (g *MetadataTypeGenerator) GetMapLength() int {
	return len(g.metadataIds)
}

func (g *MetadataTypeGenerator) GetMetadataTypes() sc.Sequence[MetadataType] {
	return g.metadataTypes
}

func (g *MetadataTypeGenerator) AppendMetadataTypes(types sc.Sequence[MetadataType]) {
	g.metadataTypes = append(g.metadataTypes, types...)
}

// BuildMetadataTypeRecursively Builds the metadata type (recursively) if it does not exist
func (g *MetadataTypeGenerator) BuildMetadataTypeRecursively(v reflect.Value, path *sc.Sequence[sc.Str], def *MetadataTypeDefinition, params *sc.Sequence[MetadataTypeParameter]) int {
	valueType := v.Type()
	typeName := valueType.Name()
	typeId, ok := g.GetId(typeName)
	if ok {
		return typeId
	}
	switch valueType.Kind() {
	case reflect.Struct:
		typeId, ok = g.isCompactVariation(v)
		if ok {
			return typeId
		}
		typeId = g.assignNewMetadataId(typeName)
		if strings.HasPrefix(typeName, "Option") {
			return g.constructOptionType(v, typeId)
		}
		metadataTypeFields := g.constructTypeFields(v)
		metadataTypeDef := NewMetadataTypeDefinitionComposite(metadataTypeFields)
		metadataTypePath := sc.Sequence[sc.Str]{}
		metadataTypeParams := sc.Sequence[MetadataTypeParameter]{}
		metadataDocs := constructTypeDocs(typeName)
		if def != nil {
			metadataTypeDef = *def
		}
		if path != nil {
			metadataTypePath = *path
		}
		if params != nil {
			metadataTypeParams = *params
		}
		newMetadataType := NewMetadataTypeWithParams(typeId, metadataDocs, metadataTypePath, metadataTypeDef, metadataTypeParams)
		g.metadataTypes = append(g.metadataTypes, newMetadataType)
		return typeId
	case reflect.Slice:
		return g.buildSequenceType(v, path, def)
	case reflect.Array:
		return g.metadataIds[typeName]
	default:
		return typeId
	}
}

// BuildCallsMetadata returns metadata calls type of a module
func (g *MetadataTypeGenerator) BuildCallsMetadata(moduleName string, moduleFunctions map[sc.U8]Call, params *sc.Sequence[MetadataTypeParameter]) int {
	callsMetadataId := g.assignNewMetadataId(moduleName + "Calls")

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
			constructFunctionName(functionName),
			fields,
			sc.U8(i),
			f.Docs())
		functionVariants = append(functionVariants, functionVariant)
	}

	variant := NewMetadataTypeDefinitionVariant(functionVariants)

	g.metadataTypes = append(g.metadataTypes, NewMetadataTypeWithParams(callsMetadataId, moduleName+" calls", sc.Sequence[sc.Str]{sc.Str("pallet_" + strings.ToLower(moduleName)), "pallet", "Call"}, variant, *params))

	return callsMetadataId
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
			g.metadataTypes = append(g.metadataTypes, NewMetadataTypeWithPath(errorsTypeId,
				"frame_system pallet Error",
				sc.Sequence[sc.Str]{"frame_system", "pallet", "Error"}, *definition))
		}
	}
	return errorsTypeId
}

func (g *MetadataTypeGenerator) BuildModuleConstants(config reflect.Value) sc.Sequence[MetadataModuleConstant] {
	var constants sc.Sequence[MetadataModuleConstant]
	configType := config.Type()

	typeNumFields := configType.NumField()
	for i := 0; i < typeNumFields; i++ {
		fieldValue := config.Field(i)
		fieldName := configType.Field(i).Name
		fieldTypeName := configType.Field(i).Type.Name()

		var fieldId int
		fieldValueNumFields := fieldValue.NumField()
		valueEncodable, ok := fieldValue.Interface().(sc.Encodable)
		if ok && fieldValueNumFields == 1 {
			encodableField := fieldValue.Field(0)
			valueEncodable, ok = encodableField.Interface().(sc.Encodable)
			if ok {
				fieldPkgPath := encodableField.Type().PkgPath()
				fieldId, ok = g.metadataIds[encodableField.Type().Name()]
				if ok && fieldPkgPath != goscalePath {
					fieldId, ok = g.metadataIds[fieldTypeName]
					if !ok {
						fieldId = g.BuildMetadataTypeRecursively(fieldValue, nil, nil, nil)
					}
				}
			}
		} else {
			fieldId, ok = g.metadataIds[fieldTypeName]
			if !ok {
				fieldId = g.BuildMetadataTypeRecursively(config.Field(i), nil, nil, nil)
			}
		}

		var docs string
		describerValue, ok := fieldValue.Interface().(Describer)
		if ok {
			docs = describerValue.Docs()
		}

		constant := NewMetadataModuleConstant(
			fieldName,
			sc.ToCompact(fieldId),
			sc.BytesToSequenceU8(valueEncodable.Bytes()),
			docs,
		)
		constants = append(constants, constant)
	}

	return constants
}

func (g *MetadataTypeGenerator) assignNewMetadataId(name string) int {
	newId := len(g.metadataIds) + 1
	g.metadataIds[name] = newId
	return newId
}

func (g *MetadataTypeGenerator) buildSequenceType(v reflect.Value, path *sc.Sequence[sc.Str], def *MetadataTypeDefinition) int {
	valueType := v.Type()
	typeName := valueType.Name()
	sequenceType := valueType.Elem().Name()
	var typeId int
	if sequenceType == encodableTypeName { // TransactionSource (alias for sc.VaryingData)
		typeId = g.assignNewMetadataId(typeName)
		newMetadataType := NewMetadataTypeWithParams(typeId, typeName, *path, *def, sc.Sequence[MetadataTypeParameter]{})
		g.metadataTypes = append(g.metadataTypes, newMetadataType)
	} else {
		sequenceName := "Sequence"
		sequence := sequenceName + sequenceType
		if strings.HasPrefix(sequenceType, "Sequence") { // We are dealing with double sequence (e.g. SequenceSequenceU8)
			sequence = strings.Replace(sequence, goscalePathTrim, "", 1)
		}

		sequenceTypeId, ok := g.metadataIds[sequenceType] // ApiItem
		if !ok {
			n := reflect.Zero(valueType.Elem())
			sequenceTypeId = g.BuildMetadataTypeRecursively(n, path, nil, nil)
		}

		sequenceId, ok := g.metadataIds[sequence]
		if !ok {
			sequenceId = g.assignNewMetadataId(sequence)
			newMetadataType := NewMetadataType(sequenceId, sequence, NewMetadataTypeDefinitionSequence(sc.ToCompact(sequenceTypeId)))
			g.metadataTypes = append(g.metadataTypes, newMetadataType)
		}
		typeId = sequenceId
	}

	return typeId
}

func (g *MetadataTypeGenerator) constructTypeFields(v reflect.Value) sc.Sequence[MetadataTypeDefinitionField] {
	valueType := v.Type()
	typeNumFields := valueType.NumField()
	metadataFields := sc.Sequence[MetadataTypeDefinitionField]{}
	for i := 0; i < typeNumFields; i++ {
		fieldName := valueType.Field(i).Name
		fieldTypeName := valueType.Field(i).Type.Name()
		if isIgnoredName(fieldName) || isIgnoredType(fieldTypeName) {
			continue
		}
		fieldId, ok := g.metadataIds[fieldTypeName]
		if !ok {
			fieldId = g.BuildMetadataTypeRecursively(v.Field(i), nil, nil, nil)
		}
		if strings.HasPrefix(fieldTypeName, "Sequence") {
			fieldName = "Vec<" + fieldName + ">"
		}
		metadataFields = append(metadataFields, NewMetadataTypeDefinitionFieldWithName(fieldId, sc.Str(fieldName)))
	}
	return metadataFields
}

func (g *MetadataTypeGenerator) constructOptionType(v reflect.Value, typeId int) int {
	optionTypeName := v.FieldByName("Value").Type().Name()
	typeParameterId, _ := g.GetId(optionTypeName)
	metadataTypeParams := append(sc.Sequence[MetadataTypeParameter]{}, NewMetadataTypeParameter(typeParameterId, "T"))
	metadataTypeDef := optionTypeDefinition(optionTypeName, typeParameterId)
	metadataDocs := "Option<" + optionTypeName + ">"
	metadataTypePath := sc.Sequence[sc.Str]{"Option"}

	newMetadataType := NewMetadataTypeWithParams(typeId, metadataDocs, metadataTypePath, metadataTypeDef, metadataTypeParams)
	g.metadataTypes = append(g.metadataTypes, newMetadataType)
	return typeId
}

func (g *MetadataTypeGenerator) isCompactVariation(v reflect.Value) (int, bool) {
	field := v.FieldByName("Number")
	if field.IsValid() {
		if v.Type() == reflect.TypeOf(sc.Compact{}) {
			switch field.Elem().Type() {
			case reflect.TypeOf(*new(sc.U128)):
				typeId, ok := g.metadataIds["CompactU128"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU128")
					g.metadataTypes = append(g.metadataTypes, NewMetadataType(typeId, "CompactU128", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))))
				}
				return typeId, true
			case reflect.TypeOf(*new(sc.U64)):
				typeId, ok := g.metadataIds["CompactU64"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU64")
					g.metadataTypes = append(g.metadataTypes, NewMetadataType(typeId, "CompactU64", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))))
				}
				return typeId, true
			case reflect.TypeOf(*new(sc.U32)):
				typeId, ok := g.metadataIds["CompactU32"]
				if !ok {
					typeId = g.assignNewMetadataId("CompactU32")
					g.metadataTypes = append(g.metadataTypes, NewMetadataType(typeId, "CompactU32", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))))
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

func constructTypeDocs(typeName string) string {
	metadataDocs := strings.Replace(typeName, primitivesPackagePath, "", 1)
	return strings.Replace(metadataDocs, goscalePathTrim, "", 1)
}

// constructFunctionName constructs the formal name of a function call for the module metadata type given its struct name as an input (e.g. callTransferAll -> transfer_all)
func constructFunctionName(input string) string {
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

func isIgnoredType(t string) bool {
	return t == moduleTypeName || t == hookOnChargeTypeName || t == varyingDataTypeName
}

func isIgnoredName(name string) bool {
	return name == additionalSignedTypeName
}
