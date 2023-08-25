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
	refTime := sc.DecodeCompact(buffer)
	proofSize := sc.DecodeCompact(buffer)

	return Weight{
		RefTime:   sc.To[sc.U64](sc.U128(refTime)),
		ProofSize: sc.To[sc.U64](sc.U128(proofSize)),
	}
}

func (w Weight) Bytes() []byte {
	return sc.EncodedBytes(w)
}

func (w Weight) Add(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.Add(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.Add(rhs.ProofSize).(sc.U64),
	}
}

func (w Weight) SaturatingAdd(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.SaturatingAdd(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.SaturatingAdd(rhs.ProofSize).(sc.U64),
	}
}

// Saturating [`Weight`] subtraction. Computes `self - rhs`, saturating at the numeric bounds
// of all fields instead of overflowing.
func (w Weight) SaturatingSub(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.SaturatingSub(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.SaturatingSub(rhs.ProofSize).(sc.U64),
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
	refTime, err := w.RefTime.CheckedAdd(rhs.RefTime)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	proofSize, err := w.ProofSize.CheckedAdd(rhs.ProofSize)
	if err != nil {
		return sc.NewOption[Weight](nil)
	}

	return sc.NewOption[Weight](Weight{refTime.(sc.U64), proofSize.(sc.U64)})
}

func (w Weight) Sub(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.Sub(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.Sub(rhs.ProofSize).(sc.U64),
	}
}

func (w Weight) Mul(b sc.U64) Weight {
	return Weight{
		RefTime:   w.RefTime.Mul(b).(sc.U64),
		ProofSize: w.ProofSize.Mul(b).(sc.U64),
	}
}

func (w Weight) SaturatingMul(b sc.U64) Weight {
	return Weight{
		RefTime:   w.RefTime.SaturatingMul(b).(sc.U64),
		ProofSize: w.ProofSize.SaturatingMul(b).(sc.U64),
	}
}

// Min Get the conservative min of `self` and `other` weight.
func (w Weight) Min(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.Min(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.Min(rhs.ProofSize).(sc.U64),
	}
}

// Max Get the aggressive max of `self` and `other` weight.
func (w Weight) Max(rhs Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.Max(rhs.RefTime).(sc.U64),
		ProofSize: w.ProofSize.Max(rhs.ProofSize).(sc.U64),
	}
}

// AllGt Returns true if all of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AllGt(rhs Weight) sc.Bool {
	return sc.Bool(w.RefTime.Gt(rhs.RefTime) && w.ProofSize.Gt(rhs.ProofSize))
}

// AnyGt Returns true if any of `self`'s constituent weights is strictly greater than that of the
// `other`'s, otherwise returns false.
func (w Weight) AnyGt(otherW Weight) sc.Bool {
	return sc.Bool(w.RefTime.Gt(otherW.RefTime) || w.ProofSize.Gt(otherW.ProofSize))
}

// Construct [`Weight`] from weight parts, namely reference time and proof size weights.
func WeightFromParts(refTime sc.U64, proofSize sc.U64) Weight {
	return Weight{refTime, proofSize}
}

// Return a [`Weight`] where all fields are zero.
func WeightZero() Weight {
	return Weight{RefTime: 0, ProofSize: 0}
}
