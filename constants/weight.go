package constants

import (
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TODO: needs to be benchmarked

const FiveMbPerBlockPerExtrinsic sc.U32 = 5 * 1024 * 1024
const WeightRefTimePerSecond sc.U64 = 1_000_000_000_000
const WeightRefTimePerNanos sc.U64 = 1_000

// We assume that ~10% of the block weight is consumed by `on_initialize` handlers.
// This is used to limit the maximal weight of a single extrinsic.
var AverageOnInitializeRatio types.Perbill = types.Perbill{Percentage: 10}

// We allow `Normal` extrinsics to fill up the block up to 75%, the rest can be used
// by  Operational  extrinsics.
var NormalDispatchRatio types.Perbill = types.Perbill{Percentage: 75}

// Block resource limits configuration structures.
//
// FRAME defines two resources that are limited within a block:
// - Weight (execution cost/time)
// - Length (block size)
//
// `frame_system` tracks consumption of each of these resources separately for each
// `DispatchClass`. This module contains configuration object for both resources,
// which should be passed to `frame_system` configuration when runtime is being set up.

// A ratio of `Normal` dispatch class within block, used as default value for
// `BlockWeight` and `BlockLength`. The `Default` impls are provided mostly for convenience
// to use in tests.

// ExtrinsicBaseWeight is the time to execute a NO-OP extrinsic, for example `System::remark`.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
var ExtrinsicBaseWeight = baseExtrinsicWeight(WeightRefTimePerNanos)

// Time to execute an empty block.
// Calculated by multiplying the *Average* with `1.0` and adding `0`.
var BlockExecutionWeight = blockExecutionWeight(WeightRefTimePerNanos)

// MaximumBlockWeight is the maximum weight 2 seconds of compute with a 6-second average block time, with maximum proof size.
var MaximumBlockWeight = types.WeightFromParts(
	sc.SaturatingMulU64(WeightRefTimePerSecond, 2),
	math.MaxUint64,
)

// RocksDbWeight for RocksDB, used throughout the runtime.
var RocksDbWeight = types.RuntimeDbWeight{
	Read:  sc.U64(25_000) * WeightRefTimePerNanos,
	Write: sc.U64(100_000) * WeightRefTimePerNanos,
}
