package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

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

func (w Weight) Sub(otherW Weight) Weight {
	return Weight{
		RefTime:   w.RefTime.Sub(otherW.RefTime),
		ProofSize: w.ProofSize.Sub(otherW.ProofSize),
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

// Construct [`Weight`] with reference time weight and 0 storage size weight.
func WeightFromRefTime(refTime sc.U64) Weight {
	return Weight{
		RefTime:   refTime,
		ProofSize: 0,
	}
}

// Return a [`Weight`] where all fields are zero.
func WeightZero() Weight {
	return Weight{RefTime: 0, ProofSize: 0}
}
