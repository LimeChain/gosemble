package types

import (
	"reflect"
	"strings"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

type MetadataGenerator interface {
	assignNewMetadataId(name string) (id int)
	BuildMetadataTypeRecursively(t reflect.Value) int
	IdsMap() map[string]int
	constructFunctionName(input string) string
	CallsMetadata(moduleName string, moduleFunctions map[sc.U8]Call, params *sc.Sequence[MetadataTypeParameter]) (MetadataType, int)
}

type MetadataTypeGenerator struct {
	MetadataIds   map[string]int
	MetadataTypes sc.Sequence[MetadataType]
}

func NewMetadataTypeGenerator() MetadataGenerator {
	return &MetadataTypeGenerator{
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
		"CompactU128":  metadata.TypesCompactU128,
	}
}

func (g MetadataTypeGenerator) IdsMap() map[string]int {
	return g.MetadataIds
}

func (g MetadataTypeGenerator) assignNewMetadataId(name string) int {
	lastIndex := len(g.MetadataIds)
	newId := lastIndex + 1
	g.MetadataIds[name] = newId
	return newId
}

// BuildMetadataTypeRecursively Builds the metadata type (recursively) if it does not exist
func (g MetadataTypeGenerator) BuildMetadataTypeRecursively(v reflect.Value) int {
	valueType := v.Type()
	typeName := valueType.Name()
	var typeId int
	var ok bool
	switch valueType.Kind() {
	case reflect.Struct:
		typeId, ok = g.MetadataIds[typeName]
		if !ok {
			typeId = g.assignNewMetadataId(typeName)
			typeNumFields := valueType.NumField()
			metadataFields := sc.Sequence[MetadataTypeDefinitionField]{}
			for i := 0; i < typeNumFields; i++ {
				fieldName := valueType.Field(i).Name
				fieldTypeName := valueType.Field(i).Type.Name()
				fieldId, ok := g.MetadataIds[fieldTypeName]
				if !ok {
					fieldId = g.BuildMetadataTypeRecursively(v.Field(i))
				}
				metadataFields = append(metadataFields, NewMetadataTypeDefinitionFieldWithName(fieldId, sc.Str(fieldName)))
			}

			g.MetadataTypes = append(g.MetadataTypes,
				NewMetadataType(
					typeId,
					typeName,
					NewMetadataTypeDefinitionComposite(metadataFields)))
		}
	case reflect.Slice:
		sequenceName := "Sequence"
		sequenceType := sequenceName + valueType.Elem().Name()
		sequenceTypeId, ok := g.MetadataIds[sequenceType]
		if !ok {
			sequenceTypeId = g.BuildMetadataTypeRecursively(v.Elem())
		}
		typeId = sequenceTypeId
	case reflect.Array: // Compact type
		compactLen := v.Len()
		switch compactLen {
		case 2: // CompactU128
			if valueType.Name() == "Compact" {
				typeId, ok = g.MetadataIds["CompactU128"]
				if !ok {
					typeId = g.BuildMetadataTypeRecursively(v.Elem())
				}
			} else {
				typeId = g.MetadataIds[typeName]
			}
		}
	case reflect.Uint64:
		if strings.HasPrefix(typeName, "Compact") {
			typeId, ok = g.MetadataIds["CompactU64"]
			if !ok {
				typeId = g.assignNewMetadataId(typeName)
				g.MetadataTypes = append(g.MetadataTypes, NewMetadataType(typeId, "CompactU64", NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))))
			}
		}
	}
	return typeId
}

// constructFunctionName constructs the formal name of a function call for the module metadata type given its struct name as an input (e.g. callTransferAll -> transfer_all)
func (g MetadataTypeGenerator) constructFunctionName(input string) string {
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
func (g MetadataTypeGenerator) CallsMetadata(moduleName string, moduleFunctions map[sc.U8]Call, params *sc.Sequence[MetadataTypeParameter]) (MetadataType, int) {
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
				currentArgId := g.BuildMetadataTypeRecursively(currentArg)
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
