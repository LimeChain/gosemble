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

func (rm RuntimeMetadataV14) Encode(buffer *bytes.Buffer) {
	rm.Types.Encode(buffer)
	rm.Modules.Encode(buffer)
	rm.Extrinsic.Encode(buffer)
	rm.Type.Encode(buffer)
}

func DecodeRuntimeMetadataV14(buffer *bytes.Buffer) RuntimeMetadataV14 {
	return RuntimeMetadataV14{
		Types:     sc.DecodeSequenceWith(buffer, DecodeMetadataType),
		Modules:   sc.DecodeSequenceWith(buffer, DecodeMetadataModuleV14),
		Extrinsic: DecodeMetadataExtrinsicV14(buffer),
		Type:      sc.DecodeCompact(buffer),
	}
}

func (rm RuntimeMetadataV14) Bytes() []byte {
	return sc.EncodedBytes(rm)
}
