package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckWeight struct {
	systemModule system.Module
}

func NewCheckWeight(systemModule system.Module) CheckWeight {
	return CheckWeight{
		systemModule: systemModule,
	}
}

func (cw CheckWeight) Encode(*bytes.Buffer) {}

func (cw CheckWeight) Decode(*bytes.Buffer) {}

func (cw CheckWeight) Bytes() []byte {
	return sc.EncodedBytes(cw)
}

func (cw CheckWeight) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return primitives.AdditionalSigned{}, nil
}

func (cw CheckWeight) Validate(_who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return cw.doValidate(info, length)
}

func (cw CheckWeight) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return cw.doValidate(info, length)
}

func (cw CheckWeight) PreDispatch(_who *primitives.Address32, _call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := cw.doPreDispatch(info, length)
	return primitives.Pre{}, err
}

func (cw CheckWeight) PreDispatchUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cw.doPreDispatch(info, length)
	return err
}

func (cw CheckWeight) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	unspent := postInfo.CalcUnspent(info)
	if unspent.AnyGt(primitives.WeightZero()) {
		currentWeight := cw.systemModule.Storage.BlockWeight.Get()
		currentWeight.Reduce(unspent, info.Class)
		cw.systemModule.Storage.BlockWeight.Put(currentWeight)
	}
	return nil
}

// Do the validate checks. This can be applied to both signed and unsigned.
//
// It only checks that the block weight and length limit will not exceed.
func (cw CheckWeight) doValidate(info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	// ignore the next length. If they return `Ok`, then it is below the limit.
	_, err := cw.checkBlockLength(info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	// during validation we skip block limit check. Since the `validate_transaction`
	// call runs on an empty block anyway, by this we prevent `on_initialize` weight
	// consumption from causing false negatives.
	_, err = cw.checkExtrinsicWeight(info)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	return primitives.DefaultValidTransaction(), err
}

func (cw CheckWeight) doPreDispatch(info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	nextLength, err := cw.checkBlockLength(info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	nextWeight, err := cw.checkBlockWeight(info)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	_, err = cw.checkExtrinsicWeight(info)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	cw.systemModule.Storage.AllExtrinsicsLen.Put(nextLength)
	cw.systemModule.Storage.BlockWeight.Put(nextWeight)

	return primitives.ValidTransaction{}, err
}

// Checks if the current extrinsic can fit into the block with respect to block length limits.
//
// Upon successes, it returns the new block length as a `Result`.
func (cw CheckWeight) checkBlockLength(info *primitives.DispatchInfo, length sc.Compact) (sc.U32, primitives.TransactionValidityError) {
	lengthLimit := cw.systemModule.Constants.BlockLength
	currentLen := cw.systemModule.Storage.AllExtrinsicsLen.Get()
	addedLen := sc.U32(sc.U128(length).ToBigInt().Uint64())

	nextLen := currentLen + addedLen // saturating_add

	var maxLimit sc.U32
	if info.Class.Is(primitives.DispatchClassNormal) {
		maxLimit = lengthLimit.Max.Normal
	} else if info.Class.Is(primitives.DispatchClassOperational) {
		maxLimit = lengthLimit.Max.Operational
	} else if info.Class.Is(primitives.DispatchClassMandatory) {
		maxLimit = lengthLimit.Max.Mandatory
	} else {
		log.Critical("invalid DispatchClass type in CheckBlockLength()")
	}

	if nextLen > maxLimit {
		return sc.U32(0), primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
	}

	return sc.U32(sc.ToCompact(nextLen).ToBigInt().Uint64()), nil
}

// Checks if the current extrinsic can fit into the block with respect to block weight limits.
//
// Upon successes, it returns the new block weight as a `Result`.
func (cw CheckWeight) checkBlockWeight(info *primitives.DispatchInfo) (primitives.ConsumedWeight, primitives.TransactionValidityError) {
	maximumWeight := cw.systemModule.Constants.BlockWeights
	allWeight := cw.systemModule.Storage.BlockWeight.Get()
	return cw.calculateConsumedWeight(maximumWeight, allWeight, info)
}

// Checks if the current extrinsic does not exceed the maximum weight a single extrinsic
// with given `DispatchClass` can have.
func (cw CheckWeight) checkExtrinsicWeight(info *primitives.DispatchInfo) (sc.Empty, primitives.TransactionValidityError) {
	max := cw.systemModule.Constants.BlockWeights.Get(info.Class).MaxExtrinsic

	if max.HasValue {
		if info.Weight.AnyGt(max.Value) {
			return sc.Empty{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	}

	return sc.Empty{}, nil
}

func (cw CheckWeight) calculateConsumedWeight(maximumWeight system.BlockWeights, allConsumedWeight primitives.ConsumedWeight, info *primitives.DispatchInfo) (primitives.ConsumedWeight, primitives.TransactionValidityError) {
	limitPerClass := maximumWeight.Get(info.Class)
	extrinsicWeight := info.Weight.SaturatingAdd(limitPerClass.BaseExtrinsic)

	// add the weight. If class is unlimited, use saturating add instead of checked one.
	if !limitPerClass.MaxTotal.HasValue && !limitPerClass.Reserved.HasValue {
		allConsumedWeight.Accrue(extrinsicWeight, info.Class)
	} else {
		_, e := allConsumedWeight.CheckedAccrue(extrinsicWeight, info.Class)
		if e != nil {
			return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	}

	consumedPerClass := allConsumedWeight.Get(info.Class)

	// Check if we don't exceed per-class allowance
	if limitPerClass.MaxTotal.HasValue {
		max := limitPerClass.MaxTotal.Value
		if consumedPerClass.AnyGt(max) {
			return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	} else {
		// There is no `max_total` limit (`None`),
		// or we are below the limit.
	}

	// In cases total block weight is exceeded, we need to fall back
	// to `reserved` pool if there is any.
	if allConsumedWeight.Total().AnyGt(maximumWeight.MaxBlock) {
		if limitPerClass.Reserved.HasValue {
			reserved := limitPerClass.Reserved.Value
			if consumedPerClass.AnyGt(reserved) {
				return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
			}
		} else {
			// There is either no limit in reserved pool (`None`),
			// or we are below the limit.
		}
	}

	return allConsumedWeight, nil
}

func (cw CheckWeight) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckWeight,
			"CheckWeight",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_weight", "CheckWeight"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
		),
		primitives.NewMetadataSignedExtension("CheckWeight", metadata.CheckWeight, metadata.TypesEmptyTuple)
}
