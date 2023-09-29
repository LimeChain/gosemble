package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

type RuntimeMetadataV15 struct {
	Types      sc.Sequence[MetadataType]
	Modules    sc.Sequence[MetadataModule]
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

func DecodeRuntimeMetadataV15(buffer *bytes.Buffer) RuntimeMetadataV15 {
	return RuntimeMetadataV15{
		Types:      sc.DecodeSequenceWith(buffer, DecodeMetadataType),
		Modules:    sc.DecodeSequenceWith(buffer, DecodeMetadataModule),
		Extrinsic:  DecodeMetadataExtrinsicV15(buffer),
		Type:       sc.DecodeCompact(buffer),
		Apis:       sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMetadata),
		OuterEnums: DecodeOuterEnums(buffer),
		Custom:     DecodeCustomMetadata(buffer),
	}
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

func DecodeCustomMetadata(buffer *bytes.Buffer) CustomMetadata {
	return CustomMetadata{
		Map: sc.DecodeDictionary[sc.Str, CustomValueMetadata](buffer),
	}
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

func DecodeOuterEnums(buffer *bytes.Buffer) OuterEnums {
	return OuterEnums{
		CallEnumType:  sc.DecodeCompact(buffer),
		EventEnumType: sc.DecodeCompact(buffer),
		ErrorEnumType: sc.DecodeCompact(buffer),
	}
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

func DecodeRuntimeApiMethodParamMetadata(buffer *bytes.Buffer) RuntimeApiMethodParamMetadata {
	return RuntimeApiMethodParamMetadata{
		Name: sc.DecodeStr(buffer),
		Type: sc.DecodeCompact(buffer),
	}
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

func DecodeRuntimeApiMethodMetadata(buffer *bytes.Buffer) RuntimeApiMethodMetadata {
	return RuntimeApiMethodMetadata{
		Name:   sc.DecodeStr(buffer),
		Inputs: sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodParamMetadata),
		Output: sc.DecodeCompact(buffer),
		Docs:   sc.DecodeSequence[sc.Str](buffer),
	}
}

func (ramm RuntimeApiMethodMetadata) Bytes() []byte {
	return sc.EncodedBytes(ramm)
}

type RuntimeApiMetadata struct {
	Name    sc.Str
	Methods sc.Sequence[RuntimeApiMethodMetadata]
	Docs    sc.Sequence[sc.Str]
}

func (ram RuntimeApiMetadata) Encode(buffer *bytes.Buffer) {
	ram.Name.Encode(buffer)
	ram.Methods.Encode(buffer)
	ram.Docs.Encode(buffer)
}

func DecodeRuntimeApiMetadata(buffer *bytes.Buffer) RuntimeApiMetadata {
	return RuntimeApiMetadata{
		Name:    sc.DecodeStr(buffer),
		Methods: sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodMetadata),
		Docs:    sc.DecodeSequence[sc.Str](buffer),
	}
}

func (ram RuntimeApiMetadata) Bytes() []byte {
	return sc.EncodedBytes(ram)
}
