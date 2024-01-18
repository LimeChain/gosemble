package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeMetadataV14 struct {
	Types     sc.Sequence[MetadataType]
	Modules   sc.Sequence[MetadataModuleV14]
	Extrinsic MetadataExtrinsicV14
	Type      sc.Compact
}

func (rm RuntimeMetadataV14) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		rm.Types,
		rm.Modules,
		rm.Extrinsic,
		rm.Type,
	)
}

func DecodeRuntimeMetadataV14(buffer *bytes.Buffer) (RuntimeMetadataV14, error) {
	types, err := sc.DecodeSequenceWith(buffer, DecodeMetadataType)
	if err != nil {
		return RuntimeMetadataV14{}, err
	}
	modules, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleV14)
	if err != nil {
		return RuntimeMetadataV14{}, err
	}
	extrinsic, err := DecodeMetadataExtrinsicV14(buffer)
	if err != nil {
		return RuntimeMetadataV14{}, err
	}
	typeId, err := sc.DecodeCompact[sc.U128](buffer)
	if err != nil {
		return RuntimeMetadataV14{}, err
	}
	return RuntimeMetadataV14{
		Types:     types,
		Modules:   modules,
		Extrinsic: extrinsic,
		Type:      typeId,
	}, nil
}

func (rm RuntimeMetadataV14) Bytes() []byte {
	return sc.EncodedBytes(rm)
}
