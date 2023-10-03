package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
)

// MaxWithNormalRatio Create new `BlockLength` with `max` for `Operational` & `Mandatory`
// and `normal * max` for `Normal`.
func MaxWithNormalRatio(max sc.U32, normal types.Perbill) types.BlockLength {
	return types.BlockLength{
		Max: types.PerDispatchClass[sc.U32]{
			Normal:      normal.Mul(max).(sc.U32),
			Operational: max,
			Mandatory:   max,
		},
	}
}

// WithSensibleDefaults Create a sensible default weights system given only expected maximal block weight and the
// ratio that `Normal` extrinsics should occupy.
//
// Assumptions:
//   - Average block initialization is assumed to be `10%`.
//   - `Operational` transactions have reserved allowance (`1.0 - normal_ratio`)
func WithSensibleDefaults(expectedBlockWeight types.Weight, normalRatio types.Perbill) types.BlockWeights {
	normalWeight := normalRatio.Mul(expectedBlockWeight)
	return NewBlockWeightsBuilder().
		ForClass([]types.DispatchClass{types.NewDispatchClassNormal()}, func(weights *types.WeightsPerClass) {
			weights.MaxTotal = sc.NewOption[types.Weight](normalWeight)
		}).
		ForClass([]types.DispatchClass{types.NewDispatchClassOperational()}, func(weights *types.WeightsPerClass) {
			weights.MaxTotal = sc.NewOption[types.Weight](expectedBlockWeight)
			weights.Reserved = sc.NewOption[types.Weight](expectedBlockWeight.Sub(normalWeight.(types.Weight)))
		}).
		AvgBlockInitialization(constants.AverageOnInitializeRatio).
		Build()
	// TODO: builder.Expect("Sensible defaults are tested to be valid")
}

// An opinionated builder for `Weights` object.
type BlockWeightsBuilder struct {
	Weights  types.BlockWeights
	InitCost sc.Option[types.Perbill]
}

// Start constructing new `BlockWeights` object.
//
// By default all kinds except of `Mandatory` extrinsics are disallowed.
func NewBlockWeightsBuilder() *BlockWeightsBuilder {
	// Start constructing new `BlockWeights` object.
	//
	// By default all kinds except of `Mandatory` extrinsics are disallowed.
	WeightsForNormalAndOperational := types.WeightsPerClass{
		BaseExtrinsic: constants.ExtrinsicBaseWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](types.WeightZero()),
		Reserved:      sc.NewOption[types.Weight](types.WeightZero()),
	}

	WeightsForMandatory := types.WeightsPerClass{
		BaseExtrinsic: constants.ExtrinsicBaseWeight,
		MaxExtrinsic:  sc.NewOption[types.Weight](nil),
		MaxTotal:      sc.NewOption[types.Weight](nil),
		Reserved:      sc.NewOption[types.Weight](nil),
	}

	weightsPerClass := types.PerDispatchClass[types.WeightsPerClass]{
		Mandatory:   WeightsForMandatory,
		Normal:      WeightsForNormalAndOperational,
		Operational: WeightsForNormalAndOperational,
	}

	return &BlockWeightsBuilder{
		Weights: types.BlockWeights{
			BaseBlock: constants.ExtrinsicBaseWeight,
			MaxBlock:  types.WeightZero(),
			PerClass:  weightsPerClass,
		},
		InitCost: sc.NewOption[types.Perbill](nil),
	}
}

// Set base block weight.
func (b *BlockWeightsBuilder) BaseBlock(baseBlock types.Weight) *BlockWeightsBuilder {
	b.Weights.BaseBlock = baseBlock
	return b
}

// ForClass Set parameters for particular class.
//
// Note: `None` values of `max_extrinsic` will be overwritten in `build` in case
// `avg_block_initialization` rate is set to a non-zero value.
func (b *BlockWeightsBuilder) ForClass(classes []types.DispatchClass, action func(_ *types.WeightsPerClass)) *BlockWeightsBuilder {
	for _, cl := range classes {
		action(b.Weights.PerClass.Get(cl))
	}
	return b
}

// AvgBlockInitialization Average block initial ization weight cost.
//
// This value is used to derive maximal allowed extrinsic weight for each
// class, based on the allowance.
//
// This is to make sure that extrinsics don't stay forever in the pool,
// because they could seamingly fit the block (since they are below `max_block`),
// but the cost of calling `on_initialize` always prevents them from being included.
func (b *BlockWeightsBuilder) AvgBlockInitialization(initCost types.Perbill) *BlockWeightsBuilder {
	b.InitCost = sc.NewOption[types.Perbill](initCost)
	return b
}

// Construct the `BlockWeights` object.
func (b *BlockWeightsBuilder) Build() types.BlockWeights {
	// compute max extrinsic size
	weights, initCost := b.Weights, b.InitCost

	// compute max block size.
	for _, class := range types.DispatchClassAll() {
		if (*weights.PerClass.Get(class)).MaxTotal.HasValue {
			max := (*weights.PerClass.Get(class)).MaxTotal.Value
			weights.MaxBlock = max.Max(weights.MaxBlock)
		}
	}

	// compute max size of single extrinsic
	var initWeight sc.Option[types.Weight]
	if initCost.HasValue {
		initWeight = sc.NewOption[types.Weight](initCost.Value.Mul(weights.MaxBlock))
	} else {
		initWeight = sc.NewOption[types.Weight](nil)
	}

	if initWeight.HasValue {
		for _, class := range types.DispatchClassAll() {
			perClass := *(weights.PerClass.Get(class))
			if !perClass.MaxExtrinsic.HasValue && initCost.HasValue {
				if perClass.MaxTotal.HasValue {
					perClass.MaxExtrinsic = sc.NewOption[types.Weight](perClass.MaxTotal.Value.SaturatingSub(initWeight.Value).SaturatingSub(perClass.BaseExtrinsic))
				} else {
					perClass.MaxExtrinsic = sc.NewOption[types.Weight](nil)
				}
			}
		}
	}

	// Validate the result
	// TODO: weights.Validate()
	return weights
}
