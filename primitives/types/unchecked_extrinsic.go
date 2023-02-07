/*
Implementation of an unchecked (pre-verification) extrinsic.
*/
package types

import (
	"bytes"
)

// A extrinsic right from the external world. This is unchecked and so can contain a signature.
type UncheckedExtrinsic struct{}

func DecodeUncheckedExtrinsic(buffer *bytes.Buffer) UncheckedExtrinsic {
	return UncheckedExtrinsic{}
}

func (uxt UncheckedExtrinsic) Bytes() []byte {
	return []byte{}
}

func (uxt UncheckedExtrinsic) Check() (xt CheckedExtrinsic, err TransactionValidityError) {
	return xt, err
}
