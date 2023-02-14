package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type ExtrinsicEra struct {
	// TODO:
}

func (e ExtrinsicEra) Encode(buffer *bytes.Buffer) {
	// TODO:
}

func DecodeExtrinsicEra(buffer *bytes.Buffer) ExtrinsicEra {
	// TODO:
	return ExtrinsicEra{}
}

func (e ExtrinsicEra) Bytes() []byte {
	return sc.EncodedBytes(e)
}
