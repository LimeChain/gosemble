package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type AccountIdExtra struct {
	Address32
	SignedExtra
}

func (ae AccountIdExtra) Encode(buffer *bytes.Buffer) {
	ae.Address32.Encode(buffer)
	ae.SignedExtra.Encode(buffer)
}

func (ae AccountIdExtra) Bytes() []byte {
	return sc.EncodedBytes(ae)
}
