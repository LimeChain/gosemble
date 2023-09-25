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

func (w Weight) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(w.RefTime).Encode(buffer)
	sc.ToCompact(w.ProofSize).Encode(buffer)
}

func DecodeWeight(buffer *bytes.Buffer) Weight {
	refTime := sc.DecodeCompact(buffer).ToBigInt()
	proofSize := sc.DecodeCompact(buffer).ToBigInt()

	return Weight{
		RefTime:   sc.U64(refTime.Uint64()),
		ProofSize: sc.U64(proofSize.Uint64()),
	}
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
		RefTime:   w.RefTime + rhs.RefTime,     // saturating_add
		ProofSize: w.ProofSize + rhs.ProofSize, // saturating_add
	}
}

// Saturating [`Weight`] subtraction. Computes `self - rhs`, saturating at the numeric bounds
// of all fields instead of overflowing.
func (w Weight) SaturatingSub(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime - rhs.RefTime,     // saturating_sub
		ProofSize: w.ProofSize - rhs.ProofSize, // saturating_sub
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

// Checked [`Weight`] addition. Computes `self + rhs`, returning `None` if overflow occurred.
func (w Weight) CheckedAdd(rhs Weight) sc.Option[Weight] {
	refTime := w.RefTime + rhs.RefTime       // checked_add
	proofSize := w.ProofSize + rhs.ProofSize // checked_адд

	return sc.NewOption[Weight](Weight{refTime, proofSize})
}

func (w Weight) Sub(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime - rhs.RefTime,
		ProofSize: w.ProofSize - rhs.ProofSize,
	}
}

func (w Weight) Mul(b sc.U64) Weight {
	return Weight{
		RefTime:   w.RefTime * b,
		ProofSize: w.ProofSize * b,
	}
}

func (w Weight) SaturatingMul(b sc.U64) Weight {
	return Weight{
		RefTime:   w.RefTime * b,   // saturating_mul
		ProofSize: w.ProofSize * b, // saturating_mul
	}
}

// Min Get the conservative min of `self` and `other` weight.
func (w Weight) Min(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.U64(minUint64(uint64(w.RefTime), uint64(rhs.RefTime))),     // min
		ProofSize: sc.U64(minUint64(uint64(w.ProofSize), uint64(rhs.ProofSize))), // min
	}
}

// Max Get the aggressive max of `self` and `other` weight.
func (w Weight) Max(rhs Weight) Weight {
	return Weight{
		RefTime:   sc.U64(maxUint64(uint64(w.RefTime), uint64(rhs.RefTime))),     // max
		ProofSize: sc.U64(maxUint64(uint64(w.ProofSize), uint64(rhs.ProofSize))), // max
	}
}

// AllGt Returns true if all of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AllGt(rhs Weight) sc.Bool {
	return w.RefTime > rhs.RefTime && w.ProofSize > rhs.ProofSize
}

// AnyGt Returns true if any of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AnyGt(otherW Weight) sc.Bool {
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
