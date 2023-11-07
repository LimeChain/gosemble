package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

const (
	MetadataTypeDefinitionComposite sc.U8 = iota
	MetadataTypeDefinitionVariant
	MetadataTypeDefinitionSequence
	MetadataTypeDefinitionFixedSequence
	MetadataTypeDefinitionTuple
	MetadataTypeDefinitionPrimitive
	MetadataTypeDefinitionCompact
	MetadataTypeDefinitionBitSequence
)

type MetadataTypeDefinition = sc.VaryingData

func NewMetadataTypeDefinitionComposite(fields sc.Sequence[MetadataTypeDefinitionField]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionComposite, fields)
}

func NewMetadataTypeDefinitionVariant(variants sc.Sequence[MetadataDefinitionVariant]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionVariant, variants)
}

func NewMetadataTypeDefinitionSequence(compact sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionSequence, compact)
}

func NewMetadataTypeDefinitionFixedSequence(length sc.U32, typeId sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionFixedSequence, length, typeId)
}

func NewMetadataTypeDefinitionTuple(compacts sc.Sequence[sc.Compact]) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionTuple, compacts)
}

func NewMetadataTypeDefinitionPrimitive(primitive MetadataDefinitionPrimitive) MetadataTypeDefinition {
	// TODO: type safety
	return sc.NewVaryingData(MetadataTypeDefinitionPrimitive, primitive)
}

func NewMetadataTypeDefinitionCompact(compact sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionCompact, compact)
}

func NewMetadataTypeDefinitionBitSequence(storeOrder, orderType sc.Compact) MetadataTypeDefinition {
	return sc.NewVaryingData(MetadataTypeDefinitionBitSequence, storeOrder, orderType)
}

func DecodeMetadataTypeDefinition(buffer *bytes.Buffer) (MetadataTypeDefinition, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataTypeDefinition{}, err
	}

	switch b {
	case MetadataTypeDefinitionComposite:
		fields, err := sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionField)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionComposite(fields), nil
	case MetadataTypeDefinitionVariant:
		variants, err := sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionVariant)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionVariant(variants), nil
	case MetadataTypeDefinitionSequence:
		cmpct, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionSequence(cmpct), nil
	case MetadataTypeDefinitionFixedSequence:
		len, err := sc.DecodeU32(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		id, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionFixedSequence(len, id), nil
	case MetadataTypeDefinitionTuple:
		cmpcts, err := sc.DecodeSequence[sc.Compact](buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionTuple(cmpcts), nil
	case MetadataTypeDefinitionPrimitive:
		prim, err := DecodeMetadataDefinitionPrimitive(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionPrimitive(prim), nil
	case MetadataTypeDefinitionCompact:
		cmpct, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionCompact(cmpct), nil
	case MetadataTypeDefinitionBitSequence:
		storeOrder, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		orderType, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MetadataTypeDefinition{}, err
		}
		return NewMetadataTypeDefinitionBitSequence(storeOrder, orderType), nil
	default:
		return MetadataTypeDefinition{}, newTypeError("MetadataTypeDefinition")
	}
}

type MetadataTypeDefinitionField struct {
	Name     sc.Option[sc.Str]
	Type     sc.Compact
	TypeName sc.Option[sc.Str]
	Docs     sc.Sequence[sc.Str]
}

func NewMetadataTypeDefinitionField(id int) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](nil),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](nil),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func NewMetadataTypeDefinitionFieldWithNames(id int, name sc.Str, idName sc.Str) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](name),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](idName),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func NewMetadataTypeDefinitionFieldWithName(id int, idName sc.Str) MetadataTypeDefinitionField {
	return MetadataTypeDefinitionField{
		Name:     sc.NewOption[sc.Str](nil),
		Type:     sc.ToCompact(id),
		TypeName: sc.NewOption[sc.Str](idName),
		Docs:     sc.Sequence[sc.Str]{},
	}
}

func (mtdf MetadataTypeDefinitionField) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		mtdf.Name,
		mtdf.Type,
		mtdf.TypeName,
		mtdf.Docs,
	)
}

func DecodeMetadataTypeDefinitionField(buffer *bytes.Buffer) (MetadataTypeDefinitionField, error) {
	name, err := sc.DecodeOption[sc.Str](buffer)
	if err != nil {
		return MetadataTypeDefinitionField{}, err
	}
	t, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataTypeDefinitionField{}, err
	}
	tName, err := sc.DecodeOption[sc.Str](buffer)
	if err != nil {
		return MetadataTypeDefinitionField{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataTypeDefinitionField{}, err
	}

	return MetadataTypeDefinitionField{
		Name:     name,
		Type:     t,
		TypeName: tName,
		Docs:     docs,
	}, nil
}

