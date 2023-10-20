package types

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectBytesValidTransaction, _ = hex.DecodeString("0100000000000000040c020304080c0506070c08090a0b0000000000000001")
)

var (
	targetValidTransaction = ValidTransaction{
		Priority: 1,
		Requires: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.Sequence[sc.U8]{
				2, 3, 4,
			},
		},
		Provides: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.Sequence[sc.U8]{
				5, 6, 7,
			},
			sc.Sequence[sc.U8]{
				8, 9, 10,
			},
		},
		Longevity: 11,
		Propagate: true,
	}
)

func Test_ValidTransaction_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetValidTransaction.Encode(buffer)

	assert.Equal(t, expectBytesValidTransaction, buffer.Bytes())
}

func Test_DecodeValidTransaction(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesValidTransaction)

	result := DecodeValidTransaction(buffer)

	assert.Equal(t, targetValidTransaction, result)
}

func Test_ValidTransaction_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesValidTransaction, targetValidTransaction.Bytes())
}

func Test_DefaultValidTransaction(t *testing.T) {
	expect := ValidTransaction{
		Priority:  0,
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}

	assert.Equal(t, expect, DefaultValidTransaction())
}

func Test_ValidTransaction_CombineWith(t *testing.T) {
	other := ValidTransaction{
		Priority: math.MaxUint64 - 1,
		Requires: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.Sequence[sc.U8]{
				100, 111, 112,
			},
		},
		Provides: sc.Sequence[sc.Sequence[sc.U8]]{
			sc.Sequence[sc.U8]{
				113, 114, 115,
			},
		},
		Longevity: 1000,
		Propagate: false,
	}

	expect := ValidTransaction{
		Priority:  math.MaxUint64,
		Requires:  append(targetValidTransaction.Requires, other.Requires...),
		Provides:  append(targetValidTransaction.Provides, other.Provides...),
		Longevity: targetValidTransaction.Longevity,
		Propagate: false,
	}

	assert.Equal(t, expect, targetValidTransaction.CombineWith(other))
}
