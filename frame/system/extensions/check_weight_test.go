package extensions

import (
	"bytes"
	"errors"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	dispatchInfo = &primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(1, 2),
		Class:   primitives.NewDispatchClassNormal(),
		PaysFee: primitives.PaysYes,
	}
	blockWeight = primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(7, 8),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}
	blockLength = primitives.BlockLength{
		Max: primitives.PerDispatchClass[sc.U32]{
			Normal:      10,
			Operational: 20,
			Mandatory:   30,
		},
	}
	length         = sc.U32(5)
	storageLen     = sc.U32(1)
	consumedWeight = primitives.ConsumedWeight{
		Normal:      primitives.WeightFromParts(2, 2),
		Operational: primitives.WeightFromParts(3, 3),
		Mandatory:   primitives.WeightFromParts(4, 4),
	}
)

var (
	invalidTransactionExhaustsResources = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
)

func Test_CheckWeight_AdditionalSigned(t *testing.T) {
	target := setupCheckWeight()

	result, err := target.AdditionalSigned()

	assert.Nil(t, err)
	assert.Equal(t, primitives.AdditionalSigned{}, result)
}

func Test_CheckWeight_Encode(t *testing.T) {
	target := setupCheckWeight()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, 0, buffer.Len())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_CheckWeight_Decode(t *testing.T) {
	target := setupCheckWeight()
	value := []byte{1, 2, 3}
	buffer := bytes.NewBuffer(value)

	target.Decode(buffer)

	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, bytes.NewBuffer(value), buffer)
}

func Test_CheckWeight_Bytes(t *testing.T) {
	target := setupCheckWeight()

	result := target.Bytes()

	assert.Equal(t, []byte(nil), result)
}

func Test_CheckCheckWeight_DeepCopy(t *testing.T) {
	target := setupCheckWeight()

	result := target.DeepCopy()

	assert.Equal(t, &target, result)

	target.typesInfoAdditionalSignedData = nil
	assert.NotEqual(t, &target, result)
}

func Test_CheckWeight_Validate(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)

	result, err := target.Validate(oneAccountId, nil, dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_ValidateUnsigned(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)

	result, err := target.ValidateUnsigned(nil, dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_PreDispatch(t *testing.T) {
	target := setupCheckWeight()
	expectNewStorageWeight := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockModule.On("StorageAllExtrinsicsLenSet", length+storageLen).Return()
	mockModule.On("StorageBlockWeightSet", expectNewStorageWeight).Return()

	result, err := target.PreDispatch(oneAccountId, nil, dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, primitives.Pre{}, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLenSet", length+storageLen)
	mockModule.AssertCalled(t, "StorageBlockWeightSet", expectNewStorageWeight)
}

func Test_CheckWeight_PreDispatchUnsigned(t *testing.T) {
	target := setupCheckWeight()
	expectNewStorageWeight := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockModule.On("StorageAllExtrinsicsLenSet", length+storageLen).Return()
	mockModule.On("StorageBlockWeightSet", expectNewStorageWeight).Return()

	result := target.PreDispatchUnsigned(nil, dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLenSet", length+storageLen)
	mockModule.AssertCalled(t, "StorageBlockWeightSet", expectNewStorageWeight)
}

func Test_CheckWeight_PostDispatch_Unspent(t *testing.T) {
	postInfo := &primitives.PostDispatchInfo{
		ActualWeight: sc.NewOption[primitives.Weight](primitives.WeightFromParts(3, 1)),
	}
	expectedStorageWeight := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Sub(primitives.WeightFromParts(0, 1)),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}
	target := setupCheckWeight()

	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockModule.On("StorageBlockWeightSet", expectedStorageWeight).Return()

	result := target.PostDispatch(sc.Option[primitives.Pre]{}, dispatchInfo, postInfo, sc.Compact{}, nil)

	assert.Nil(t, result)

	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertCalled(t, "StorageBlockWeightSet", expectedStorageWeight)
}

func Test_CheckWeight_PostDispatch_NoUnspent(t *testing.T) {
	postInfo := &primitives.PostDispatchInfo{
		ActualWeight: sc.NewOption[primitives.Weight](primitives.WeightFromParts(1, 2)),
	}
	target := setupCheckWeight()

	result := target.PostDispatch(sc.Option[primitives.Pre]{}, dispatchInfo, postInfo, sc.Compact{}, nil)

	assert.Nil(t, result)

	mockModule.AssertNotCalled(t, "StorageBlockWeight")
}

func Test_CheckWeight_PostDispatch_StorageBlockWeightError(t *testing.T) {
	postInfo := &primitives.PostDispatchInfo{
		ActualWeight: sc.NewOption[primitives.Weight](primitives.WeightFromParts(3, 1)),
	}
	expectedStorageWeight := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Sub(primitives.WeightFromParts(0, 1)),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}
	target := setupCheckWeight()

	expectedErr := errors.New("error")
	mockModule.On("StorageBlockWeight").Return(consumedWeight, expectedErr)
	mockModule.On("StorageBlockWeightSet", expectedStorageWeight).Return()

	result := target.PostDispatch(sc.Option[primitives.Pre]{}, dispatchInfo, postInfo, sc.Compact{}, nil)

	assert.Equal(t, expectedErr, result)

	mockModule.AssertCalled(t, "StorageBlockWeight")
}

func Test_CheckWeight_PostDispatch_Reduce_Error(t *testing.T) {
	dispatchInfoInvalidClass := &primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(1, 2),
		Class:   primitives.DispatchClass{VaryingData: sc.NewVaryingData(sc.U8(99))},
		PaysFee: primitives.PaysYes,
	}
	postInfo := &primitives.PostDispatchInfo{
		ActualWeight: sc.NewOption[primitives.Weight](primitives.WeightFromParts(3, 1)),
	}
	target := setupCheckWeight()

	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)

	err := target.PostDispatch(sc.Option[primitives.Pre]{}, dispatchInfoInvalidClass, postInfo, sc.Compact{}, nil)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())

	mockModule.AssertCalled(t, "StorageBlockWeight")
}

