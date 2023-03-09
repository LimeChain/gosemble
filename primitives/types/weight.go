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
	RefTime sc.U64 // codec compact

	// The weight of storage space used by proof of validity.
	ProofSize sc.U64 // codec compact
}

func (w Weight) Encode(buffer *bytes.Buffer) {
	w.RefTime.Encode(buffer)
	w.ProofSize.Encode(buffer)
}

func DecodeWeight(buffer *bytes.Buffer) Weight {
	w := Weight{}
	w.RefTime = sc.DecodeU64(buffer)
	w.ProofSize = sc.DecodeU64(buffer)
	return w
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

func (w Weight) AnyGt(otherW Weight) sc.Bool {
	// TODO: check if this is correct
	return w.RefTime > otherW.RefTime // || w.ProofSize > otherW.ProofSize
}

// Construct [`Weight`] from weight parts, namely reference time and proof size weights.
func WeightFromParts(refTime sc.U64, proofSize sc.U64) Weight {
	return Weight{refTime, proofSize}
}

// Return a [`Weight`] where all fields are zero.
func WeightZero() Weight {
	return Weight{RefTime: 0, ProofSize: 0}
}
