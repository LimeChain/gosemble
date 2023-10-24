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

func DecodeMetadataExtrinsicV14(buffer *bytes.Buffer) (MetadataExtrinsicV14, error) {
	t, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV14{}, err
	}
	version, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataExtrinsicV14{}, err
	}
	se, err := sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension)
	if err != nil {
		return MetadataExtrinsicV14{}, err
	}

	return MetadataExtrinsicV14{
		Type:             t,
		Version:          version,
		SignedExtensions: se,
	}, nil
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

func DecodeMetadataExtrinsicV15(buffer *bytes.Buffer) (MetadataExtrinsicV15, error) {
	version, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	addr, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	call, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	sig, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	extra, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	se, err := sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}

	return MetadataExtrinsicV15{
		Version:          version,
		Address:          addr,
		Call:             call,
		Signature:        sig,
		Extra:            extra,
		SignedExtensions: se,
	}, nil
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

func DecodeMetadataSignedExtension(buffer *bytes.Buffer) (MetadataSignedExtension, error) {
	id, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	t, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	as, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	return MetadataSignedExtension{
		Identifier:       id,
		Type:             t,
		AdditionalSigned: as,
	}, nil
}

func (mse MetadataSignedExtension) Bytes() []byte {
	return sc.EncodedBytes(mse)
}
