package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type BlockLength struct {
	//  Maximal total length in bytes for each extrinsic class.
	//
	// In the worst case, the total block length is going to be:
	// `MAX(max)`
	Max PerDispatchClassU32
}

func (bl BlockLength) Encode(buffer *bytes.Buffer) error {
	return bl.Max.Encode(buffer)
}

func (bl BlockLength) Bytes() []byte {
	return sc.EncodedBytes(bl)
}

func (bl BlockLength) Docs() string {
	return "The maximum length of a block (in bytes)."
}
