package executive

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	runtimeVersion = &primitives.RuntimeVersion{
		SpecName:           "new-version",
		ImplName:           "new-version",
		AuthoringVersion:   1,
		SpecVersion:        100,
		ImplVersion:        1,
		TransactionVersion: 1,
		StateVersion:       1,
	}

	oldUpgradeInfo = primitives.LastRuntimeUpgradeInfo{
		SpecVersion: 99,
		SpecName:    "old-version",
	}

	currentUpgradeInfo = primitives.LastRuntimeUpgradeInfo{
		SpecVersion: 100,
		SpecName:    "new-version",
	}

	blockWeights = primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(1, 1),
		MaxBlock:  primitives.WeightFromParts(7, 7),
	}

	consumedWeight = primitives.ConsumedWeight{
		Normal:      primitives.WeightFromParts(1, 1),
		Operational: primitives.WeightFromParts(2, 2),
		Mandatory:   primitives.WeightFromParts(3, 3),
	}

	baseWeight = primitives.WeightFromParts(1, 1)

	dispatchClassNormal      = primitives.NewDispatchClassNormal()
	dispatchClassOperational = primitives.NewDispatchClassOperational()
	dispatchClassMandatory   = primitives.NewDispatchClassMandatory()

	dispatchInfo = primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(2, 2),
		Class:   dispatchClassNormal,
		PaysFee: primitives.NewPaysYes(),
	}

	dispatchResultWithPostInfo = &primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	unsignedValidator primitives.UnsignedValidator

	defaultAccountIdLookup = primitives.DefaultAccountIdLookup()

	txSource = primitives.NewTransactionSourceExternal()

	defaultDigest = primitives.Digest{}

	blockNumber = sc.U64(1)

	blake256Hash = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	blockHash    = primitives.NewBlake2bHash(sc.BytesToSequenceU8(blake256Hash)...)

	header = primitives.Header{
		Number:     blockNumber,
		ParentHash: blockHash,
		Digest:     testDigest(),
	}

	block = types.NewBlock(header, sc.Sequence[primitives.UncheckedExtrinsic]{})

	encodedExtrinsic    = []byte{0, 1, 2, 3, 4, 5}
	encodedExtrinsicLen = sc.ToCompact(len(encodedExtrinsic))

	signer = sc.Option[primitives.Address32]{}
)

var (
	transactionValidityError primitives.TransactionValidityError

	unknownTransactionCannotLookupError = primitives.NewTransactionValidityError(
		primitives.NewUnknownTransactionCannotLookup(),
	)

	invalidTransactionExhaustsResourcesError = primitives.NewTransactionValidityError(
		primitives.NewInvalidTransactionExhaustsResources(),
	)

	invalidTransactionBadMandatory = primitives.NewTransactionValidityError(
		primitives.NewInvalidTransactionBadMandatory(),
	)

	invalidTransactionMandatoryValidation = primitives.NewTransactionValidityError(
		primitives.NewInvalidTransactionMandatoryValidation(),
	)

	defaultDispatchOutcome  = primitives.DispatchOutcome{}
	defaultValidTransaction = primitives.ValidTransaction{}
)

var (
	target module

	mockSystemModule                  *mocks.SystemModule
	mockRuntimeExtrinsic              *mocks.RuntimeExtrinsic
	mockOnRuntimeUpgradeHook          *mocks.DefaultOnRuntimeUpgrade
	mockUncheckedExtrinsic            *mocks.UncheckedExtrinsic
	mockSignedExtra                   *mocks.SignedExtra
	mockCheckedExtrinsic              *mocks.CheckedExtrinsic
	mockCall                          *mocks.Call
	mockStorageLastRuntimeUpgradeInfo *mocks.StorageValue[primitives.LastRuntimeUpgradeInfo]
	mockStorageBlockHash              *mocks.StorageMap[sc.U64, primitives.Blake2bHash]
	mockStorageBlockNumber            *mocks.StorageValue[sc.U64]
	mockStorageBlockWeight            *mocks.StorageValue[primitives.ConsumedWeight]
	mockIoHashing                     *mocks.IoHashing
)

