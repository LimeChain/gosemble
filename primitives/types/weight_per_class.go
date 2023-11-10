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

func (cl WeightsPerClass) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		cl.BaseExtrinsic,
		cl.MaxExtrinsic,
		cl.MaxTotal,
		cl.Reserved,
	)
}

func DecodeWeightsPerClass(buffer *bytes.Buffer) (WeightsPerClass, error) {
	cl := WeightsPerClass{}
	baseExtrinsic, err := DecodeWeight(buffer)
	if err != nil {
		return WeightsPerClass{}, err
	}
	maxExtrinsic, err := sc.DecodeOptionWith(buffer, DecodeWeight)
	if err != nil {
		return WeightsPerClass{}, err
	}
	maxTotal, err := sc.DecodeOptionWith(buffer, DecodeWeight)
	if err != nil {
		return WeightsPerClass{}, err
	}
	reserved, err := sc.DecodeOptionWith(buffer, DecodeWeight)
	if err != nil {
		return WeightsPerClass{}, err
	}

	cl.BaseExtrinsic = baseExtrinsic
	cl.MaxExtrinsic = maxExtrinsic
	cl.MaxTotal = maxTotal
	cl.Reserved = reserved
	return cl, nil
}

func (cl WeightsPerClass) Bytes() []byte {
	return sc.EncodedBytes(cl)
}
