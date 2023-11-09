package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Weight struct {
	// The weight of computational time used based on some reference hardware.
	RefTime sc.U64
	// The weight of storage space used by proof of validity.
	ProofSize sc.U64
}

func (w Weight) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		sc.ToCompact(w.RefTime),
		sc.ToCompact(w.ProofSize),
	)
}

func DecodeWeight(buffer *bytes.Buffer) (Weight, error) {
	refTime, err := sc.DecodeCompact(buffer)
	if err != nil {
		return Weight{}, err
	}
	proofSize, err := sc.DecodeCompact(buffer)
	if err != nil {
		return Weight{}, err
	}

	return Weight{
		RefTime:   sc.U64(refTime.ToBigInt().Uint64()),
		ProofSize: sc.U64(proofSize.ToBigInt().Uint64()),
	}, nil
}

func (w Weight) Bytes() []byte {
	return sc.EncodedBytes(w)
}

func (w Weight) Add(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime + rhs.RefTime,
		ProofSize: w.ProofSize + rhs.ProofSize,
	}
}

func (w Weight) SaturatingAdd(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.SaturatingAddU64(w.RefTime, rhs.RefTime),
		ProofSize: sc.SaturatingAddU64(w.ProofSize, rhs.ProofSize),
	}
}

// Checked [`Weight`] addition. Computes `self + rhs`, returning `None` if overflow occurred.
func (w Weight) CheckedAdd(rhs Weight) sc.Option[Weight] {
	refTime, err := sc.CheckedAddU64(w.RefTime, rhs.RefTime)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	proofSize, err := sc.CheckedAddU64(w.ProofSize, rhs.ProofSize)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	return sc.NewOption[Weight](Weight{refTime, proofSize})
}

func (w Weight) Sub(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime - rhs.RefTime,
		ProofSize: w.ProofSize - rhs.ProofSize,
	}
}

// Saturating [`Weight`] subtraction. Computes `self - rhs`, saturating at the numeric bounds
// of all fields instead of overflowing.
func (w Weight) SaturatingSub(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.SaturatingSubU64(w.RefTime, rhs.RefTime),
		ProofSize: sc.SaturatingSubU64(w.ProofSize, rhs.ProofSize),
	}
}

// Increment [`Weight`] by `amount` via saturating addition.
func (w *Weight) SaturatingAccrue(amount Weight) {
	*w = w.SaturatingAdd(amount)
}

// Reduce [`Weight`] by `amount` via saturating subtraction.
func (w *Weight) SaturatingReduce(amount Weight) {
	*w = w.SaturatingSub(amount)
}

func (w Weight) Mul(b sc.U64) Weight {
	return Weight{
		RefTime:   w.RefTime * b,
		ProofSize: w.ProofSize * b,
	}
}

func (w Weight) SaturatingMul(b sc.U64) Weight {
	return Weight{
		RefTime:   sc.SaturatingMulU64(w.RefTime, b),
		ProofSize: sc.SaturatingMulU64(w.ProofSize, b),
	}
}

// Min Get the conservative min of `self` and `other` weight.
func (w Weight) Min(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.Min64(w.RefTime, rhs.RefTime),
		ProofSize: sc.Min64(w.ProofSize, rhs.ProofSize),
	}
}

// Max Get the aggressive max of `self` and `other` weight.
func (w Weight) Max(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.Max64(w.RefTime, rhs.RefTime),
		ProofSize: sc.Max64(w.ProofSize, rhs.ProofSize),
	}
}

// AllGt Returns true if all of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AllGt(rhs Weight) bool {
	return w.RefTime > rhs.RefTime && w.ProofSize > rhs.ProofSize
}

// AnyGt Returns true if any of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AnyGt(otherW Weight) bool {
	return w.RefTime > otherW.RefTime || w.ProofSize > otherW.ProofSize
}

// Construct [`Weight`] from weight parts, namely reference time and proof size weights.
func WeightFromParts(refTime sc.U64, proofSize sc.U64) Weight {
	return Weight{refTime, proofSize}
}

// Return a [`Weight`] where all fields are zero.
func WeightZero() Weight {
	return Weight{RefTime: 0, ProofSize: 0}
}