func Test_CheckWeight_doValidate_Success(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)

	result, err := target.doValidate(dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, primitives.DefaultValidTransaction(), result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_doValidate_InvalidExtrinsicLength(t *testing.T) {
	target := setupCheckWeight()

	blockWeight := primitives.BlockWeights{
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				MaxExtrinsic: sc.NewOption[primitives.Weight](primitives.WeightFromParts(1, 0)),
			},
		},
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)

	result, err := target.doValidate(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, primitives.ValidTransaction{}, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_doValidate_InvalidBlockLength(t *testing.T) {
	storageLen := sc.U32(10)
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result, err := target.doValidate(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, primitives.ValidTransaction{}, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertNotCalled(t, "BlockWeights")
}

func Test_CheckWeight_doPreDispatch_Success(t *testing.T) {
	target := setupCheckWeight()
	expectNewStorageWeight := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockModule.On("StorageAllExtrinsicsLenSet", length+storageLen).Return()
	mockModule.On("StorageBlockWeightSet", expectNewStorageWeight).Return()

	result := target.doPreDispatch(dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertCalled(t, "BlockWeights")
	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLenSet", length+storageLen)
	mockModule.AssertCalled(t, "StorageBlockWeightSet", expectNewStorageWeight)
}

func Test_CheckWeight_doPreDispatch_InvalidExtrinsicWeight(t *testing.T) {
	target := setupCheckWeight()

	blockWeight := primitives.BlockWeights{
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				MaxExtrinsic: sc.NewOption[primitives.Weight](primitives.WeightFromParts(1, 0)),
			},
		},
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)

	result := target.doPreDispatch(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertNumberOfCalls(t, "BlockWeights", 2)
	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertNotCalled(t, "StorageAllExtrinsicsLenSet", mock.Anything)
	mockModule.AssertNotCalled(t, "StorageBlockWeightSet", mock.Anything)
}

func Test_CheckWeight_doPreDispatch_InvalidBlockWeight(t *testing.T) {
	target := setupCheckWeight()
	blockWeight := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(math.MaxUint64, 8),
				MaxTotal:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(10, 12)),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)
	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)

	result := target.doPreDispatch(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertNumberOfCalls(t, "BlockWeights", 1)
	mockModule.AssertCalled(t, "StorageBlockWeight")
	mockModule.AssertNotCalled(t, "StorageAllExtrinsicsLenSet", mock.Anything)
	mockModule.AssertNotCalled(t, "StorageBlockWeightSet", mock.Anything)
}

