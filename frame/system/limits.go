package system

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"

	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Block resource limits configuration structures.
//
// FRAME defines two resources that are limited within a block:
// - Weight (execution cost/time)
// - Length (block size)
//
// `frame_system` tracks consumption of each of these resources separately for each
// `DispatchClass`. This module contains configuration object for both resources,
// which should be passed to `frame_system` configuration when runtime is being set up.

// TODO: needs to be benchmarked

const FiveMbPerBlockPerExtrinsic sc.U32 = 1024 // TODO: 5 * 1024 * 1024

const WeightRefTimePerSecond sc.U64 = 1_000_000_000_000
const WeightRefTimePerNanos sc.U64 = 1_000

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
var BlockExecutionWeight types.Weight = types.WeightFromRefTime(WeightRefTimePerNanos.SaturatingMul(381_015))

// A ratio of `Normal` dispatch class within block, used as default value for
// `BlockWeight` and `BlockLength`. The `Default` impls are provided mostly for convenience
// to use in tests.
var DefaultNormalRatio = Perbill{Percentage: 75}

type Perbill struct {
	Percentage sc.U32
}

func (p Perbill) Encode(buffer *bytes.Buffer) {
	p.Percentage.Encode(buffer)
}

func DecodePerbill(buffer *bytes.Buffer) Perbill {
	p := Perbill{}
	p.Percentage = sc.DecodeU32(buffer)
	return p
}

func (p Perbill) Bytes() []byte {
	return sc.EncodedBytes(p)
}

func (p Perbill) Mul(v sc.U32) sc.U32 {
	return (v / 100) * p.Percentage
}

type BlockLength struct {
	//  Maximal total length in bytes for each extrinsic class.
	//
	// In the worst case, the total block length is going to be:
	// `MAX(max)`
	Max types.PerDispatchClass[sc.U32]
}

func DefaultBlockLength() BlockLength {
	return MaxWithNormalRatio(FiveMbPerBlockPerExtrinsic, DefaultNormalRatio)
}

// Create new `BlockLength` with `max` for `Operational` & `Mandatory`
// and `normal * max` for `Normal`.
func MaxWithNormalRatio(max sc.U32, normal Perbill) BlockLength {
	return BlockLength{
		Max: types.PerDispatchClass[sc.U32]{
			Normal:      normal.Mul(max),
			Operational: max,
			Mandatory:   max,
		},
	}
}

// / `DispatchClass`-specific weight configuration.
type WeightsPerClass struct {
	// Base weight of single extrinsic of given class.
	BaseExtrinsic types.Weight

	// Maximal weight of single extrinsic. Should NOT include `base_extrinsic` cost.
	//
	// `None` indicates that this class of extrinsics doesn't have a limit.
	MaxExtrinsic sc.Option[types.Weight]

	// Block maximal total weight for all extrinsics of given class.
	//
	// `None` indicates that weight sum of this class of extrinsics is not
	// restricted. Use this value carefully, since it might produce heavily oversized
	// blocks.
	//
	// In the worst case, the total weight consumed by the class is going to be:
	// `MAX(max_total) + MAX(reserved)`.
	MaxTotal sc.Option[types.Weight]

	// Block reserved allowance for all extrinsics of a particular class.
	//
	// Setting to `None` indicates that extrinsics of that class are allowed
	// to go over total block weight (but at most `max_total` for that class).
	// Setting to `Some(x)` guarantees that at least `x` weight of particular class
	// is processed in every block.
	Reserved sc.Option[types.Weight]
}

func (cl WeightsPerClass) Encode(buffer *bytes.Buffer) {
	cl.BaseExtrinsic.Encode(buffer)
	cl.MaxExtrinsic.Encode(buffer)
	cl.MaxTotal.Encode(buffer)
	cl.Reserved.Encode(buffer)
}

func DecodeWeightsPerClass(buffer *bytes.Buffer) WeightsPerClass {
	cl := WeightsPerClass{}
	cl.BaseExtrinsic = types.DecodeWeight(buffer)
	cl.MaxExtrinsic = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	cl.MaxTotal = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	cl.Reserved = sc.DecodeOptionWith(buffer, types.DecodeWeight)
	return cl
}

func (cl WeightsPerClass) Bytes() []byte {
	return sc.EncodedBytes(cl)
}

type BlockWeights struct {
	// Base weight of block execution.
	BaseBlock types.Weight

	// Maximal total weight consumed by all kinds of extrinsics (without `reserved` space).
	MaxBlock types.Weight

	// Weight limits for extrinsics of given dispatch class.
	PerClass types.PerDispatchClass[WeightsPerClass]
}

