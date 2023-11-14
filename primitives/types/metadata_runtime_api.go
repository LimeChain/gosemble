package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeApiMetadata struct {
	Name    sc.Str
	Methods sc.Sequence[RuntimeApiMethodMetadata]
	Docs    sc.Sequence[sc.Str]
}

func (ram RuntimeApiMetadata) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		ram.Name,
		ram.Methods,
		ram.Docs,
	)
}

func DecodeRuntimeApiMetadata(buffer *bytes.Buffer) (RuntimeApiMetadata, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	methods, err := sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodMetadata)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	return RuntimeApiMetadata{
		Name:    name,
		Methods: methods,
		Docs:    docs,
	}, nil
}

func (ram RuntimeApiMetadata) Bytes() []byte {
	return sc.EncodedBytes(ram)
}
