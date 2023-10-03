package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// WeightsPerClass `DispatchClass`-specific weight configuration.
type WeightsPerClass struct {
	// Base weight of single extrinsic of given class.
	BaseExtrinsic Weight

	// Maximal weight of single extrinsic. Should NOT include `base_extrinsic` cost.
	//
	// `None` indicates that this class of extrinsics doesn't have a limit.
	MaxExtrinsic sc.Option[Weight]

	// Block maximal total weight for all extrinsics of given class.
	//
	// `None` indicates that weight sum of this class of extrinsics is not
	// restricted. Use this value carefully, since it might produce heavily oversized
	// blocks.
	//
	// In the worst case, the total weight consumed by the class is going to be:
	// `MAX(max_total) + MAX(reserved)`.
	MaxTotal sc.Option[Weight]

	// Block reserved allowance for all extrinsics of a particular class.
	//
	// Setting to `None` indicates that extrinsics of that class are allowed
	// to go over total block weight (but at most `max_total` for that class).
	// Setting to `Some(x)` guarantees that at least `x` weight of particular class
	// is processed in every block.
	Reserved sc.Option[Weight]
}

func (cl WeightsPerClass) Encode(buffer *bytes.Buffer) {
	cl.BaseExtrinsic.Encode(buffer)
	cl.MaxExtrinsic.Encode(buffer)
	cl.MaxTotal.Encode(buffer)
	cl.Reserved.Encode(buffer)
}

func DecodeWeightsPerClass(buffer *bytes.Buffer) WeightsPerClass {
	cl := WeightsPerClass{}
	cl.BaseExtrinsic = DecodeWeight(buffer)
	cl.MaxExtrinsic = sc.DecodeOptionWith(buffer, DecodeWeight)
	cl.MaxTotal = sc.DecodeOptionWith(buffer, DecodeWeight)
	cl.Reserved = sc.DecodeOptionWith(buffer, DecodeWeight)
	return cl
}

func (cl WeightsPerClass) Bytes() []byte {
	return sc.EncodedBytes(cl)
}
