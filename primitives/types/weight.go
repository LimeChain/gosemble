package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// The weight of database operations that the runtime can invoke.
//
// NOTE: This is currently only measured in computational time, and will probably
// be updated all together once proof size is accounted for.
type RuntimeDbWeight struct {
	Read  sc.U64
	Write sc.U64
}

func (dbw RuntimeDbWeight) Reads(r sc.U64) Weight {
	return WeightFromParts(dbw.Read.SaturatingMul(r), 0)
}

func (dbw RuntimeDbWeight) Writes(w sc.U64) Weight {
	return WeightFromParts(dbw.Write.SaturatingMul(w), 0)
}

func (dbw RuntimeDbWeight) ReadsWrites(r, w sc.U64) Weight {
	readWeight := dbw.Read.SaturatingMul(r)
	writeWeight := dbw.Write.SaturatingMul(w)
	return WeightFromParts(readWeight.SaturatingAdd(writeWeight), 0)
}

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

func (lhs Weight) Add(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime + rhs.RefTime,
		ProofSize: lhs.ProofSize + rhs.ProofSize,
	}
}

func (lhs Weight) SaturatingAdd(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime.SaturatingAdd(rhs.RefTime),
		ProofSize: lhs.ProofSize.SaturatingAdd(rhs.ProofSize),
	}
}

// Saturating [`Weight`] subtraction. Computes `self - rhs`, saturating at the numeric bounds
// of all fields instead of overflowing.
func (lhs Weight) SaturatingSub(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime.SaturatingSub(rhs.RefTime),
		ProofSize: lhs.ProofSize.SaturatingSub(rhs.ProofSize),
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
func (lhs Weight) CheckedAdd(rhs Weight) sc.Option[Weight] {
	refTime, err := lhs.RefTime.CheckedAdd(rhs.RefTime)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	proofSize, err := lhs.ProofSize.CheckedAdd(rhs.ProofSize)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	return sc.NewOption[Weight](Weight{refTime, proofSize})
}

func (lhs Weight) Sub(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime - rhs.RefTime,
		ProofSize: lhs.ProofSize - rhs.ProofSize,
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
		RefTime:   w.RefTime.SaturatingMul(b),
		ProofSize: w.ProofSize.SaturatingMul(b),
	}
}

// Get the conservative min of `self` and `other` weight.
func (lhs Weight) Min(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime.Min(rhs.RefTime),
		ProofSize: lhs.ProofSize.Min(rhs.ProofSize),
	}
}

// Get the aggressive max of `self` and `other` weight.
func (lhs Weight) Max(rhs Weight) Weight {
	return Weight{
		RefTime:   lhs.RefTime.Max(rhs.RefTime),
		ProofSize: lhs.ProofSize.Max(rhs.ProofSize),
	}
}

// Returns true if all of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (lhs Weight) AllGt(rhs Weight) sc.Bool {
	return lhs.RefTime > rhs.RefTime && lhs.ProofSize > rhs.ProofSize
}

// / Returns true if any of `self`'s constituent weights is strictly greater than that of the
// / `other`'s, otherwise returns false.
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
