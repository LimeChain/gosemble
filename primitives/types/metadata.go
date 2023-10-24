package types

import (
	"bytes"
	"errors"
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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

func DecodeMetadata(buffer *bytes.Buffer) (Metadata, error) {
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	metaReserved, err := sc.DecodeU32(buffer)
	if err != nil {
		return Metadata{}, err
	}
	if metaReserved != MetadataReserved {
		log.Critical("metadata reserved mismatch: expect [" + strconv.Itoa(int(MetadataReserved)) + "], actual [" + strconv.Itoa(int(metaReserved)) + "]")
		return Metadata{}, errors.New("metadata reserved mismatch: expect [" + strconv.Itoa(int(MetadataReserved)) + "], actual [" + strconv.Itoa(int(metaReserved)) + "]")
	}

	version, err := sc.DecodeU8(buffer)
	if err != nil {
		return Metadata{}, err
	}

	switch version {
	case MetadataVersion14:
		data14, err := DecodeRuntimeMetadataV14(buffer)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{Version: MetadataVersion14, DataV14: data14}, nil
	case MetadataVersion15:
		data15, err := DecodeRuntimeMetadataV15(buffer)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{Version: MetadataVersion15, DataV15: data15}, nil
	default:
		log.Critical("metadata version mismatch: expect [" + strconv.Itoa(int(MetadataVersion14)) + "or" + strconv.Itoa(int(MetadataVersion15)) + "] , actual [" + strconv.Itoa(int(version)) + "]")
		return Metadata{}, errors.New("metadata version mismatch: expect [" + strconv.Itoa(int(MetadataVersion14)) + "or" + strconv.Itoa(int(MetadataVersion15)) + "] , actual [" + strconv.Itoa(int(version)) + "]")
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

func DecodeMetadataType(buffer *bytes.Buffer) (MetadataType, error) {
	id, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataType{}, err
	}
	path, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataType{}, err
	}
	params, err := sc.DecodeSequenceWith(buffer, DecodeMetadataTypeParameter)
	if err != nil {
		return MetadataType{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataType{}, err
	}
	def, err := DecodeMetadataTypeDefinition(buffer)
	if err != nil {
		return MetadataType{}, err
	}
	return MetadataType{
		Id:         id,
		Path:       path,
		Params:     params,
		Definition: def,
		Docs:       docs,
	}, nil
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

func DecodeMetadataTypeParameter(buffer *bytes.Buffer) (MetadataTypeParameter, error) {
	text, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataTypeParameter{}, err
	}
	t, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataTypeParameter{}, err
	}
	return MetadataTypeParameter{
		Text: text,
		Type: t,
	}, nil
}

func (mtp MetadataTypeParameter) Bytes() []byte {
	return sc.EncodedBytes(mtp)
}