func Test_CheckWeight_doPreDispatch_InvalidBlockLength(t *testing.T) {
	storageLen := sc.U32(10)
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result := target.doPreDispatch(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
	mockModule.AssertNotCalled(t, "BlockWeights")
	mockModule.AssertNotCalled(t, "StorageBlockWeight")
	mockModule.AssertNotCalled(t, "StorageAllExtrinsicsLenSet", mock.Anything)
	mockModule.AssertNotCalled(t, "StorageBlockWeightSet", mock.Anything)
}

func Test_CheckWeight_checkBlockLength_DispatchNormal(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, length+storageLen, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockLength_DispatchOperational(t *testing.T) {
	dispatchInfo := &primitives.DispatchInfo{
		Class: primitives.NewDispatchClassOperational(),
	}
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, length+storageLen, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockLength_DispatchMandatory(t *testing.T) {
	dispatchInfo := &primitives.DispatchInfo{
		Class: primitives.NewDispatchClassMandatory(),
	}
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))

	assert.Nil(t, err)
	assert.Equal(t, length+storageLen, result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockLength_InvalidDispatch(t *testing.T) {
	dispatchInfo := &primitives.DispatchInfo{
		Class: primitives.DispatchClass{
			VaryingData: sc.NewVaryingData(sc.U8(4)),
		},
	}
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	_, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))
	assert.Equal(t, errInvalidDispatchClass, err)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockLength_ExhaustsResources(t *testing.T) {
	storageLen := sc.U32(10)
	target := setupCheckWeight()

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, nil)

	result, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, sc.U32(0), result)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockLength_StorageAllExtrinsicsLenError(t *testing.T) {
	storageLen := sc.U32(10)
	target := setupCheckWeight()
	expectedErr := errors.New("error")

	mockModule.On("BlockLength").Return(blockLength)
	mockModule.On("StorageAllExtrinsicsLen").Return(storageLen, expectedErr)

	_, err := target.checkBlockLength(dispatchInfo, sc.ToCompact(length))

	assert.Equal(t, expectedErr, err)

	mockModule.AssertCalled(t, "BlockLength")
	mockModule.AssertCalled(t, "StorageAllExtrinsicsLen")
}

func Test_CheckWeight_checkBlockWeight(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, nil)

	expect := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	result, err := target.checkBlockWeight(dispatchInfo)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)

	mockModule.AssertCalled(t, "BlockWeights")
	mockModule.AssertCalled(t, "StorageBlockWeight")
}

func Test_CheckWeight_checkBlockWeight_StorageBlockWeightError(t *testing.T) {
	target := setupCheckWeight()
	expectedErr := errors.New("error")

	mockModule.On("BlockWeights").Return(blockWeight)
	mockModule.On("StorageBlockWeight").Return(consumedWeight, expectedErr)

	_, err := target.checkBlockWeight(dispatchInfo)

	assert.Equal(t, expectedErr, err)

	mockModule.AssertCalled(t, "BlockWeights")
	mockModule.AssertCalled(t, "StorageBlockWeight")
}

func Test_CheckWeight_checkExtrinsicWeight_ExhaustsResources(t *testing.T) {
	target := setupCheckWeight()

	blockWeight := primitives.BlockWeights{
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				MaxExtrinsic: sc.NewOption[primitives.Weight](primitives.WeightFromParts(1, 0)),
			},
		},
	}

	mockModule.On("BlockWeights").Return(blockWeight)

	result := target.checkExtrinsicWeight(dispatchInfo)

	assert.Equal(t, invalidTransactionExhaustsResources, result)

	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_checkExtrinsicWeight_NoMax(t *testing.T) {
	target := setupCheckWeight()

	mockModule.On("BlockWeights").Return(blockWeight)

	result := target.checkExtrinsicWeight(dispatchInfo)

	assert.Nil(t, result)
	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_checkExtrinsicWeight_BlockWeights_Get_Error(t *testing.T) {
	target := setupCheckWeight()

	dispatchInfoInvalidClass := &primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(1, 2),
		Class:   primitives.DispatchClass{VaryingData: sc.NewVaryingData()},
		PaysFee: primitives.PaysYes,
	}

	mockModule.On("BlockWeights").Return(blockWeight)

	err := target.checkExtrinsicWeight(dispatchInfoInvalidClass)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())

	mockModule.AssertCalled(t, "BlockWeights")
}

