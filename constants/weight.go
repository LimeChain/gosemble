package constants

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TODO: needs to be benchmarked
const FiveMbPerBlockPerExtrinsic sc.U32 = 1024 // TODO: 5 * 1024 * 1024

const WeightRefTimePerSecond sc.U64 = 1_000_000_000_000
const WeightRefTimePerNanos sc.U64 = 1_000_000_000

// TODO: update according to the DB used
// this will be the weight used throughout the runtime.
var DbWeight types.RuntimeDbWeight = types.RuntimeDbWeight{
	Read:  25_000 * WeightRefTimePerNanos,
	Write: 100_000 * WeightRefTimePerNanos,
}

// Time to execute a NO-OP extrinsic, for example `System::remark`.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
//
// Stats nanoseconds:
//
//	Min, Max: 99_481, 103_304
//	Average:  99_840
//	Median:   99_795
//	Std-Dev:  376.17
//
// Percentiles nanoseconds:
//
//	99th: 100_078
//	95th: 100_051
//	75th: 99_916
var ExtrinsicBaseWeight types.Weight = types.WeightFromParts(WeightRefTimePerNanos.SaturatingMul(99_840), 0)

// Time to execute an empty block.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
//
// Stats nanoseconds:
//
//	Min, Max: 377_722, 414_752
//	Average:  381_015
//	Median:   379_751
//	Std-Dev:  5462.64
//
// Percentiles nanoseconds:
//
//	99th: 413_074
//	95th: 384_876
//	75th: 380_642
var BlockExecutionWeight types.Weight = types.WeightFromParts(WeightRefTimePerNanos.SaturatingMul(381_015), 0)
