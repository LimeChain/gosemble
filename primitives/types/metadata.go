package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"strconv"
)

const (
	MetadataReserved  sc.U32 = 0x6174656d // "meta"
	MetadataVersion14 sc.U8  = 14
	MetadataVersion15 sc.U8  = 15
)

type Metadata15 struct {
	Data RuntimeMetadataV15
}

func (m Metadata15) Bytes() []byte {
	return sc.EncodedBytes(m)
}

func (m Metadata15) Encode(buffer *bytes.Buffer) {
	MetadataReserved.Encode(buffer)
	MetadataVersion15.Encode(buffer)
	m.Data.Encode(buffer)
}

type Metadata14 struct {
	Data RuntimeMetadataV14
}

func (m Metadata14) Bytes() []byte {
	return sc.EncodedBytes(m)
}

func (m Metadata14) Encode(buffer *bytes.Buffer) {
	MetadataReserved.Encode(buffer)
	MetadataVersion14.Encode(buffer)
	m.Data.Encode(buffer)
}

type Metadata struct {
	Version sc.U8
	DataV14 RuntimeMetadataV14
	DataV15 RuntimeMetadataV15
}

func NewMetadataV14(data RuntimeMetadataV14) Metadata14 {
	return Metadata14{Data: data}
}

func NewMetadataV15(data RuntimeMetadataV15) Metadata15 {
	return Metadata15{Data: data}
}

func (m Metadata) Encode(buffer *bytes.Buffer) {
	MetadataReserved.Encode(buffer)

	switch m.Version {
	case MetadataVersion14:
		MetadataVersion14.Encode(buffer)
		m.DataV14.Encode(buffer)
	case MetadataVersion15:
		MetadataVersion15.Encode(buffer)
		m.DataV15.Encode(buffer)
	default:
		log.Critical("Unsupported metadata version")
	}
}

func DecodeMetadata(buffer *bytes.Buffer) Metadata {
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	metaReserved := sc.DecodeU32(buffer)
	if metaReserved != MetadataReserved {
		log.Critical("metadata reserved mismatch: expect [" + strconv.Itoa(int(MetadataReserved)) + "], actual [" + strconv.Itoa(int(metaReserved)) + "]")
		return Metadata{}
	}

	version := sc.DecodeU8(buffer)

	switch version {
	case MetadataVersion14:
		return Metadata{Version: MetadataVersion14, DataV14: DecodeRuntimeMetadataV14(buffer)}
	case MetadataVersion15:
		return Metadata{Version: MetadataVersion15, DataV15: DecodeRuntimeMetadataV15(buffer)}
	default:
		log.Critical("metadata version mismatch: expect [" + strconv.Itoa(int(MetadataVersion14)) + "or" + strconv.Itoa(int(MetadataVersion15)) + "] , actual [" + strconv.Itoa(int(version)) + "]")
		return Metadata{}
	}
}

func (m Metadata) Bytes() []byte {
	return sc.EncodedBytes(m)
}

type MetadataType struct {
	Id         sc.Compact
	Path       sc.Sequence[sc.Str]
	Params     sc.Sequence[MetadataTypeParameter]
	Definition MetadataTypeDefinition
	Docs       sc.Sequence[sc.Str]
}

func NewMetadataType(id int, docs string, definition MetadataTypeDefinition) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithPath(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       path,
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithParam(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition, param MetadataTypeParameter) MetadataType {
	return MetadataType{
		Id:   sc.ToCompact(id),
		Path: path,
		Params: sc.Sequence[MetadataTypeParameter]{
			param,
		},
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func NewMetadataTypeWithParams(id int, docs string, path sc.Sequence[sc.Str], definition MetadataTypeDefinition, params sc.Sequence[MetadataTypeParameter]) MetadataType {
	return MetadataType{
		Id:         sc.ToCompact(id),
		Path:       path,
		Params:     params,
		Definition: definition,
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mt MetadataType) Encode(buffer *bytes.Buffer) {
	mt.Id.Encode(buffer)
	mt.Path.Encode(buffer)
	mt.Params.Encode(buffer)
	mt.Definition.Encode(buffer)
	mt.Docs.Encode(buffer)
}

func DecodeMetadataType(buffer *bytes.Buffer) MetadataType {
	return MetadataType{
		Id:         sc.DecodeCompact(buffer),
		Path:       sc.DecodeSequence[sc.Str](buffer),
		Params:     sc.DecodeSequenceWith(buffer, DecodeMetadataTypeParameter),
		Definition: DecodeMetadataTypeDefinition(buffer),
		Docs:       sc.DecodeSequence[sc.Str](buffer),
	}
}

func (mt MetadataType) Bytes() []byte {
	return sc.EncodedBytes(mt)
}

type MetadataTypeParameter struct {
	Text sc.Str
	Type sc.Option[sc.Compact]
}

func NewMetadataTypeParameter(id int, text string) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.Str(text),
		Type: sc.NewOption[sc.Compact](sc.ToCompact(id)),
	}
}

func NewMetadataTypeParameterCompactId(id sc.Compact, text string) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.Str(text),
		Type: sc.NewOption[sc.Compact](id),
	}
}

func NewMetadataEmptyTypeParameter(text string) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.Str(text),
		Type: sc.NewOption[sc.Compact](nil),
	}
}

func (mtp MetadataTypeParameter) Encode(buffer *bytes.Buffer) {
	mtp.Text.Encode(buffer)
	mtp.Type.Encode(buffer)
}

func DecodeMetadataTypeParameter(buffer *bytes.Buffer) MetadataTypeParameter {
	return MetadataTypeParameter{
		Text: sc.DecodeStr(buffer),
		Type: sc.DecodeOption[sc.Compact](buffer),
	}
}

func (mtp MetadataTypeParameter) Bytes() []byte {
	return sc.EncodedBytes(mtp)
}
