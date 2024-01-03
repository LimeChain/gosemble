package types

import (
	"reflect"
	"strings"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

const (
	additionalSignedTypeName = "typesInfoAdditionalSignedData"
	moduleTypeName           = "Module"
	hookOnChargeTypeName     = "OnChargeTransaction"
)

type MetadataTypeGenerator struct {
	MetadataIds   map[string]int
	MetadataTypes sc.Sequence[MetadataType]
}

func NewMetadataTypeGenerator() MetadataTypeGenerator {
	return MetadataTypeGenerator{
		MetadataIds:   BuildMetadataTypesIdsMap(),
		MetadataTypes: sc.Sequence[MetadataType]{},
	}
}

func BuildMetadataTypesIdsMap() map[string]int {
	return map[string]int{
		"Bool":         metadata.PrimitiveTypesBool,
		"String":       metadata.PrimitiveTypesString,
		"U8":           metadata.PrimitiveTypesU8,
		"U16":          metadata.PrimitiveTypesU16,
		"U32":          metadata.PrimitiveTypesU32,
		"U64":          metadata.PrimitiveTypesU64,
		"U128":         metadata.PrimitiveTypesU128,
		"I8":           metadata.PrimitiveTypesI8,
		"I16":          metadata.PrimitiveTypesI16,
		"I32":          metadata.PrimitiveTypesI32,
		"I64":          metadata.PrimitiveTypesI64,
		"I128":         metadata.PrimitiveTypesI128,
		"H256":         metadata.TypesH256,
		"SequenceU8":   metadata.TypesSequenceU8,
		"MultiAddress": metadata.TypesMultiAddress,
	}
}

func (g *MetadataTypeGenerator) IdsMap() map[string]int {
	return g.MetadataIds
}

func (g *MetadataTypeGenerator) GetMetadataTypes() sc.Sequence[MetadataType] {
	return g.MetadataTypes
}

func (g *MetadataTypeGenerator) AppendMetadataTypes(types sc.Sequence[MetadataType]) {
	g.MetadataTypes = append(g.MetadataTypes, types...)
}

func (g *MetadataTypeGenerator) assignNewMetadataId(name string) int {
	newId := len(g.MetadataIds) + 1
	g.MetadataIds[name] = newId
	return newId
}

func (g *MetadataTypeGenerator) isCompactVariation(t reflect.Type) (int, bool) {
	switch t {
	case reflect.TypeOf(*new(sc.Compact[sc.U128])):
		typeId := g.assignNewMetadataId(t.Name())
		g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU128", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))))
		return typeId, true
	case reflect.TypeOf(*new(sc.Compact[sc.U64])):
		typeId := g.assignNewMetadataId(t.Name())
		g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU64", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))))
		return typeId, true
	}
	return -1, false
}

// BuildMetadataTypeRecursively Builds the metadata type (recursively) if it does not exist
func (g *MetadataTypeGenerator) BuildMetadataTypeRecursively(v reflect.Value, path *sc.Sequence[sc.Str]) int {
	valueType := v.Type()
	typeName := valueType.Name()
	var typeId int
	var ok bool
	switch valueType.Kind() {
	case reflect.Struct:
		typeId, ok = g.MetadataIds[typeName]
		if !ok {
			typeId, ok = g.isCompactVariation(valueType)
			if ok {
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
				fieldId, ok := g.MetadataIds[fieldTypeName]
				if !ok {
					fieldId = g.BuildMetadataTypeRecursively(v.Field(i), nil)
				}
				metadataFields = append(metadataFields, NewMetadataTypeDefinitionFieldWithName(fieldId, sc.Str(fieldName)))
			}
			var newMetadataType MetadataType
			if path != nil {
				newMetadataType = NewMetadataTypeWithPath(
					typeId,
					typeName,
					*path,
					NewMetadataTypeDefinitionComposite(metadataFields))
			} else {
				newMetadataType = NewMetadataType(
					typeId,
					typeName,
					NewMetadataTypeDefinitionComposite(metadataFields))
			}
			g.MetadataTypes = append(g.MetadataTypes, newMetadataType)
		}
	case reflect.Slice:
		sequenceName := "Sequence"
		sequenceType := sequenceName + valueType.Elem().Name()
		sequenceTypeId, ok := g.MetadataIds[sequenceType]
		if !ok {
			sequenceTypeId = g.BuildMetadataTypeRecursively(v.Elem(), path)
		}
		typeId = sequenceTypeId
	case reflect.Array: // types U128 and U64
		typeId = g.MetadataIds[typeName]
	}
	return typeId
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

// CallsMetadata returns metadata calls type of a module
func (g *MetadataTypeGenerator) CallsMetadata(moduleName string, moduleFunctions map[sc.U8]Call, params *sc.Sequence[MetadataTypeParameter]) (MetadataType, int) {
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
				currentArgId := g.BuildMetadataTypeRecursively(currentArg, nil)
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

	return NewMetadataTypeWithParams(balancesCallsMetadataId, moduleName+" calls", sc.Sequence[sc.Str]{sc.Str("pallet_" + strings.ToLower(moduleName)), "pallet", "Call"}, variant, *params), balancesCallsMetadataId
}

func isIgnoredType(t string) bool {
	return t == moduleTypeName || t == hookOnChargeTypeName
}

func isIgnoredName(name string) bool {
	return name == additionalSignedTypeName
}