func setup() {
	mockSystemModule = new(mocks.SystemModule)
	mockRuntimeExtrinsic = new(mocks.RuntimeExtrinsic)
	mockOnRuntimeUpgradeHook = new(mocks.DefaultOnRuntimeUpgrade)
	mockUncheckedExtrinsic = new(mocks.UncheckedExtrinsic)
	mockSignedExtra = new(mocks.SignedExtra)
	mockCheckedExtrinsic = new(mocks.CheckedExtrinsic)
	mockCall = new(mocks.Call)
	mockStorageLastRuntimeUpgradeInfo = new(mocks.StorageValue[primitives.LastRuntimeUpgradeInfo])
	mockStorageBlockHash = new(mocks.StorageMap[sc.U64, primitives.Blake2bHash])
	mockStorageBlockNumber = new(mocks.StorageValue[sc.U64])
	mockStorageBlockWeight = new(mocks.StorageValue[primitives.ConsumedWeight])
	mockIoHashing = new(mocks.IoHashing)

	target = New(
		mockSystemModule,
		mockRuntimeExtrinsic,
		mockOnRuntimeUpgradeHook,
	).(module)
	target.hashing = mockIoHashing

	unsignedValidator = extrinsic.NewUnsignedValidatorForChecked(mockRuntimeExtrinsic)
}

func testDigest() primitives.Digest {
	digest := primitives.Digest{}

	digest[primitives.DigestTypePreRuntime] = sc.FixedSequence[primitives.DigestItem]{
		primitives.DigestItem{
			Engine:  sc.BytesToFixedSequenceU8([]byte{'a', 'u', 'r', 'a'}),
			Payload: sc.BytesToSequenceU8(sc.U64(0).Bytes()),
		},
	}

	return digest
}

