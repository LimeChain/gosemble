package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeMetadataV15 struct {
	Types      sc.Sequence[MetadataType]
	Modules    sc.Sequence[MetadataModuleV15]
	Extrinsic  MetadataExtrinsicV15
	Type       sc.Compact
	Apis       sc.Sequence[RuntimeApiMetadata]
	OuterEnums OuterEnums
	Custom     CustomMetadata
}

func (rm RuntimeMetadataV15) Encode(buffer *bytes.Buffer) {
	rm.Types.Encode(buffer)
	rm.Modules.Encode(buffer)
	rm.Extrinsic.Encode(buffer)
	rm.Type.Encode(buffer)
	rm.Apis.Encode(buffer)
	rm.OuterEnums.Encode(buffer)
	rm.Custom.Encode(buffer)
}

func DecodeRuntimeMetadataV15(buffer *bytes.Buffer) (RuntimeMetadataV15, error) {
	types, err := sc.DecodeSequenceWith(buffer, DecodeMetadataType)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	modules, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleV15)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	extrinsic, err := DecodeMetadataExtrinsicV15(buffer)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	typeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	apis, err := sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMetadata)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	outerEnums, err := DecodeOuterEnums(buffer)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}
	customMd, err := DecodeCustomMetadata(buffer)
	if err != nil {
		return RuntimeMetadataV15{}, err
	}

	return RuntimeMetadataV15{
		Types:      types,
		Modules:    modules,
		Extrinsic:  extrinsic,
		Type:       typeId,
		Apis:       apis,
		OuterEnums: outerEnums,
		Custom:     customMd,
	}, nil
}

func (rm RuntimeMetadataV15) Bytes() []byte {
	return sc.EncodedBytes(rm)
}

type CustomValueMetadata struct {
	Type  sc.Compact
	Value sc.Sequence[sc.U8]
}

func (cvm CustomValueMetadata) Encode(buffer *bytes.Buffer) {
	cvm.Type.Encode(buffer)
	cvm.Value.Encode(buffer)
}

func (cvm CustomValueMetadata) Bytes() []byte {
	return sc.EncodedBytes(cvm)
}

type CustomMetadata struct {
	Map sc.Dictionary[sc.Str, CustomValueMetadata]
}

func (cm CustomMetadata) Encode(buffer *bytes.Buffer) {
	cm.Map.Encode(buffer)
}

func DecodeCustomMetadata(buffer *bytes.Buffer) (CustomMetadata, error) {
	m, err := sc.DecodeDictionary[sc.Str, CustomValueMetadata](buffer)
	if err != nil {
		return CustomMetadata{}, err
	}
	return CustomMetadata{
		Map: m,
	}, nil
}

func (cm CustomMetadata) Bytes() []byte {
	return sc.EncodedBytes(cm)
}

type OuterEnums struct {
	CallEnumType  sc.Compact
	EventEnumType sc.Compact
	ErrorEnumType sc.Compact
}

func (oe OuterEnums) Encode(buffer *bytes.Buffer) {
	oe.CallEnumType.Encode(buffer)
	oe.EventEnumType.Encode(buffer)
	oe.ErrorEnumType.Encode(buffer)
}

func DecodeOuterEnums(buffer *bytes.Buffer) (OuterEnums, error) {
	callEnum, err := sc.DecodeCompact(buffer)
	if err != nil {
		return OuterEnums{}, err
	}
	eventEnum, err := sc.DecodeCompact(buffer)
	if err != nil {
		return OuterEnums{}, err
	}
	errorEnum, err := sc.DecodeCompact(buffer)
	if err != nil {
		return OuterEnums{}, err
	}
	return OuterEnums{
		CallEnumType:  callEnum,
		EventEnumType: eventEnum,
		ErrorEnumType: errorEnum,
	}, nil
}

func (oe OuterEnums) Bytes() []byte {
	return sc.EncodedBytes(oe)
}

type RuntimeApiMethodParamMetadata struct {
	Name sc.Str
	Type sc.Compact
}

func (rampm RuntimeApiMethodParamMetadata) Encode(buffer *bytes.Buffer) {
	rampm.Name.Encode(buffer)
	rampm.Type.Encode(buffer)
}

func DecodeRuntimeApiMethodParamMetadata(buffer *bytes.Buffer) (RuntimeApiMethodParamMetadata, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeApiMethodParamMetadata{}, err
	}
	typeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return RuntimeApiMethodParamMetadata{}, err
	}
	return RuntimeApiMethodParamMetadata{
		Name: name,
		Type: typeId,
	}, nil
}

func (rampm RuntimeApiMethodParamMetadata) Bytes() []byte {
	return sc.EncodedBytes(rampm)
}

type RuntimeApiMethodMetadata struct {
	Name   sc.Str
	Inputs sc.Sequence[RuntimeApiMethodParamMetadata]
	Output sc.Compact
	Docs   sc.Sequence[sc.Str]
}

func (ramm RuntimeApiMethodMetadata) Encode(buffer *bytes.Buffer) {
	ramm.Name.Encode(buffer)
	ramm.Inputs.Encode(buffer)
	ramm.Output.Encode(buffer)
	ramm.Docs.Encode(buffer)
}

func DecodeRuntimeApiMethodMetadata(buffer *bytes.Buffer) (RuntimeApiMethodMetadata, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeApiMethodMetadata{}, err
	}
	inputs, err := sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodParamMetadata)
	if err != nil {
		return RuntimeApiMethodMetadata{}, err
	}
	output, err := sc.DecodeCompact(buffer)
	if err != nil {
		return RuntimeApiMethodMetadata{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return RuntimeApiMethodMetadata{}, err
	}
	return RuntimeApiMethodMetadata{
		Name:   name,
		Inputs: inputs,
		Output: output,
		Docs:   docs,
	}, nil
}

func (ramm RuntimeApiMethodMetadata) Bytes() []byte {
	return sc.EncodedBytes(ramm)
}
