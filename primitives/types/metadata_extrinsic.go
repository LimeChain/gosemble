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

func (me MetadataExtrinsicV14) Encode(buffer *bytes.Buffer) error {
	err := me.Type.Encode(buffer)
	if err != nil {
		return err
	}
	err = me.Version.Encode(buffer)
	if err != nil {
		return err
	}
	return me.SignedExtensions.Encode(buffer)
}

func DecodeMetadataExtrinsicV14(buffer *bytes.Buffer) (MetadataExtrinsicV14, error) {
	typeId, err := sc.DecodeCompact(buffer)
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
		Type:             typeId,
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

func (me MetadataExtrinsicV15) Encode(buffer *bytes.Buffer) error {
	err := me.Version.Encode(buffer)
	if err != nil {
		return err
	}
	err = me.Address.Encode(buffer)
	if err != nil {
		return err
	}
	err = me.Call.Encode(buffer)
	if err != nil {
		return err
	}
	err = me.Signature.Encode(buffer)
	if err != nil {
		return err
	}
	err = me.Extra.Encode(buffer)
	if err != nil {
		return err
	}
	return me.SignedExtensions.Encode(buffer)
}

func DecodeMetadataExtrinsicV15(buffer *bytes.Buffer) (MetadataExtrinsicV15, error) {
	version, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	addrTypeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	callTypeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	sigTypeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	extraTypeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}
	seTypeId, err := sc.DecodeSequenceWith(buffer, DecodeMetadataSignedExtension)
	if err != nil {
		return MetadataExtrinsicV15{}, err
	}

	return MetadataExtrinsicV15{
		Version:          version,
		Address:          addrTypeId,
		Call:             callTypeId,
		Signature:        sigTypeId,
		Extra:            extraTypeId,
		SignedExtensions: seTypeId,
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

func (mse MetadataSignedExtension) Encode(buffer *bytes.Buffer) error {
	err := mse.Identifier.Encode(buffer)
	if err != nil {
		return err
	}
	err = mse.Type.Encode(buffer)
	if err != nil {
		return err
	}
	return mse.AdditionalSigned.Encode(buffer)
}

func DecodeMetadataSignedExtension(buffer *bytes.Buffer) (MetadataSignedExtension, error) {
	id, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	typeId, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	as, err := sc.DecodeCompact(buffer)
	if err != nil {
		return MetadataSignedExtension{}, err
	}
	return MetadataSignedExtension{
		Identifier:       id,
		Type:             typeId,
		AdditionalSigned: as,
	}, nil
}

func (mse MetadataSignedExtension) Bytes() []byte {
	return sc.EncodedBytes(mse)
}