func Test_Executive_InitializeBlock_VersionUpgraded(t *testing.T) {
	setup()

	mockSystemModule.On("ResetEvents").Return()
	mockSystemModule.On("StorageLastRuntimeUpgrade").Return(oldUpgradeInfo, nil)
	mockSystemModule.On("Version").Return(*runtimeVersion)
	mockSystemModule.On("StorageLastRuntimeUpgradeSet", currentUpgradeInfo)
	mockOnRuntimeUpgradeHook.On("OnRuntimeUpgrade").Return(primitives.WeightFromParts(1, 1))
	mockRuntimeExtrinsic.On("OnRuntimeUpgrade").Return(primitives.WeightFromParts(2, 2))
	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.On("OnInitialize", header.Number).Return(primitives.WeightFromParts(3, 3), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("RegisterExtraWeightUnchecked", primitives.WeightFromParts(7, 7), dispatchClassMandatory)
	mockSystemModule.On("NoteFinishedInitialize")

	target.InitializeBlock(header)

	mockSystemModule.AssertCalled(t, "ResetEvents")
	mockSystemModule.AssertCalled(t, "StorageLastRuntimeUpgradeSet", currentUpgradeInfo)
	mockOnRuntimeUpgradeHook.AssertCalled(t, "OnRuntimeUpgrade")
	mockRuntimeExtrinsic.AssertCalled(t, "OnRuntimeUpgrade")
	mockSystemModule.AssertCalled(t, "Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.AssertCalled(t, "OnInitialize", header.Number)
	mockSystemModule.AssertCalled(t, "RegisterExtraWeightUnchecked", primitives.WeightFromParts(7, 7), dispatchClassMandatory)
	mockSystemModule.AssertCalled(t, "NoteFinishedInitialize")
}

func Test_Executive_InitializeBlock_VersionNotUpgraded(t *testing.T) {
	setup()

	mockSystemModule.On("ResetEvents").Return()
	mockSystemModule.On("StorageLastRuntimeUpgrade").Return(currentUpgradeInfo, nil)
	mockSystemModule.On("Version").Return(*runtimeVersion)
	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.On("OnInitialize", header.Number).Return(primitives.WeightFromParts(3, 3), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.On("NoteFinishedInitialize")

	target.InitializeBlock(header)

	mockSystemModule.AssertCalled(t, "ResetEvents")
	mockStorageLastRuntimeUpgradeInfo.AssertNotCalled(t, "Put", currentUpgradeInfo)
	mockOnRuntimeUpgradeHook.AssertNotCalled(t, "OnRuntimeUpgrade")
	mockRuntimeExtrinsic.AssertNotCalled(t, "OnRuntimeUpgrade")
	mockSystemModule.AssertCalled(t, "Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.AssertCalled(t, "OnInitialize", header.Number)
	mockSystemModule.AssertCalled(t, "RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.AssertCalled(t, "NoteFinishedInitialize")
}

func Test_Executive_ExecuteBlock_InvalidParentHash(t *testing.T) {
	setup()

	mockSystemModule.On("ResetEvents").Return()
	mockSystemModule.On("StorageLastRuntimeUpgrade").Return(currentUpgradeInfo, nil)
	mockSystemModule.On("Version").Return(*runtimeVersion)
	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.On("OnInitialize", header.Number).Return(primitives.WeightFromParts(3, 3), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.On("NoteFinishedInitialize")

	invalidParentHash := primitives.NewBlake2bHash(sc.BytesToSequenceU8([]byte("abcdefghijklmnopqrstuvwxyz123450"))...)
	mockSystemModule.On("StorageBlockHash", header.Number-1).Return(invalidParentHash, nil)

	assert.PanicsWithValue(t, "parent hash should be valid", func() {
		target.ExecuteBlock(block)
	})
}

func Test_Executive_ExecuteBlock_InvalidInherentPosition(t *testing.T) {
	setup()

	header := primitives.Header{
		Number:     sc.U64(0),
		ParentHash: blockHash,
		Digest:     testDigest(),
	}

	block := types.NewBlock(header, sc.Sequence[primitives.UncheckedExtrinsic]{})

	mockSystemModule.On("ResetEvents").Return()
	mockSystemModule.On("StorageLastRuntimeUpgrade").Return(currentUpgradeInfo, nil)
	mockSystemModule.On("Version").Return(*runtimeVersion)
	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.On("OnInitialize", header.Number).Return(primitives.WeightFromParts(3, 3), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.On("NoteFinishedInitialize")
	mockRuntimeExtrinsic.On("EnsureInherentsAreFirst", block).Return(0)

	assert.PanicsWithValue(t, "invalid inherent position for extrinsic at index 0", func() {
		target.ExecuteBlock(block)
	})
}

func Test_Executive_ExecuteBlock_Success(t *testing.T) {
	setup()

	blockWeights := primitives.BlockWeights{
		BaseBlock: primitives.WeightFromParts(1, 1),
		MaxBlock:  primitives.WeightFromParts(6, 6),
	}
	header := primitives.Header{
		Number:     sc.U64(0),
		ParentHash: blockHash,
		Digest:     testDigest(),
	}

	block := types.NewBlock(header, sc.Sequence[primitives.UncheckedExtrinsic]{})

	mockSystemModule.On("ResetEvents").Return()
	mockSystemModule.On("StorageLastRuntimeUpgrade").Return(currentUpgradeInfo, nil)
	mockSystemModule.On("Version").Return(*runtimeVersion)
	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockRuntimeExtrinsic.On("OnInitialize", header.Number).Return(primitives.WeightFromParts(3, 3), nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockSystemModule.On("RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.On("NoteFinishedInitialize")
	mockSystemModule.On("StorageBlockHash", header.Number-1).Return(blockHash, nil)
	mockRuntimeExtrinsic.On("EnsureInherentsAreFirst", block).Return(-1)
	mockSystemModule.On("NoteFinishedExtrinsics")
	mockSystemModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	mockRuntimeExtrinsic.On("OnFinalize", blockNumber-1).Return(nil)
	mockSystemModule.On("Finalize").Return(header, nil)

	target.ExecuteBlock(block)

	mockSystemModule.AssertCalled(t, "ResetEvents")
	mockSystemModule.AssertNotCalled(t, "StorageLastRuntimeUpgradeSet", currentUpgradeInfo)
	mockOnRuntimeUpgradeHook.AssertNotCalled(t, "OnRuntimeUpgrade")
	mockRuntimeExtrinsic.AssertNotCalled(t, "OnRuntimeUpgrade")
	mockSystemModule.AssertCalled(t, "Initialize", header.Number, header.ParentHash, header.Digest)
	mockSystemModule.AssertCalled(t, "RegisterExtraWeightUnchecked", primitives.WeightFromParts(4, 4), dispatchClassMandatory)
	mockSystemModule.AssertCalled(t, "NoteFinishedInitialize")
	mockSystemModule.AssertNotCalled(t, "StorageBlockHash", header.Number-1)
	mockRuntimeExtrinsic.AssertCalled(t, "EnsureInherentsAreFirst", block)
	mockSystemModule.AssertCalled(t, "NoteFinishedExtrinsics")
	mockSystemModule.AssertCalled(t, "StorageBlockWeight")
	mockSystemModule.AssertCalled(t, "BlockWeights")
	mockRuntimeExtrinsic.AssertNotCalled(t, "OnIdle")
	mockSystemModule.AssertNotCalled(t, "RegisterExtraWeightUnchecked")
	mockRuntimeExtrinsic.AssertCalled(t, "OnFinalize", blockNumber-1)
	mockSystemModule.AssertCalled(t, "Finalize")
}

func Test_Executive_ApplyExtrinsic_UnknownTransactionCannotLookupError(t *testing.T) {
	setup()

	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(nil, unknownTransactionCannotLookupError)

	outcome, err := target.ApplyExtrinsic(mockUncheckedExtrinsic)

	assert.Equal(t, defaultDispatchOutcome, outcome)
	assert.Equal(t, unknownTransactionCannotLookupError, err)
}

func Test_Executive_ApplyExtrinsic_InvalidTransactionExhaustsResourcesError(t *testing.T) {
	setup()

	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(mockCheckedExtrinsic, nil)
	mockSystemModule.On("NoteExtrinsic", mockUncheckedExtrinsic.Bytes())
	mockCheckedExtrinsic.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfo.Weight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfo.Class)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfo.PaysFee)

	mockCheckedExtrinsic.On("Apply", unsignedValidator, &dispatchInfo, encodedExtrinsicLen).
		Return(primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, invalidTransactionExhaustsResourcesError)

	outcome, err := target.ApplyExtrinsic(mockUncheckedExtrinsic)

	mockSystemModule.AssertCalled(t, "NoteExtrinsic", mockUncheckedExtrinsic.Bytes())
	mockSystemModule.AssertNotCalled(t, "NoteAppliedExtrinsic", mock.Anything, mock.Anything)
	assert.Equal(t, defaultDispatchOutcome, outcome)
	assert.Equal(t, invalidTransactionExhaustsResourcesError, err)
}

func Test_Executive_ApplyExtrinsic_InvalidTransactionBadMandatoryError(t *testing.T) {
	setup()

	dispatchInfo := primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(2, 2),
		Class:   dispatchClassMandatory,
		PaysFee: primitives.NewPaysYes(),
	}

	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(mockCheckedExtrinsic, nil)
	mockSystemModule.On("NoteExtrinsic", mockUncheckedExtrinsic.Bytes())
	mockCheckedExtrinsic.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfo.Weight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfo.Class)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfo.PaysFee)
	mockCheckedExtrinsic.On("Apply", unsignedValidator, &dispatchInfo, encodedExtrinsicLen).Return(*dispatchResultWithPostInfo, nil)

	outcome, err := target.ApplyExtrinsic(mockUncheckedExtrinsic)

	mockSystemModule.AssertCalled(t, "NoteExtrinsic", mockUncheckedExtrinsic.Bytes())
	mockSystemModule.AssertNotCalled(t, "NoteAppliedExtrinsic", mock.Anything, mock.Anything)
	assert.Equal(t, defaultDispatchOutcome, outcome)
	assert.Equal(t, invalidTransactionBadMandatory, err)
}

func Test_Executive_ApplyExtrinsic_Success(t *testing.T) {
	setup()

	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(mockCheckedExtrinsic, nil)
	mockSystemModule.On("NoteExtrinsic", mockUncheckedExtrinsic.Bytes())

	mockCheckedExtrinsic.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfo.Weight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfo.Class)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfo.PaysFee)
	mockCheckedExtrinsic.On("Apply", unsignedValidator, &dispatchInfo, encodedExtrinsicLen).
		Return(*dispatchResultWithPostInfo, transactionValidityError)
	mockSystemModule.On("NoteAppliedExtrinsic", dispatchResultWithPostInfo, dispatchInfo)

	outcome, err := target.ApplyExtrinsic(mockUncheckedExtrinsic)

	mockSystemModule.AssertCalled(t, "NoteExtrinsic", mockUncheckedExtrinsic.Bytes())
	mockSystemModule.AssertCalled(t, "NoteAppliedExtrinsic", dispatchResultWithPostInfo, dispatchInfo)
	assert.Equal(t, primitives.NewDispatchOutcome(dispatchResultWithPostInfo.Err.Error), outcome)
	assert.Equal(t, transactionValidityError, err)
}

func Test_Executive_FinalizeBlock(t *testing.T) {
	setup()

	blockNumber := sc.U64(3)
	header := primitives.Header{
		Number:     blockNumber,
		ParentHash: blockHash,
		Digest:     testDigest(),
	}

	mockSystemModule.On("NoteFinishedExtrinsics")
	mockSystemModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockSystemModule.On("StorageBlockWeight").Return(consumedWeight, nil)
	mockSystemModule.On("BlockWeights").Return(blockWeights)
	remainingWeight := primitives.WeightFromParts(1, 1)
	usedWeight := primitives.WeightFromParts(6, 6)
	mockRuntimeExtrinsic.On("OnIdle", blockNumber, remainingWeight).Return(usedWeight)
	mockSystemModule.On("RegisterExtraWeightUnchecked", usedWeight, dispatchClassMandatory)
	mockRuntimeExtrinsic.On("OnFinalize", blockNumber).Return(nil)
	mockSystemModule.On("Finalize").Return(header, nil)

	target.FinalizeBlock()

	mockSystemModule.AssertCalled(t, "NoteFinishedExtrinsics")
	mockSystemModule.AssertCalled(t, "StorageBlockNumber")
	mockSystemModule.AssertCalled(t, "StorageBlockWeight")
	mockSystemModule.AssertCalled(t, "BlockWeights")
	mockRuntimeExtrinsic.AssertCalled(t, "OnIdle", blockNumber, remainingWeight)
	mockSystemModule.AssertCalled(t, "RegisterExtraWeightUnchecked", usedWeight, dispatchClassMandatory)
	mockRuntimeExtrinsic.AssertCalled(t, "OnFinalize", blockNumber)
	mockSystemModule.AssertCalled(t, "Finalize")
}

func Test_Executive_ValidateTransaction_UnknownTransactionCannotLookupError(t *testing.T) {
	setup()

	mockSystemModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockSystemModule.On("Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(nil, unknownTransactionCannotLookupError)

	outcome, err := target.ValidateTransaction(txSource, mockUncheckedExtrinsic, header.ParentHash)

	mockSystemModule.AssertCalled(t, "StorageBlockNumber")
	mockSystemModule.AssertCalled(t, "Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.AssertCalled(t, "Bytes")
	mockUncheckedExtrinsic.AssertCalled(t, "Check", defaultAccountIdLookup)
	mockCheckedExtrinsic.AssertNotCalled(t, "Validate", unsignedValidator, txSource, &dispatchInfo, encodedExtrinsicLen)
	assert.Equal(t, defaultValidTransaction, outcome)
	assert.Equal(t, unknownTransactionCannotLookupError, err)
}

func Test_Executive_ValidateTransaction_InvalidTransactionMandatoryValidationError(t *testing.T) {
	setup()

	dispatchInfo := primitives.DispatchInfo{
		Weight:  primitives.WeightFromParts(2, 2),
		Class:   dispatchClassMandatory,
		PaysFee: primitives.NewPaysYes(),
	}

	mockSystemModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockSystemModule.On("Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(mockCheckedExtrinsic, nil)
	mockCheckedExtrinsic.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfo.Weight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfo.Class)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfo.PaysFee)

	outcome, err := target.ValidateTransaction(txSource, mockUncheckedExtrinsic, header.ParentHash)

	mockSystemModule.AssertCalled(t, "StorageBlockNumber")
	mockSystemModule.AssertCalled(t, "Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.AssertCalled(t, "Bytes")
	mockUncheckedExtrinsic.AssertCalled(t, "Check", defaultAccountIdLookup)
	mockCheckedExtrinsic.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockCheckedExtrinsic.AssertNotCalled(t, "Validate", unsignedValidator, txSource, &dispatchInfo, encodedExtrinsicLen)
	assert.Equal(t, defaultValidTransaction, outcome)
	assert.Equal(t, invalidTransactionMandatoryValidation, err)
}

func Test_Executive_ValidateTransaction(t *testing.T) {
	setup()

	mockSystemModule.On("StorageBlockNumber").Return(blockNumber, nil)
	mockSystemModule.On("Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.On("Bytes").Return(encodedExtrinsic)
	mockUncheckedExtrinsic.On("Check", defaultAccountIdLookup).Return(mockCheckedExtrinsic, transactionValidityError)
	mockCheckedExtrinsic.On("Function").Return(mockCall)
	mockCall.On("BaseWeight").Return(baseWeight)
	mockCall.On("WeighData", baseWeight).Return(dispatchInfo.Weight)
	mockCall.On("ClassifyDispatch", baseWeight).Return(dispatchInfo.Class)
	mockCall.On("PaysFee", baseWeight).Return(dispatchInfo.PaysFee)
	mockCheckedExtrinsic.On("Validate", unsignedValidator, txSource, &dispatchInfo, encodedExtrinsicLen).
		Return(defaultValidTransaction, nil)

	outcome, err := target.ValidateTransaction(txSource, mockUncheckedExtrinsic, header.ParentHash)

	mockSystemModule.AssertCalled(t, "StorageBlockNumber")
	mockSystemModule.AssertCalled(t, "Initialize", blockNumber+1, header.ParentHash, defaultDigest)
	mockUncheckedExtrinsic.AssertCalled(t, "Bytes")
	mockUncheckedExtrinsic.AssertCalled(t, "Check", defaultAccountIdLookup)
	mockCheckedExtrinsic.AssertCalled(t, "Function")
	mockCall.AssertCalled(t, "BaseWeight")
	mockCall.AssertCalled(t, "WeighData", baseWeight)
	mockCall.AssertCalled(t, "ClassifyDispatch", baseWeight)
	mockCall.AssertCalled(t, "PaysFee", baseWeight)
	mockCheckedExtrinsic.AssertCalled(t, "Validate", unsignedValidator, txSource, &dispatchInfo, encodedExtrinsicLen)
	assert.Equal(t, defaultValidTransaction, outcome)
	assert.Equal(t, transactionValidityError, err)
}

func Test_Executive_OffchainWorker(t *testing.T) {
	setup()

	mockSystemModule.On("Initialize", header.Number, header.ParentHash, header.Digest)
	mockIoHashing.On("Blake256", header.Bytes()).Return(blake256Hash)
	mockSystemModule.On("StorageBlockHashSet", header.Number, blockHash)
	mockRuntimeExtrinsic.On("OffchainWorker", header.Number)

	target.OffchainWorker(header)

	mockSystemModule.AssertCalled(t, "Initialize", header.Number, header.ParentHash, header.Digest)
	mockSystemModule.AssertCalled(t, "StorageBlockHashSet", header.Number, blockHash)
	mockRuntimeExtrinsic.AssertCalled(t, "OffchainWorker", header.Number)
}
