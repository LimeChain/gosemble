package types

import (
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

type MetadataGenerator interface {
	AssignNewMetadataId(name string) (id int)
	BuildMetadataTypeRecursively(t reflect.Type) int
	GetMap() map[string]int
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
		"Bool":       metadata.PrimitiveTypesBool,
		"String":     metadata.PrimitiveTypesString,
		"U8":         metadata.PrimitiveTypesU8,
		"U16":        metadata.PrimitiveTypesU16,
		"U32":        metadata.PrimitiveTypesU32,
		"U64":        metadata.PrimitiveTypesU64,
		"U128":       metadata.PrimitiveTypesU128,
		"I8":         metadata.PrimitiveTypesI8,
		"I16":        metadata.PrimitiveTypesI16,
		"I32":        metadata.PrimitiveTypesI32,
		"I64":        metadata.PrimitiveTypesI64,
		"I128":       metadata.PrimitiveTypesI128,
		"H256":       metadata.TypesH256,
		"SequenceU8": metadata.TypesSequenceU8,
	}
}

func (g MetadataTypeGenerator) GetMap() map[string]int {
	return g.MetadataIds
}

func (g MetadataTypeGenerator) AssignNewMetadataId(name string) int {
	lastIndex := len(g.MetadataIds)
	newId := lastIndex + 1
	g.MetadataIds[name] = newId
	return newId
}

// BuildMetadataTypeRecursively Builds the metadata type (recursively) if it does not exist
func (g MetadataTypeGenerator) BuildMetadataTypeRecursively(t reflect.Type) int {
	typeName := t.Name()
	typeKind := t.Kind()
	var typeId int
	var ok bool
	switch typeKind {
	case reflect.Struct:
		typeId, ok = g.MetadataIds[typeName]
		if !ok {
			typeId = g.AssignNewMetadataId(typeName)
			typeNumFields := t.NumField()
			metadataFields := sc.Sequence[MetadataTypeDefinitionField]{}
			for i := 0; i < typeNumFields; i++ {
				fieldName := t.Field(i).Name
				fieldTypeName := t.Field(i).Type.Name()
				fieldId, ok := g.MetadataIds[fieldTypeName]
				if !ok {
					fieldId = g.BuildMetadataTypeRecursively(t.Field(i).Type)
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
		sequenceType := sequenceName + t.Elem().Name()
		sequenceTypeId, ok := g.MetadataIds[sequenceType]
		if !ok {
			sequenceTypeId = g.BuildMetadataTypeRecursively(t.Elem())
		}
		typeId = sequenceTypeId
	}
	return typeId
}