func DefaultBlockWeights() BlockWeights {
	WithSensibleDefaults(
		types.WeightFromParts(WeightRefTimePerSecond, math.MaxUint64),
		DefaultNormalRatio,
	)
	return BlockWeights{}
}

// Get per-class weight settings.
func (bw BlockWeights) Get(class types.DispatchClass) *WeightsPerClass {
	if class.Is(types.DispatchClassNormal) {
		return &bw.PerClass.Normal
	} else if class.Is(types.DispatchClassOperational) {
		return &bw.PerClass.Operational
	} else if class.Is(types.DispatchClassMandatory) {
		return &bw.PerClass.Mandatory
	} else {
		log.Critical("Invalid dispatch class")
	}

	panic("unreachable")
}

// Create a sensible default weights system given only expected maximal block weight and the
// ratio that `Normal` extrinsics should occupy.
//
// Assumptions:
//   - Average block initialization is assumed to be `10%`.
//   - `Operational` transactions have reserved allowance (`1.0 - normal_ratio`)
func WithSensibleDefaults(expectedBlockWeight types.Weight, normalRatio Perbill) BlockWeights {
	// Start constructing new `BlockWeights` object.
	//
	// By default all kinds except of `Mandatory` extrinsics are disallowed.
	WeightsForNormalAndOperational := WeightsPerClass{
		BaseExtrinsic: BlockExecutionWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](nil),
		Reserved:      sc.NewOption[types.Weight](nil),
	}

	WeightsForMandatory := WeightsPerClass{
		BaseExtrinsic: BlockExecutionWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](types.WeightZero()),
		Reserved:      sc.NewOption[types.Weight](types.WeightZero()),
	}

	weightsPerClass := types.PerDispatchClass[WeightsPerClass]{
		Mandatory:   WeightsForMandatory,
		Normal:      WeightsForNormalAndOperational,
		Operational: WeightsForNormalAndOperational,
	}

	builder := BlockWeightsBuilder{
		Weights: BlockWeights{
			BaseBlock: BlockExecutionWeight,
			MaxBlock:  types.WeightZero(),
			PerClass:  weightsPerClass,
		},
		InitCost: sc.NewOption[Perbill](nil),
	}

	normalWeight := expectedBlockWeight // TODO: normalRatio *

	// Set parameters for particular class.
	//
	// Note: `None` values of `max_extrinsic` will be overwritten in `build` in case
	// `avg_block_initialization` rate is set to a non-zero value.
	builder.Weights.PerClass.Normal.MaxTotal = sc.NewOption[types.Weight](normalWeight)
	builder.Weights.PerClass.Operational.MaxTotal = sc.NewOption[types.Weight](expectedBlockWeight)
	builder.Weights.PerClass.Operational.Reserved = sc.NewOption[types.Weight](expectedBlockWeight.Sub(normalWeight))

	// Average block initialization weight cost.
	//
	// This value is used to derive maximal allowed extrinsic weight for each
	// class, based on the allowance.
	//
	// This is to make sure that extrinsics don't stay forever in the pool,
	// because they could seamingly fit the block (since they are below `max_block`),
	// but the cost of calling `on_initialize` always prevents them from being included.
	builder.InitCost = sc.NewOption[Perbill](Perbill{10})

	// .build()

	return BlockWeights{}
}

// Construct the `BlockWeights` object.
// func (bwb BlockWeightsBuilder) Build()  (ok BlockWeights, err error) {
// 	// compute max extrinsic size
// 	let Self { mut weights, init_cost } = self;

// 	// compute max block size.
// 	for class in DispatchClass::all() {
// 		weights.max_block = match weights.per_class.get(*class).max_total {
// 			Some(max) => max.max(weights.max_block),
// 			_ => weights.max_block,
// 		};
// 	}
// 	// compute max size of single extrinsic
// 	if let Some(init_weight) = init_cost.map(|rate| rate * weights.max_block) {
// 		for class in DispatchClass::all() {
// 			let per_class = weights.per_class.get_mut(*class);
// 			if per_class.max_extrinsic.is_none() && init_cost.is_some() {
// 				per_class.max_extrinsic = per_class
// 					.max_total
// 					.map(|x| x.saturating_sub(init_weight))
// 					.map(|x| x.saturating_sub(per_class.base_extrinsic));
// 			}
// 		}
// 	}

// 	// Validate the result
// 	weights.validate()
// }

// An opinionated builder for `Weights` object.
type BlockWeightsBuilder struct {
	Weights  BlockWeights
	InitCost sc.Option[Perbill]
}