func Test_CheckWeight_calculateConsumedWeight_Success(t *testing.T) {
	target := setupCheckWeight()

	expect := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	result, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfo)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
}

func Test_CheckWeight_calculateConsumedWeight_MaxTotal_ExhaustsResources(t *testing.T) {
	target := setupCheckWeight()
	blockWeight := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(math.MaxUint64, 8),
				MaxTotal:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(10, 12)),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}

	result, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfo)

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, primitives.ConsumedWeight{}, result)
}

func Test_CheckWeight_calculateConsumsedWeight_TotalMoreThanMaxReserved_ExhaustsResources(t *testing.T) {
	target := setupCheckWeight()
	blockWeight := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(7, 8),
				MaxTotal:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(10, 12)),
				Reserved:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(0, 1)),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}

	result, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfo)

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, primitives.ConsumedWeight{}, result)
}

func Test_CheckWeight_calculateConsumedWeight_LessThanMaxTotal_ExhaustsResources(t *testing.T) {
	target := setupCheckWeight()
	blockWeight := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(7, 8),
				MaxTotal:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(0, 1)),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}

	result, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfo)

	assert.Equal(t, invalidTransactionExhaustsResources, err)
	assert.Equal(t, primitives.ConsumedWeight{}, result)
}

func Test_CheckWeight_calculateConsumedWeight_MaxTotal_Success(t *testing.T) {
	target := setupCheckWeight()
	blockWeight := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(3, 4),
		MaxBlock:  primitives.WeightFromParts(5, 6),
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(7, 8),
				MaxTotal:      sc.NewOption[primitives.Weight](primitives.WeightFromParts(10, 12)),
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(9, 10),
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.WeightFromParts(11, 12),
			},
		},
	}

	expect := primitives.ConsumedWeight{
		Normal:      consumedWeight.Normal.Add(dispatchInfo.Weight).Add(blockWeight.PerClass.Normal.BaseExtrinsic),
		Operational: consumedWeight.Operational,
		Mandatory:   consumedWeight.Mandatory,
	}

	result, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfo)

	assert.Nil(t, err)
	assert.Equal(t, expect, result)
}

func Test_CheckWeight_calculateConsumedWeight_MaximumWeight_Get_Error(t *testing.T) {
	target := setupCheckWeight()

	dispatchInfoInvalidClass := &primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(1, 2),
		Class:   primitives.DispatchClass{VaryingData: sc.NewVaryingData()},
		PaysFee: primitives.PaysYes,
	}
	_, err := target.calculateConsumedWeight(blockWeight, consumedWeight, dispatchInfoInvalidClass)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
}

func Test_CheckWeight_maxLimit_ClassIs_Error(t *testing.T) {
	dispatchInfoInvalidClass := &primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(1, 2),
		Class:   primitives.DispatchClass{VaryingData: sc.NewVaryingData()},
		PaysFee: primitives.PaysYes,
	}
	_, err := maxLimit(blockLength, dispatchInfoInvalidClass)
	assert.Equal(t, "not a valid 'DispatchClass' type", err.Error())
}

func Test_CheckWeight_ModulePath(t *testing.T) {
	target := setupCheckWeight()

	expectedModulePath := "frame_system"
	actualModulePath := target.ModulePath()

	assert.Equal(t, expectedModulePath, actualModulePath)
}

func setupCheckWeight() CheckWeight {
	mockModule = new(mocks.SystemModule)
	extension, ok := NewCheckWeight(mockModule).(*CheckWeight)
	if !ok {
		panic("invalid type assert for *CheckWeight")
	}
	return *extension
}
