package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type MetadataExtrinsicV14 struct {
	Type             sc.Compact
	Version          sc.U8
	SignedExtensions sc.Sequence[MetadataSignedExtension]
}

func (me MetadataExtrinsicV14) Encode(buffer *bytes.Buffer) {
	me.Type.Encode(buffer)
	me.Version.Encode(buffer)
	me.SignedExtensions.Encode(buffer)
}

func DecodeMetadataExtrinsicV14(buffer *bytes.Buffer) MetadataExtrinsicV14 {
	return MetadataExtrinsicV14{
		Type:             sc.DecodeCompact(buffer),
		Version:          sc.DecodeU8(buffer),
		SignedExtensions: sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension),
	}
}

func (me MetadataExtrinsicV14) Bytes() []byte {
	return sc.EncodedBytes(me)
}

type MetadataExtrinsicV15 struct {
	Version          sc.U8
	Address          sc.Compact
	Call             sc.Compact
	Signature        sc.Compact
	Extra            sc.Compact
	SignedExtensions sc.Sequence[MetadataSignedExtension]
}

func (me MetadataExtrinsicV15) Encode(buffer *bytes.Buffer) {
	me.Version.Encode(buffer)
	me.Address.Encode(buffer)
	me.Call.Encode(buffer)
	me.Signature.Encode(buffer)
	me.Extra.Encode(buffer)
	me.SignedExtensions.Encode(buffer)
}

func DecodeMetadataExtrinsicV15(buffer *bytes.Buffer) MetadataExtrinsicV15 {
	return MetadataExtrinsicV15{
		Version:          sc.DecodeU8(buffer),
		Address:          sc.DecodeCompact(buffer),
		Call:             sc.DecodeCompact(buffer),
		Signature:        sc.DecodeCompact(buffer),
		Extra:            sc.DecodeCompact(buffer),
		SignedExtensions: sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension),
	}
}

func (me MetadataExtrinsicV15) Bytes() []byte {
	return sc.EncodedBytes(me)
}

type MetadataSignedExtension struct {
	Identifier       sc.Str
	Type             sc.Compact
	AdditionalSigned sc.Compact
}

func NewMetadataSignedExtension(identifier sc.Str, typeIndex, additionalSigned int) MetadataSignedExtension {
	return MetadataSignedExtension{
		Identifier:       identifier,
		Type:             sc.ToCompact(typeIndex),
		AdditionalSigned: sc.ToCompact(additionalSigned),
	}
}

func (mse MetadataSignedExtension) Encode(buffer *bytes.Buffer) {
	mse.Identifier.Encode(buffer)
	mse.Type.Encode(buffer)
	mse.AdditionalSigned.Encode(buffer)
}

func DecodeMetadataSignedExtension(buffer *bytes.Buffer) MetadataSignedExtension {
	return MetadataSignedExtension{
		Identifier:       sc.DecodeStr(buffer),
		Type:             sc.DecodeCompact(buffer),
		AdditionalSigned: sc.DecodeCompact(buffer),
	}
}

func (mse MetadataSignedExtension) Bytes() []byte {
	return sc.EncodedBytes(mse)
}
