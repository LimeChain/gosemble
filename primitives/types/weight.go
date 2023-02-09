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

// Construct [`Weight`] with reference time weight and 0 storage size weight.
func WeightFromRefTime(refTime sc.U64) Weight {
	return Weight{
		RefTime:   refTime,
		ProofSize: 0,
	}
}
