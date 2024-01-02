package extensions

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	errInvalidDispatchClass = errors.New("invalid DispatchClass type in CheckBlockLength()")
)

type CheckWeight struct {
	systemModule                  system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckWeight(systemModule system.Module) primitives.SignedExtension {
	return &CheckWeight{
		systemModule:                  systemModule,
		typesInfoAdditionalSignedData: sc.NewVaryingData(),
	}
}

func (cw CheckWeight) Encode(*bytes.Buffer) error {
	return nil
}

func (cw CheckWeight) Decode(*bytes.Buffer) error { return nil }

func (cw CheckWeight) Bytes() []byte {
	return sc.EncodedBytes(cw)
}

func (cw CheckWeight) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return primitives.AdditionalSigned{}, nil
}

func (cw CheckWeight) Validate(_who primitives.AccountId, _call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	return cw.doValidate(info, length)
}

func (cw CheckWeight) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	return cw.doValidate(info, length)
}

func (cw CheckWeight) PreDispatch(_who primitives.AccountId, _call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, error) {
	return primitives.Pre{}, cw.doPreDispatch(info, length)
}

func (cw CheckWeight) PreDispatchUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) error {
	return cw.doPreDispatch(info, length)
}

func (cw CheckWeight) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _dispatchErr error) error {
	unspent := postInfo.CalcUnspent(info)
	if unspent.AnyGt(primitives.WeightZero()) {
		currentWeight, err := cw.systemModule.StorageBlockWeight()
		if err != nil {
			return err
		}
		err = currentWeight.Reduce(unspent, info.Class)
		if err != nil {
			return err
		}
		cw.systemModule.StorageBlockWeightSet(currentWeight)
	}
	return nil
}

// Do the validate checks. This can be applied to both signed and unsigned.
//
// It only checks that the block weight and length limit will not exceed.
func (cw CheckWeight) doValidate(info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	// ignore the next length. If they return `Ok`, then it is below the limit.
	_, err := cw.checkBlockLength(info, length)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	// during validation, we skip block limit check. Since the `validate_transaction`
	// call runs on an empty block anyway, by this we prevent `on_initialize` weight
	// consumption from causing false negatives.
	err = cw.checkExtrinsicWeight(info)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	return primitives.DefaultValidTransaction(), nil
}

func (cw CheckWeight) doPreDispatch(info *primitives.DispatchInfo, length sc.Compact) error {
	nextLength, err := cw.checkBlockLength(info, length)
	if err != nil {
		return err
	}

	nextWeight, err := cw.checkBlockWeight(info)
	if err != nil {
		return err
	}

	err = cw.checkExtrinsicWeight(info)
	if err != nil {
		return err
	}

	cw.systemModule.StorageAllExtrinsicsLenSet(nextLength)
	cw.systemModule.StorageBlockWeightSet(nextWeight)

	return nil
}

// Checks if the current extrinsic can fit into the block with respect to block length limits.
//
// Upon successes, it returns the new block length as a `Result`.
func (cw CheckWeight) checkBlockLength(info *primitives.DispatchInfo, length sc.Compact) (sc.U32, error) {
	lengthLimit := cw.systemModule.BlockLength()
	currentLen, err := cw.systemModule.StorageAllExtrinsicsLen()
	if err != nil {
		return 0, err
	}
	addedLen := sc.U32(length.ToBigInt().Uint64())

	nextLen := sc.SaturatingAddU32(currentLen, addedLen)

	maxLimit, err := maxLimit(lengthLimit, info)
	if err != nil {
		return 0, err
	}

	if nextLen > maxLimit {
		return sc.U32(0), primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
	}

	return nextLen, nil
}

// Checks if the current extrinsic can fit into the block with respect to block weight limits.
//
// Upon successes, it returns the new block weight as a `Result`.
func (cw CheckWeight) checkBlockWeight(info *primitives.DispatchInfo) (primitives.ConsumedWeight, error) {
	maximumWeight := cw.systemModule.BlockWeights()
	allWeight, err := cw.systemModule.StorageBlockWeight()
	if err != nil {
		return primitives.ConsumedWeight{}, err
	}
	return cw.calculateConsumedWeight(maximumWeight, allWeight, info)
}

// Checks if the current extrinsic does not exceed the maximum weight a single extrinsic
// with given `DispatchClass` can have.
func (cw CheckWeight) checkExtrinsicWeight(info *primitives.DispatchInfo) error {
	dispatchClass, err := cw.systemModule.BlockWeights().Get(info.Class)
	if err != nil {
		return err
	}

	max := dispatchClass.MaxExtrinsic

	if max.HasValue {
		if info.Weight.AnyGt(max.Value) {
			return primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	}

	return nil
}

func (cw CheckWeight) calculateConsumedWeight(maximumWeight primitives.BlockWeights, allConsumedWeight primitives.ConsumedWeight, info *primitives.DispatchInfo) (primitives.ConsumedWeight, error) {
	limitPerClass, err := maximumWeight.Get(info.Class)
	if err != nil {
		return primitives.ConsumedWeight{}, err
	}

	extrinsicWeight := info.Weight.SaturatingAdd(limitPerClass.BaseExtrinsic)

	// add the weight. If class is unlimited, use saturating add instead of checked one.
	if !limitPerClass.MaxTotal.HasValue && !limitPerClass.Reserved.HasValue {
		allConsumedWeight.Accrue(extrinsicWeight, info.Class)
	} else {
		err := allConsumedWeight.CheckedAccrue(extrinsicWeight, info.Class)
		if err != nil {
			return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	}

	consumedPerClass, perClassErr := allConsumedWeight.Get(info.Class)
	if perClassErr != nil {
		return primitives.ConsumedWeight{}, perClassErr
	}

	// Check if we don't exceed per-class allowance
	if limitPerClass.MaxTotal.HasValue {
		max := limitPerClass.MaxTotal.Value
		if consumedPerClass.AnyGt(max) {
			return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
		}
	}

	// In cases total block weight is exceeded, we need to fall back
	// to `reserved` pool if there is any.
	total, totalWeightErr := allConsumedWeight.Total()
	if totalWeightErr != nil {
		return primitives.ConsumedWeight{}, totalWeightErr
	}

	if total.AnyGt(maximumWeight.MaxBlock) {
		if limitPerClass.Reserved.HasValue {
			reserved := limitPerClass.Reserved.Value
			if consumedPerClass.AnyGt(reserved) {
				return primitives.ConsumedWeight{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
			}
		}
	}

	return allConsumedWeight, nil
}

func maxLimit(lengthLimit primitives.BlockLength, info *primitives.DispatchInfo) (sc.U32, error) {
	isNormal, err := info.Class.Is(primitives.DispatchClassNormal)
	if err != nil {
		return 0, err
	}
	if isNormal {
		return lengthLimit.Max.Normal, nil
	}

	isOperational, err := info.Class.Is(primitives.DispatchClassOperational)
	if err != nil {
		return 0, err
	}
	if isOperational {
		return lengthLimit.Max.Operational, nil
	}

	isMandatory, err := info.Class.Is(primitives.DispatchClassMandatory)
	if err != nil {
		return 0, err
	}
	if isMandatory {
		return lengthLimit.Max.Mandatory, nil
	}

	return 0, errInvalidDispatchClass

}

func (cw CheckWeight) ModulePath() string {
	return systemModulePath
}