func (mtdf MetadataTypeDefinitionField) Bytes() []byte {
	return sc.EncodedBytes(mtdf)
}

type MetadataDefinitionVariant struct {
	Name   sc.Str
	Fields sc.Sequence[MetadataTypeDefinitionField]
	Index  sc.U8
	Docs   sc.Sequence[sc.Str]
}

func NewMetadataDefinitionVariant(name string, fields sc.Sequence[MetadataTypeDefinitionField], index sc.U8, docs string) MetadataDefinitionVariant {
	return NewMetadataDefinitionVariantStr(sc.Str(name), fields, index, docs)
}

func NewMetadataDefinitionVariantStr(name sc.Str, fields sc.Sequence[MetadataTypeDefinitionField], index sc.U8, docs string) MetadataDefinitionVariant {
	return MetadataDefinitionVariant{
		Name:   name,
		Fields: fields,
		Index:  index,
		Docs:   sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mdv MetadataDefinitionVariant) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		mdv.Name,
		mdv.Fields,
		mdv.Index,
		mdv.Docs,
	)
}

func DecodeMetadataTypeDefinitionVariant(buffer *bytes.Buffer) (MetadataDefinitionVariant, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataDefinitionVariant{}, err
	}
	fields, err := sc.DecodeSequenceWith(buffer, DecodeMetadataTypeDefinitionField)
	if err != nil {
		return MetadataDefinitionVariant{}, err
	}
	idx, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataDefinitionVariant{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataDefinitionVariant{}, err
	}
	return MetadataDefinitionVariant{
		Name:   name,
		Fields: fields,
		Index:  idx,
		Docs:   docs,
	}, nil
}

func (mdv MetadataDefinitionVariant) Bytes() []byte {
	return sc.EncodedBytes(mdv)
}

const (
	MetadataDefinitionPrimitiveBoolean MetadataDefinitionPrimitive = iota
	MetadataDefinitionPrimitiveChar
	MetadataDefinitionPrimitiveString
	MetadataDefinitionPrimitiveU8
	MetadataDefinitionPrimitiveU16
	MetadataDefinitionPrimitiveU32
	MetadataDefinitionPrimitiveU64
	MetadataDefinitionPrimitiveU128
	MetadataDefinitionPrimitiveU256
	MetadataDefinitionPrimitiveI8
	MetadataDefinitionPrimitiveI16
	MetadataDefinitionPrimitiveI32
	MetadataDefinitionPrimitiveI64
	MetadataDefinitionPrimitiveI128
	MetadataDefinitionPrimitiveI256
)

type MetadataDefinitionPrimitive = sc.U8

func DecodeMetadataDefinitionPrimitive(buffer *bytes.Buffer) (MetadataDefinitionPrimitive, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataDefinitionPrimitive(0), err
	}

	switch b {
	case MetadataDefinitionPrimitiveBoolean:
		return MetadataDefinitionPrimitiveBoolean, nil
	case MetadataDefinitionPrimitiveChar:
		return MetadataDefinitionPrimitiveChar, nil
	case MetadataDefinitionPrimitiveString:
		return MetadataDefinitionPrimitiveString, nil
	case MetadataDefinitionPrimitiveU8:
		return MetadataDefinitionPrimitiveU8, nil
	case MetadataDefinitionPrimitiveU16:
		return MetadataDefinitionPrimitiveU16, nil
	case MetadataDefinitionPrimitiveU32:
		return MetadataDefinitionPrimitiveU32, nil
	case MetadataDefinitionPrimitiveU64:
		return MetadataDefinitionPrimitiveU64, nil
	case MetadataDefinitionPrimitiveU128:
		return MetadataDefinitionPrimitiveU128, nil
	case MetadataDefinitionPrimitiveU256:
		return MetadataDefinitionPrimitiveU256, nil
	case MetadataDefinitionPrimitiveI8:
		return MetadataDefinitionPrimitiveI8, nil
	case MetadataDefinitionPrimitiveI16:
		return MetadataDefinitionPrimitiveI16, nil
	case MetadataDefinitionPrimitiveI32:
		return MetadataDefinitionPrimitiveI32, nil
	case MetadataDefinitionPrimitiveI64:
		return MetadataDefinitionPrimitiveI64, nil
	case MetadataDefinitionPrimitiveI128:
		return MetadataDefinitionPrimitiveI128, nil
	case MetadataDefinitionPrimitiveI256:
		return MetadataDefinitionPrimitiveI256, nil
	default:
		return MetadataDefinitionPrimitive(0), newTypeError("MetadataDefinitionPrimitive")
	}
}
