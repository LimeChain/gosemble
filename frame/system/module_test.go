package system

import (
	"math"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	accountInfo = primitives.AccountInfo{
		Nonce:       1,
		Consumers:   2,
		Providers:   3,
		Sufficients: 4,
		Data: primitives.AccountData{
			Free:       sc.NewU128(5),
			Reserved:   sc.NewU128(6),
			MiscFrozen: sc.NewU128(7),
			FeeFrozen:  sc.NewU128(8),
		},
	}
	blockHashCount = sc.U64(5)
	blockWeights   = primitives.BlockWeights{
		BaseBlock: primitives.Weight{
			RefTime:   1,
			ProofSize: 2,
		},
		MaxBlock: primitives.Weight{
			RefTime:   3,
			ProofSize: 4,
		},
		PerClass: primitives.PerDispatchClass[primitives.WeightsPerClass]{
			Normal: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.Weight{
					RefTime:   5,
					ProofSize: 6,
				},
			},
			Operational: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.Weight{
					RefTime:   7,
					ProofSize: 8,
				},
			},
			Mandatory: primitives.WeightsPerClass{
				BaseExtrinsic: primitives.Weight{
					RefTime:   9,
					ProofSize: 10,
				},
			},
		},
	}
	blockLength = primitives.BlockLength{
		Max: primitives.PerDispatchClass[sc.U32]{
			Normal:      11,
			Operational: 12,
			Mandatory:   13,
		},
	}
	dbWeight = primitives.RuntimeDbWeight{
		Read:  1,
		Write: 2,
	}
	version = primitives.RuntimeVersion{
		SpecName:           "test-spec",
		ImplName:           "test-impl",
		AuthoringVersion:   1,
		SpecVersion:        2,
		ImplVersion:        3,
		TransactionVersion: 4,
		StateVersion:       5,
	}
	parentHash = primitives.Blake2bHash{
		FixedSequence: sc.BytesToFixedSequenceU8(
			common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
		)}
	blockNumber = sc.U64(5)
	digest      = testDigest()
)

var (
	mockStorageAccount            *mocks.StorageMap[primitives.PublicKey, primitives.AccountInfo]
	mockStorageBlockWeight        *mocks.StorageValue[primitives.ConsumedWeight]
	mockStorageBlockHash          *mocks.StorageMap[sc.U64, primitives.Blake2bHash]
	mockStorageBlockNumber        *mocks.StorageValue[sc.U64]
	mockStorageAllExtrinsicsLen   *mocks.StorageValue[sc.U32]
	mockStorageExtrinsicIndex     *mocks.StorageValue[sc.U32]
	mockStorageExtrinsicData      *mocks.StorageMap[sc.U32, sc.Sequence[sc.U8]]
	mockStorageExtrinsicCount     *mocks.StorageValue[sc.U32]
	mockStorageParentHash         *mocks.StorageValue[primitives.Blake2bHash]
	mockStorageDigest             *mocks.StorageValue[primitives.Digest]
	mockStorageEvents             *mocks.StorageValue[primitives.EventRecord]
	mockStorageEventCount         *mocks.StorageValue[sc.U32]
	mockStorageEventTopics        *mocks.StorageMap[primitives.H256, sc.VaryingData]
	mockStorageLastRuntimeUpgrade *mocks.StorageValue[primitives.LastRuntimeUpgradeInfo]
	mockStorageExecutionPhase     *mocks.StorageValue[primitives.ExtrinsicPhase]

	mockIoStorage *mocks.IoStorage
	mockIoTrie    *mocks.IoTrie

	mockTypeMutateAccountInfo       = mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]")
	mockTypeMutateOptionAccountInfo = mock.AnythingOfType("func(*goscale.Option[github.com/LimeChain/gosemble/primitives/types.AccountInfo]) goscale.Result[github.com/LimeChain/goscale.Encodable]")
)

func Test_Module_GetIndex(t *testing.T) {
	assert.Equal(t, sc.U8(moduleId), setupModule().GetIndex())
}

func Test_Module_Functions(t *testing.T) {
	target := setupModule()
	functions := target.Functions()

	assert.Equal(t, 1, len(functions))
}

func Test_Module_PreDispatch(t *testing.T) {
	target := setupModule()
	mockCall := new(mocks.Call)

	result, err := target.PreDispatch(mockCall)

	assert.Nil(t, err)
	assert.Equal(t, sc.Empty{}, result)
}

func Test_Module_ValidateUnsigned(t *testing.T) {
	target := setupModule()
	mockCall := new(mocks.Call)

	result, err := target.ValidateUnsigned(primitives.TransactionSource{}, mockCall)

	assert.Equal(t, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator()), err)
	assert.Equal(t, primitives.ValidTransaction{}, result)
}

func Test_Module_BlockHashCount(t *testing.T) {
	assert.Equal(t, blockHashCount, setupModule().BlockHashCount())
}

func Test_Module_BlockLength(t *testing.T) {
	assert.Equal(t, blockLength, setupModule().BlockLength())
}

func Test_Module_BlockWeights(t *testing.T) {
	assert.Equal(t, blockWeights, setupModule().BlockWeights())
}

func Test_Module_DbWeight(t *testing.T) {
	assert.Equal(t, dbWeight, setupModule().DbWeight())
}

func Test_Module_Version(t *testing.T) {
	assert.Equal(t, version, setupModule().Version())
}

func Test_Module_StorageDigest(t *testing.T) {
	target := setupModule()

	mockStorageDigest.On("Get").Return(digest)

	result := target.StorageDigest()

	assert.Equal(t, digest, result)
	mockStorageDigest.AssertCalled(t, "Get")
}

func Test_Module_StorageBlockWeight(t *testing.T) {
	blockWeight := primitives.ConsumedWeight{
		Normal:      primitives.WeightFromParts(1, 2),
		Operational: primitives.WeightFromParts(3, 4),
		Mandatory:   primitives.WeightFromParts(5, 6),
	}
	target := setupModule()

	mockStorageBlockWeight.On("Get").Return(blockWeight)

	result := target.StorageBlockWeight()

	assert.Equal(t, blockWeight, result)
	mockStorageBlockWeight.AssertCalled(t, "Get")
}

func Test_Module_StorageBlockWeightSet(t *testing.T) {
	blockWeight := primitives.ConsumedWeight{
		Normal:      primitives.WeightFromParts(1, 2),
		Operational: primitives.WeightFromParts(3, 4),
		Mandatory:   primitives.WeightFromParts(5, 6),
	}
	target := setupModule()

	mockStorageBlockWeight.On("Put", blockWeight).Return()

	target.StorageBlockWeightSet(blockWeight)

	mockStorageBlockWeight.AssertCalled(t, "Put", blockWeight)
}

func Test_Module_StorageBlockHash(t *testing.T) {
	key := sc.U64(0)
	target := setupModule()

	mockStorageBlockHash.On("Get", key).Return(parentHash)

	result := target.StorageBlockHash(key)

	assert.Equal(t, parentHash, result)
	mockStorageBlockHash.AssertCalled(t, "Get", key)
}

func Test_Module_StorageBlockHashSet(t *testing.T) {
	key := sc.U64(0)
	target := setupModule()

	mockStorageBlockHash.On("Put", key, parentHash).Return()

	target.StorageBlockHashSet(key, parentHash)

	mockStorageBlockHash.AssertCalled(t, "Put", key, parentHash)
}

func Test_Module_StorageBlockHashExists(t *testing.T) {
	key := sc.U64(0)
	target := setupModule()

	mockStorageBlockHash.On("Exists", key).Return(true)

	result := target.StorageBlockHashExists(key)

	assert.Equal(t, true, result)
	mockStorageBlockHash.AssertCalled(t, "Exists", key)
}

func Test_Module_StorageBlockNumber(t *testing.T) {
	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(blockNumber)

	result := target.StorageBlockNumber()

	assert.Equal(t, blockNumber, result)
	mockStorageBlockNumber.AssertCalled(t, "Get")
}

func Test_Module_StorageBlockNumberSet(t *testing.T) {
	target := setupModule()

	mockStorageBlockNumber.On("Put", blockNumber).Return()

	target.StorageBlockNumberSet(blockNumber)

	mockStorageBlockNumber.AssertCalled(t, "Put", blockNumber)
}

func Test_Module_StorageLastRuntimeUpgrade(t *testing.T) {
	lrui := primitives.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(sc.U32(1)),
		SpecName:    "test",
	}
	target := setupModule()

	mockStorageLastRuntimeUpgrade.On("Get").Return(lrui)

	result := target.StorageLastRuntimeUpgrade()

	assert.Equal(t, lrui, result)
	mockStorageLastRuntimeUpgrade.AssertCalled(t, "Get")
}

func Test_Module_StorageLastRuntimeUpgradeSet(t *testing.T) {
	lrui := primitives.LastRuntimeUpgradeInfo{
		SpecVersion: sc.ToCompact(sc.U32(1)),
		SpecName:    "test",
	}
	target := setupModule()

	mockStorageLastRuntimeUpgrade.On("Put", lrui).Return()

	target.StorageLastRuntimeUpgradeSet(lrui)

	mockStorageLastRuntimeUpgrade.AssertCalled(t, "Put", lrui)
}

func Test_Module_StorageAccount(t *testing.T) {
	target := setupModule()

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.StorageAccount(targetAccount.FixedSequence)

	assert.Equal(t, accountInfo, result)
	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
}

func Test_Module_StorageAccountSet(t *testing.T) {
	target := setupModule()

	mockStorageAccount.On("Put", targetAccount.FixedSequence, accountInfo).Return()

	target.StorageAccountSet(targetAccount.FixedSequence, accountInfo)

	mockStorageAccount.AssertCalled(t, "Put", targetAccount.FixedSequence, accountInfo)
}

func Test_Module_StorageAllExtrinsicLen(t *testing.T) {
	extrinsicLen := sc.U32(2)
	target := setupModule()

	mockStorageAllExtrinsicsLen.On("Get").Return(extrinsicLen)

	result := target.StorageAllExtrinsicsLen()

	assert.Equal(t, extrinsicLen, result)
	mockStorageAllExtrinsicsLen.AssertCalled(t, "Get")
}

func Test_Module_StorageAllExtrinsicLenSet(t *testing.T) {
	extrinsicLen := sc.U32(2)
	target := setupModule()

	mockStorageAllExtrinsicsLen.On("Put", extrinsicLen).Return()

	target.StorageAllExtrinsicsLenSet(extrinsicLen)

	mockStorageAllExtrinsicsLen.AssertCalled(t, "Put", extrinsicLen)
}

func Test_Module_Initialize(t *testing.T) {
	target := setupModule()
	executionPhase := primitives.NewExtrinsicPhaseInitialization()

	mockStorageExecutionPhase.On("Put", executionPhase).Return()
	mockStorageExtrinsicIndex.On("Put", sc.U32(0)).Return()
	mockStorageBlockNumber.On("Put", blockNumber).Return()
	mockStorageDigest.On("Put", digest).Return()
	mockStorageParentHash.On("Put", parentHash).Return()
	mockStorageBlockHash.On("Put", blockNumber-1, parentHash).Return()
	mockStorageBlockWeight.On("Clear").Return()

	target.Initialize(blockNumber, parentHash, digest)

	mockStorageExecutionPhase.AssertCalled(t, "Put", executionPhase)
	mockStorageExtrinsicIndex.AssertCalled(t, "Put", sc.U32(0))
	mockStorageBlockNumber.AssertCalled(t, "Put", blockNumber)
	mockStorageDigest.AssertCalled(t, "Put", digest)
	mockStorageParentHash.AssertCalled(t, "Put", parentHash)
	mockStorageBlockHash.AssertCalled(t, "Put", blockNumber-1, parentHash)
	mockStorageBlockWeight.AssertCalled(t, "Clear")
}

func Test_Module_RegisterExtraWeightUnchecked(t *testing.T) {
	blockWeight := primitives.ConsumedWeight{
		Normal:      primitives.WeightFromParts(1, 2),
		Operational: primitives.WeightFromParts(3, 4),
		Mandatory:   primitives.WeightFromParts(5, 6),
	}
	weight := primitives.WeightFromParts(7, 8)
	class := primitives.NewDispatchClassNormal()
	target := setupModule()
	expectCurrentWeight := primitives.ConsumedWeight{
		Normal:      blockWeight.Normal.Add(weight),
		Operational: blockWeight.Operational,
		Mandatory:   blockWeight.Mandatory,
	}

	mockStorageBlockWeight.On("Get").Return(blockWeight)
	mockStorageBlockWeight.On("Put", expectCurrentWeight)

	target.RegisterExtraWeightUnchecked(weight, class)

	mockStorageBlockWeight.AssertCalled(t, "Get")
	mockStorageBlockWeight.AssertCalled(t, "Put", expectCurrentWeight)
}

func Test_Module_NoteFinishedInitialize(t *testing.T) {
	executionPhase := primitives.NewExtrinsicPhaseApply(sc.U32(0))
	target := setupModule()

	mockStorageExecutionPhase.On("Put", executionPhase).Return()

	target.NoteFinishedInitialize()

	mockStorageExecutionPhase.AssertCalled(t, "Put", executionPhase)
}

func Test_Module_NoteExtrinsic(t *testing.T) {
	extrinsicBytes := []byte("test")
	extrinsicIndex := sc.U32(1)
	target := setupModule()

	mockStorageExtrinsicIndex.On("Get").Return(extrinsicIndex)
	mockStorageExtrinsicData.On("Put", extrinsicIndex, sc.BytesToSequenceU8(extrinsicBytes)).Return()

	target.NoteExtrinsic(extrinsicBytes)

	mockStorageExtrinsicIndex.AssertCalled(t, "Get")
	mockStorageExtrinsicData.AssertCalled(t, "Put", extrinsicIndex, sc.BytesToSequenceU8(extrinsicBytes))
}

func Test_Module_NoteAppliedExtrinsic_ExtrinsicSuccess(t *testing.T) {
	blockNum := sc.U64(5)
	eventCount := sc.U32(0)
	extrinsicIndex := sc.U32(1)
	extrinsicResult := &primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok: primitives.PostDispatchInfo{
			ActualWeight: sc.NewOption[primitives.Weight](nil),
			PaysFee:      primitives.PaysYes,
		},
	}
	dispatchInfo := primitives.DispatchInfo{
		Class:   primitives.NewDispatchClassNormal(),
		PaysFee: primitives.NewPaysYes(),
	}
	expectDispatchInfo := primitives.DispatchInfo{
		Weight:  blockWeights.PerClass.Normal.BaseExtrinsic,
		Class:   primitives.NewDispatchClassNormal(),
		PaysFee: primitives.NewPaysYes(),
	}
	expectEventRecord := primitives.EventRecord{
		Phase:  primitives.NewExtrinsicPhaseInitialization(),
		Event:  newEventExtrinsicSuccess(moduleId, expectDispatchInfo),
		Topics: []primitives.H256{},
	}

	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(blockNum)
	mockStorageExecutionPhase.On("Get").Return(primitives.NewExtrinsicPhaseInitialization())
	mockStorageEventCount.On("Get").Return(eventCount)
	mockStorageEventCount.On("Put", eventCount+1).Return()
	mockStorageEvents.On("Append", expectEventRecord).Return()

	mockStorageExtrinsicIndex.On("Get").Return(extrinsicIndex)
	mockStorageExtrinsicIndex.On("Put", extrinsicIndex+1).Return()
	mockStorageExecutionPhase.On("Put", primitives.NewExtrinsicPhaseApply(extrinsicIndex+1)).Return()

	target.NoteAppliedExtrinsic(extrinsicResult, dispatchInfo)

	mockStorageBlockNumber.AssertNumberOfCalls(t, "Get", 1)
	mockStorageExecutionPhase.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Put", eventCount+1)
	mockStorageEvents.AssertCalled(t, "Append", expectEventRecord)
	mockStorageEventTopics.AssertNotCalled(t, "Append")

	mockStorageExtrinsicIndex.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertCalled(t, "Put", primitives.NewExtrinsicPhaseApply(extrinsicIndex+1))
}

func Test_Module_NoteAppliedExtrinsic_ExtrinsicFailed(t *testing.T) {
	blockNum := sc.U64(5)
	eventCount := sc.U32(0)
	extrinsicIndex := sc.U32(1)
	extrinsicResult := &primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Ok:       primitives.PostDispatchInfo{},
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			PostInfo: primitives.PostDispatchInfo{},
			Error:    primitives.NewDispatchErrorCorruption(),
		},
	}
	dispatchInfo := primitives.DispatchInfo{
		Class:   primitives.NewDispatchClassNormal(),
		PaysFee: primitives.NewPaysYes(),
	}
	expectDispatchInfo := primitives.DispatchInfo{
		Weight:  blockWeights.PerClass.Normal.BaseExtrinsic,
		Class:   primitives.NewDispatchClassNormal(),
		PaysFee: primitives.NewPaysYes(),
	}
	expectEventRecord := primitives.EventRecord{
		Phase:  primitives.NewExtrinsicPhaseInitialization(),
		Event:  newEventExtrinsicFailed(moduleId, extrinsicResult.Err.Error, expectDispatchInfo),
		Topics: []primitives.H256{},
	}

	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(blockNum)
	mockStorageExecutionPhase.On("Get").Return(primitives.NewExtrinsicPhaseInitialization())
	mockStorageEventCount.On("Get").Return(eventCount)
	mockStorageEventCount.On("Put", eventCount+1).Return()
	mockStorageEvents.On("Append", expectEventRecord).Return()

	mockStorageExtrinsicIndex.On("Get").Return(extrinsicIndex)
	mockStorageExtrinsicIndex.On("Put", extrinsicIndex+1).Return()
	mockStorageExecutionPhase.On("Put", primitives.NewExtrinsicPhaseApply(extrinsicIndex+1)).Return()

	target.NoteAppliedExtrinsic(extrinsicResult, dispatchInfo)

	mockStorageBlockNumber.AssertNumberOfCalls(t, "Get", 2)
	mockStorageExecutionPhase.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Put", eventCount+1)
	mockStorageEvents.AssertCalled(t, "Append", expectEventRecord)
	mockStorageEventTopics.AssertNotCalled(t, "Append")

	mockStorageExtrinsicIndex.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertCalled(t, "Put", primitives.NewExtrinsicPhaseApply(extrinsicIndex+1))
}

func Test_Module_Finalize_RemovePreviousHash(t *testing.T) {
	target := setupModule()
	blockNumber := sc.U64(7)

	extrinsicCount := sc.U32(1)
	extrinsicDataBytes := []byte("extrinsicDataBytes")
	extrinsicRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a").ToBytes()
	expectExtrinsicRoot := primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(extrinsicRoot)}
	storageRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b").ToBytes()
	expectStorageRoot := primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(storageRoot)}
	blakeArgs := append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsicDataBytes...)

	expectResult := primitives.Header{
		ParentHash:     parentHash,
		Number:         blockNumber,
		StateRoot:      expectStorageRoot,
		ExtrinsicsRoot: expectExtrinsicRoot,
		Digest:         digest,
	}

	mockStorageExecutionPhase.On("Clear").Return()
	mockStorageAllExtrinsicsLen.On("Clear").Return()

	mockStorageBlockNumber.On("Get").Return(blockNumber)
	mockStorageParentHash.On("Get").Return(parentHash)
	mockStorageDigest.On("Get").Return(digest)
	mockStorageExtrinsicCount.On("Take").Return(extrinsicCount)
	mockStorageExtrinsicData.On("TakeBytes", sc.U32(0)).Return(extrinsicDataBytes)
	mockIoTrie.On("Blake2256OrderedRoot", blakeArgs, int32(constants.StorageVersion)).Return(extrinsicRoot)
	mockStorageBlockHash.On("Remove", sc.U64(1)).Return()
	mockIoStorage.On("Root", int32(version.StateVersion)).Return(storageRoot)

	result := target.Finalize()

	assert.Equal(t, expectResult, result)

	mockStorageExecutionPhase.AssertCalled(t, "Clear")
	mockStorageAllExtrinsicsLen.AssertCalled(t, "Clear")

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageParentHash.AssertCalled(t, "Get")
	mockStorageDigest.AssertCalled(t, "Get")
	mockStorageExtrinsicCount.AssertCalled(t, "Take")
	mockStorageExtrinsicData.AssertCalled(t, "TakeBytes", sc.U32(0))
	mockIoTrie.AssertCalled(t, "Blake2256OrderedRoot", blakeArgs, int32(constants.StorageVersion))
	mockStorageBlockHash.AssertCalled(t, "Remove", sc.U64(1))
	mockIoStorage.AssertCalled(t, "Root", int32(version.StateVersion))
}

func Test_Module_Finalize_Success(t *testing.T) {
	target := setupModule()
	extrinsicCount := sc.U32(1)
	extrinsicDataBytes := []byte("extrinsicDataBytes")
	extrinsicRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a").ToBytes()
	expectExtrinsicRoot := primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(extrinsicRoot)}
	storageRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b").ToBytes()
	expectStorageRoot := primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(storageRoot)}
	blakeArgs := append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsicDataBytes...)

	expectResult := primitives.Header{
		ParentHash:     parentHash,
		Number:         blockNumber,
		StateRoot:      expectStorageRoot,
		ExtrinsicsRoot: expectExtrinsicRoot,
		Digest:         digest,
	}

	mockStorageExecutionPhase.On("Clear").Return()
	mockStorageAllExtrinsicsLen.On("Clear").Return()

	mockStorageBlockNumber.On("Get").Return(blockNumber)
	mockStorageParentHash.On("Get").Return(parentHash)
	mockStorageDigest.On("Get").Return(digest)
	mockStorageExtrinsicCount.On("Take").Return(extrinsicCount)
	mockStorageExtrinsicData.On("TakeBytes", sc.U32(0)).Return(extrinsicDataBytes)
	mockIoTrie.On("Blake2256OrderedRoot", blakeArgs, int32(constants.StorageVersion)).Return(extrinsicRoot)
	mockIoStorage.On("Root", int32(version.StateVersion)).Return(storageRoot)

	result := target.Finalize()

	assert.Equal(t, expectResult, result)

	mockStorageExecutionPhase.AssertCalled(t, "Clear")
	mockStorageAllExtrinsicsLen.AssertCalled(t, "Clear")

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageParentHash.AssertCalled(t, "Get")
	mockStorageDigest.AssertCalled(t, "Get")
	mockStorageExtrinsicCount.AssertCalled(t, "Take")
	mockStorageExtrinsicData.AssertCalled(t, "TakeBytes", sc.U32(0))
	mockIoTrie.AssertCalled(t, "Blake2256OrderedRoot", blakeArgs, int32(constants.StorageVersion))
	mockStorageBlockHash.AssertNotCalled(t, "Remove", mock.Anything)
	mockIoStorage.AssertCalled(t, "Root", int32(version.StateVersion))
}

func Test_Module_NoteFinishedExtrinsics(t *testing.T) {
	extrinsicIndex := sc.U32(4)
	target := setupModule()

	mockStorageExtrinsicIndex.On("Take").Return(extrinsicIndex)
	mockStorageExtrinsicCount.On("Put", extrinsicIndex).Return()
	mockStorageExecutionPhase.On("Put", primitives.NewExtrinsicPhaseFinalization()).Return()

	target.NoteFinishedExtrinsics()

	mockStorageExtrinsicIndex.AssertCalled(t, "Take")
	mockStorageExtrinsicCount.AssertCalled(t, "Put", extrinsicIndex)
	mockStorageExecutionPhase.AssertCalled(t, "Put", primitives.NewExtrinsicPhaseFinalization())
}

func Test_Module_ResetEvents(t *testing.T) {
	target := setupModule()

	mockStorageEvents.On("Clear").Return()
	mockStorageEventCount.On("Clear").Return()
	mockStorageEventTopics.On("Clear", sc.U32(math.MaxUint32))

	target.ResetEvents()

	mockStorageEvents.AssertCalled(t, "Clear")
	mockStorageEventCount.AssertCalled(t, "Clear")
	mockStorageEventTopics.AssertCalled(t, "Clear", sc.U32(math.MaxUint32))
}

func Test_Module_CanDecProviders_ZeroConsumer(t *testing.T) {
	target := setupModule()
	accountInfo := primitives.AccountInfo{}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.CanDecProviders(targetAccount)
	assert.Equal(t, true, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
}

func Test_Module_CanDecProviders_Providers(t *testing.T) {
	target := setupModule()
	accountInfo := primitives.AccountInfo{
		Consumers: 2,
		Providers: 3,
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.CanDecProviders(targetAccount)
	assert.Equal(t, true, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
}

func Test_Module_CanDecProviders_False(t *testing.T) {
	target := setupModule()
	accountInfo := primitives.AccountInfo{
		Consumers: 2,
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.CanDecProviders(targetAccount)
	assert.Equal(t, false, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
}

func Test_Module_TryMutateExists_Error(t *testing.T) {
	target := setupModule()
	expectResult := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorBadOrigin(),
	}

	accountInfo := primitives.AccountInfo{}
	f := func(account *primitives.AccountData) sc.Result[sc.Encodable] {
		return expectResult
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.AssertNotCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_TryMutateExists_NoProviding(t *testing.T) {
	target := setupModule()
	expectResult := sc.Result[sc.Encodable]{
		Value: sc.NewU128(5),
	}

	accountInfo := primitives.AccountInfo{}
	f := func(account *primitives.AccountData) sc.Result[sc.Encodable] {
		return expectResult
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.AssertNotCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_TryMutateExists_WasProviding_NoLongerProviding_DecRefStatus_Success(t *testing.T) {
	target := setupModule()
	mockResult := sc.Result[sc.Encodable]{
		Value: primitives.DecRefStatusExists,
	}
	expectResult := sc.Result[sc.Encodable]{
		Value: sc.NewU128(5),
	}

	accountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{
			Free: sc.NewU128(1),
		},
	}
	f := func(account *primitives.AccountData) sc.Result[sc.Encodable] {
		account.Free = primitives.Balance{}
		return expectResult
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)
	mockStorageAccount.
		On(
			"TryMutateExists",
			targetAccount.FixedSequence,
			mockTypeMutateOptionAccountInfo).
		Return(mockResult)

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.
		AssertCalled(t,
			"TryMutateExists",
			targetAccount.FixedSequence,
			mockTypeMutateOptionAccountInfo)
	mockStorageAccount.AssertNotCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_TryMutateExists_WasProviding_NoLongerProviding_Error(t *testing.T) {
	target := setupModule()
	expectError := primitives.NewDispatchErrorCannotLookup()
	mockResult := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectError,
	}
	expectResult := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectError,
	}

	accountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{
			Free: sc.NewU128(1),
		},
	}
	f := func(account *primitives.AccountData) sc.Result[sc.Encodable] {
		account.Free = primitives.Balance{}
		return sc.Result[sc.Encodable]{}
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)
	mockStorageAccount.
		On(
			"TryMutateExists",
			targetAccount.FixedSequence,
			mockTypeMutateOptionAccountInfo).
		Return(mockResult)

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.
		AssertCalled(t,
			"TryMutateExists",
			targetAccount.FixedSequence,
			mockTypeMutateOptionAccountInfo)
	mockStorageAccount.AssertNotCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_TryMutateExists_WasNotProviding_IsProviding(t *testing.T) {
	target := setupModule()

	expectResult := sc.Result[sc.Encodable]{
		Value: sc.NewU128(5),
	}
	accountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{},
	}
	f := func(account *primitives.AccountData) sc.Result[sc.Encodable] {
		account.Free = sc.NewU128(5)
		return expectResult
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)
	mockStorageAccount.On(
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo).
		Return(sc.Result[sc.Encodable]{Value: primitives.IncRefStatusExisted}).Once()
	mockStorageAccount.On(
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo).
		Return(sc.Result[sc.Encodable]{Value: sc.NewU128(2)}).Once()

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.AssertNumberOfCalls(t, "Mutate", 2)
	mockStorageAccount.AssertCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_TryMutateExists_WasProviding_IsProviding_Success(t *testing.T) {
	target := setupModule()

	expectResult := sc.Result[sc.Encodable]{
		Value: sc.NewU128(5),
	}
	accountInfo := primitives.AccountInfo{
		Data: primitives.AccountData{
			Free: sc.NewU128(1),
		},
	}
	f := func(*primitives.AccountData) sc.Result[sc.Encodable] {
		return expectResult
	}

	mockStorageAccount.On("Get", targetAccount.FixedSequence).Return(accountInfo)
	mockStorageAccount.On(
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo).
		Return(sc.Result[sc.Encodable]{})

	result := target.TryMutateExists(targetAccount, f)

	assert.Equal(t, expectResult, result)

	mockStorageAccount.AssertCalled(t, "Get", targetAccount.FixedSequence)
	mockStorageAccount.AssertNumberOfCalls(t, "Mutate", 1)
	mockStorageAccount.AssertCalled(t,
		"Mutate",
		targetAccount.FixedSequence,
		mockTypeMutateAccountInfo)
}

func Test_Module_incrementProviders_RefStatusCreated(t *testing.T) {
	accountInfo := &primitives.AccountInfo{}
	expect := sc.Result[sc.Encodable]{
		HasError: false,
		Value:    primitives.IncRefStatusCreated,
	}
	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(sc.U64(0))

	result := target.incrementProviders(targetAccount, accountInfo)

	assert.Equal(t, expect, result)
	assert.Equal(t, sc.U32(1), accountInfo.Providers)

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Put", mock.Anything)
	mockStorageEventCount.AssertNotCalled(t, "Append", mock.Anything)
	mockStorageEventTopics.AssertNotCalled(t, "Append", mock.Anything, mock.Anything)
}

func Test_Module_incrementProviders_RefStatusExisted(t *testing.T) {
	accountInfo := &primitives.AccountInfo{
		Sufficients: 1,
	}
	expect := sc.Result[sc.Encodable]{
		HasError: false,
		Value:    primitives.IncRefStatusExisted,
	}
	target := setupModule()

	result := target.incrementProviders(targetAccount, accountInfo)

	assert.Equal(t, expect, result)
	assert.Equal(t, sc.U32(1), accountInfo.Providers)

	mockStorageBlockNumber.AssertNotCalled(t, "Get")
	mockStorageExecutionPhase.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Put", mock.Anything)
	mockStorageEventCount.AssertNotCalled(t, "Append", mock.Anything)
	mockStorageEventTopics.AssertNotCalled(t, "Append", mock.Anything, mock.Anything)
}

func Test_Module_DepositEvent_Success(t *testing.T) {
	firstHash := [32]sc.U8{}
	firstHash[0] = 1
	secondHash := [32]sc.U8{}
	secondHash[0] = 2
	event := newEventCodeUpdated(moduleId)
	expectEventRecord := primitives.EventRecord{
		Phase:  primitives.NewExtrinsicPhaseInitialization(),
		Event:  event,
		Topics: []primitives.H256{},
	}
	blockNum := sc.U64(1)
	eventCount := sc.U32(2)
	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(blockNum)
	mockStorageExecutionPhase.On("Get").Return(primitives.NewExtrinsicPhaseInitialization())
	mockStorageEventCount.On("Get").Return(eventCount)
	mockStorageEventCount.On("Put", eventCount+1).Return()
	mockStorageEvents.On("Append", expectEventRecord).Return()

	target.DepositEvent(event)

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Put", eventCount+1)
	mockStorageEvents.AssertCalled(t, "Append", expectEventRecord)
}

func Test_Module_depositEventIndexed_Success(t *testing.T) {
	firstHash := [32]sc.U8{}
	firstHash[0] = 1
	secondHash := [32]sc.U8{}
	secondHash[0] = 2
	topics := []primitives.H256{
		{
			firstHash[:],
		},
		{
			secondHash[:],
		},
	}
	event := newEventCodeUpdated(moduleId)
	expectEventRecord := primitives.EventRecord{
		Phase:  primitives.NewExtrinsicPhaseInitialization(),
		Event:  event,
		Topics: topics,
	}
	blockNum := sc.U64(1)
	eventCount := sc.U32(2)
	topicValue := sc.NewVaryingData(blockNum, eventCount)
	target := setupModule()

	mockStorageBlockNumber.On("Get").Return(blockNum)
	mockStorageExecutionPhase.On("Get").Return(primitives.NewExtrinsicPhaseInitialization())
	mockStorageEventCount.On("Get").Return(eventCount)
	mockStorageEventCount.On("Put", eventCount+1).Return()
	mockStorageEvents.On("Append", expectEventRecord).Return()
	mockStorageEventTopics.On("Append", topics[0], topicValue).Once()
	mockStorageEventTopics.On("Append", topics[1], topicValue).Once()

	target.depositEventIndexed(topics, event)

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Put", eventCount+1)
	mockStorageEvents.AssertCalled(t, "Append", expectEventRecord)
	mockStorageEventTopics.AssertNumberOfCalls(t, "Append", 2)
	mockStorageEventTopics.AssertCalled(t, "Append", topics[0], topicValue)
	mockStorageEventTopics.AssertCalled(t, "Append", topics[0], topicValue)
}

func Test_Module_DepositEvent_Overflow(t *testing.T) {
	target := setupModule()
	mockStorageBlockNumber.On("Get").Return(sc.U64(1))
	mockStorageExecutionPhase.On("Get").Return(primitives.NewExtrinsicPhaseInitialization())
	mockStorageEventCount.On("Get").Return(sc.U32(math.MaxUint32))

	target.DepositEvent(newEventCodeUpdated(moduleId))

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertCalled(t, "Get")
	mockStorageEventCount.AssertCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Put", mock.Anything)
	mockStorageEventCount.AssertNotCalled(t, "Append", mock.Anything)
	mockStorageEventTopics.AssertNotCalled(t, "Append", mock.Anything, mock.Anything)
}

func Test_Module_DepositEvent_ZeroBlockNumber(t *testing.T) {
	target := setupModule()
	mockStorageBlockNumber.On("Get").Return(sc.U64(0))

	target.DepositEvent(newEventCodeUpdated(moduleId))

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Put", mock.Anything)
	mockStorageEventCount.AssertNotCalled(t, "Append", mock.Anything)
	mockStorageEventTopics.AssertNotCalled(t, "Append", mock.Anything, mock.Anything)
}

func Test_Module_decrementProviders_HasAccount_NoProvidersLeft(t *testing.T) {
	target := setupModule()
	maybeAccount := sc.NewOption[primitives.AccountInfo](primitives.AccountInfo{})
	expect := sc.Result[sc.Encodable]{
		Value: primitives.DecRefStatusReaped,
	}

	mockStorageBlockNumber.On("Get").Return(sc.U64(0))

	result := target.decrementProviders(targetAccount, &maybeAccount)

	assert.Equal(t, expect, result)
	assert.Equal(t, sc.U32(1), maybeAccount.Value.Providers)

	mockStorageBlockNumber.AssertCalled(t, "Get")
	mockStorageExecutionPhase.AssertNotCalled(t, "Get")
	mockStorageEventCount.AssertNotCalled(t, "Get")
}

func Test_Module_decrementProviders_HasAccount_ConsumerRemaining(t *testing.T) {
	target := setupModule()
	accountInfo := primitives.AccountInfo{
		Consumers: 1,
		Data:      primitives.AccountData{},
	}
	maybeAccount := sc.NewOption[primitives.AccountInfo](accountInfo)
	expect := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorConsumerRemaining(),
	}

	result := target.decrementProviders(targetAccount, &maybeAccount)

	assert.Equal(t, expect, result)
	assert.Equal(t, sc.U32(1), maybeAccount.Value.Providers)
	assert.Equal(t, sc.U32(1), maybeAccount.Value.Consumers)
}

func Test_Module_decrementProviders_HasAccount_ContinueExist(t *testing.T) {
	target := setupModule()
	accountInfo := primitives.AccountInfo{
		Sufficients: 1,
		Data:        primitives.AccountData{},
	}
	maybeAccount := sc.NewOption[primitives.AccountInfo](accountInfo)
	expect := sc.Result[sc.Encodable]{
		Value: primitives.DecRefStatusExists,
	}

	result := target.decrementProviders(targetAccount, &maybeAccount)

	assert.Equal(t, expect, result)
	assert.Equal(t, sc.U32(0), maybeAccount.Value.Providers)
	assert.Equal(t, sc.U32(1), maybeAccount.Value.Sufficients)
}

func Test_Module_decrementProviders_NoAccount(t *testing.T) {
	target := setupModule()
	maybeAccount := sc.NewOption[primitives.AccountInfo](nil)
	expect := sc.Result[sc.Encodable]{
		Value: primitives.DecRefStatusReaped,
	}

	result := target.decrementProviders(targetAccount, &maybeAccount)

	assert.Equal(t, expect, result)
}

func Test_Module_mutateAccount(t *testing.T) {
	accountInfo := &primitives.AccountInfo{
		Nonce:       1,
		Consumers:   2,
		Providers:   3,
		Sufficients: 4,
		Data: primitives.AccountData{
			Free:       sc.NewU128(1),
			Reserved:   sc.NewU128(2),
			MiscFrozen: sc.NewU128(3),
			FeeFrozen:  sc.NewU128(4),
		},
	}
	accountData := primitives.AccountData{
		Free:       sc.NewU128(5),
		Reserved:   sc.NewU128(6),
		MiscFrozen: sc.NewU128(7),
		FeeFrozen:  sc.NewU128(8),
	}
	expectAccountInfo := &primitives.AccountInfo{
		Nonce:       1,
		Consumers:   2,
		Providers:   3,
		Sufficients: 4,
		Data:        accountData,
	}

	mutateAccount(accountInfo, &accountData)

	assert.Equal(t, expectAccountInfo, accountInfo)
}

func Test_Module_mutateAccount_NilData(t *testing.T) {
	accountInfo := &primitives.AccountInfo{
		Nonce:       1,
		Consumers:   2,
		Providers:   3,
		Sufficients: 4,
		Data: primitives.AccountData{
			Free:       sc.NewU128(1),
			Reserved:   sc.NewU128(2),
			MiscFrozen: sc.NewU128(3),
			FeeFrozen:  sc.NewU128(4),
		},
	}
	expectAccountInfo := &primitives.AccountInfo{
		Nonce:       1,
		Consumers:   2,
		Providers:   3,
		Sufficients: 4,
		Data:        primitives.AccountData{},
	}

	mutateAccount(accountInfo, nil)

	assert.Equal(t, expectAccountInfo, accountInfo)
}

func Test_Module_Metadata(t *testing.T) {
	target := setupModule()

	expectMetadataTypes := sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesPhase,
			"frame_system Phase",
			sc.Sequence[sc.Str]{"frame_system", "Phase"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"ApplyExtrinsic",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),
						},
						primitives.PhaseApplyExtrinsic,
						"Phase.ApplyExtrinsic"),
					primitives.NewMetadataDefinitionVariant(
						"Finalization",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.PhaseFinalization,
						"Phase.Finalization"),
					primitives.NewMetadataDefinitionVariant(
						"Initialization",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.PhaseInitialization,
						"Phase.Initialization"),
				})),
		primitives.NewMetadataType(metadata.TypesSystemEventStorage,
			"Vec<Box<EventRecord<T::RuntimeEvent, T::Hash>>>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesEventRecord))),

		primitives.NewMetadataType(metadata.TypesVecBlockNumEventIndex, "Vec<BlockNumber, EventIndex>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleU32U32))),

		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeight, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "mandatory", "T"),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight, "T"),
		),
		primitives.NewMetadataTypeWithPath(metadata.TypesWeightPerClass, "WeightPerClass", sc.Sequence[sc.Str]{"frame_system", "limits", "WeightsPerClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_extrinsic", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_extrinsic", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_total", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "reserved", "Option<Weight>"),
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeightsPerClass, "PerDispatchClass<WeightPerClass>", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "mandatory", "T"),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesWeightPerClass, "T")),

		primitives.NewMetadataTypeWithPath(metadata.TypesBlockWeights,
			"BlockWeights",
			sc.Sequence[sc.Str]{"frame_system", "limits", "BlockWeights"}, primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_block", "Weight"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "max_block", "Weight"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassWeightsPerClass, "per_class", "PerDispatchClass<WeightPerClass>"),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesDbWeight, "sp_weights RuntimeDbWeight", sc.Sequence[sc.Str]{"sp_weights", "RuntimeDbWeight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // read
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // write
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesBlockLength,
			"frame_system limits BlockLength",
			sc.Sequence[sc.Str]{"frame_system", "limits", "BlockLength"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassU32, "max", "PerDispatchClass<u32>"), // max
			})),

		primitives.NewMetadataTypeWithParams(metadata.TypesEventRecord,
			"frame_system EventRecord",
			sc.Sequence[sc.Str]{"frame_system", "EventRecord"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPhase, "phase", "Phase"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesRuntimeEvent, "event", "E"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesVecTopics, "topics", "Vec<T>"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesRuntimeEvent, "E"),
				primitives.NewMetadataTypeParameter(metadata.TypesH256, "T"),
			}),
		primitives.NewMetadataTypeWithPath(metadata.TypesSystemEvent,
			"frame_system pallet Event",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"ExtrinsicSuccess",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
						},
						EventExtrinsicSuccess,
						"Event.ExtrinsicSuccess"),
					primitives.NewMetadataDefinitionVariant(
						"ExtrinsicFailed",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchError, "dispatch_error", "DispatchError"),
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
						},
						EventExtrinsicFailed,
						"Events.ExtrinsicFailed"),
					primitives.NewMetadataDefinitionVariant(
						"CodeUpdated",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						EventCodeUpdated,
						"Events.CodeUpdated"),
					primitives.NewMetadataDefinitionVariant(
						"NewAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						},
						EventNewAccount,
						"Events.NewAccount"),
					primitives.NewMetadataDefinitionVariant(
						"KilledAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						},
						EventKilledAccount,
						"Events.KilledAccount"),
					primitives.NewMetadataDefinitionVariant(
						"Remarked",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "sender", "T::AccountId"),
							primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesH256, "hash", "T::Hash"),
						},
						EventRemarked,
						"Events.Remarked"),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesLastRuntimeUpgradeInfo,
			"LastRuntimeUpgradeInfo",
			sc.Sequence[sc.Str]{"frame_system", "LastRuntimeUpgradeInfo"}, primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU32),
					primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesSystemErrors,
			"frame_system pallet Error",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Error"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"InvalidSpecName",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorInvalidSpecName,
						"The name of specification does not match between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"SpecVersionNeedsToIncrease",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorSpecVersionNeedsToIncrease,
						"The specification version is not allowed to decrease between the current runtime and the new runtime."),
					primitives.NewMetadataDefinitionVariant(
						"FailedToExtractRuntimeVersion",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorFailedToExtractRuntimeVersion,
						"Failed to extract the runtime version from the new runtime.  Either calling `Core_version` or decoding `RuntimeVersion` failed."),
					primitives.NewMetadataDefinitionVariant(
						"NonDefaultComposite",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorNonDefaultComposite,
						"Suicide called when the account has non-default composite data."),
					primitives.NewMetadataDefinitionVariant(
						"NonZeroRefCount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorNonZeroRefCount,
						"There is a non-zero reference count preventing the account from being purged."),
					primitives.NewMetadataDefinitionVariant(
						"CallFiltered",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ErrorCallFiltered,
						"The origin filter prevent the call to be dispatched."),
				})),

		primitives.NewMetadataTypeWithParam(metadata.SystemCalls,
			"System calls",
			sc.Sequence[sc.Str]{"frame_system", "pallet", "Call"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"remark",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
						},
						functionRemarkIndex,
						"Make some on-chain remark."),
				}),
			primitives.NewMetadataEmptyTypeParameter("T")),

		primitives.NewMetadataTypeWithPath(metadata.TypesEra, "Era", sc.Sequence[sc.Str]{"sp_runtime", "generic", "era", "Era"}, primitives.NewMetadataTypeDefinitionVariant(primitives.EraTypeDefinition())),

		primitives.NewMetadataTypeWithParams(metadata.TypesBlock, "Block",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "block", "Block"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.Header, "Header"),
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceUncheckedExtrinsics, "Vec<Extrinsic>"),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.Header, "Header"),
				primitives.NewMetadataTypeParameter(metadata.UncheckedExtrinsic, "Extrinsic"),
			},
		),

		primitives.NewMetadataTypeWithPath(metadata.TypesTransactionSource, "TransactionSource", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "TransactionSource"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"InBlock",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.TransactionSourceInBlock,
						"TransactionSourceInBlock"),
					primitives.NewMetadataDefinitionVariant(
						"Local",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.TransactionSourceLocal,
						"TransactionSourceLocal"),
					primitives.NewMetadataDefinitionVariant(
						"External",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.TransactionSourceExternal,
						"TransactionSourceExternal"),
				})),

		primitives.NewMetadataTypeWithPath(metadata.TypesValidTransaction, "ValidTransaction", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "ValidTransaction"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "TransactionPriority"),
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceSequenceU8, "Vec<TransactionTag>"),
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceSequenceU8, "Vec<TransactionTag>"),
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "TransactionLongevity"),
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "bool"),
				},
			)),

		// type 871
		primitives.NewMetadataTypeWithPath(metadata.TypesInvalidTransaction, "InvalidTransaction", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "InvalidTransaction"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Call",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						CallIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Payment",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						PaymentIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Future",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						FutureIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Stale",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						StaleIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"BadProof",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						BadProofIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"AncientBirthBlock",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						AncientBirthBlockIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"ExhaustsResources",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						ExhaustsResourcesIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Custom",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),
						},
						CustomIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"BadMandatory",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						BadMandatoryIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"MandatoryValidation",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						MandatoryValidationIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"BadSigner",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						BadSignerIndex,
						""),
				},
			)),

		// type 872
		primitives.NewMetadataTypeWithPath(metadata.TypesUnknownTransaction, "UnknownTransaction", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "UnknownTransaction"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"CannotLookup",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						CannotLookupIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"NoUnsignedValidator",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						NoUnsignedValidatorIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Custom",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),
						},
						CustomUnknownIndex,
						""),
				},
			)),

		// type 870
		primitives.NewMetadataTypeWithPath(metadata.TypesTransactionValidityError, "TransactionValidityError", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "TransactionValidityError"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Invalid",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesInvalidTransaction),
						},
						InvalidTransactionIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Unknown",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesUnknownTransaction),
						},
						UnknownTransactionIndex,
						""),
				},
			)),

		primitives.NewMetadataTypeWithPath(metadata.TypesResultValidityTransaction, "Result", sc.Sequence[sc.Str]{"Result"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Ok",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesValidTransaction),
						},
						ValidTransactionIndex,
						""),
					primitives.NewMetadataDefinitionVariant(
						"Err",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesTransactionValidityError),
						},
						TransactionErrIndex,
						""),
				})),
	}

	moduleV14 := primitives.MetadataModuleV14{
		Name: name,
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: name,
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"Account",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
						sc.ToCompact(metadata.TypesAddress32),
						sc.ToCompact(metadata.TypesAccountInfo)),
					"The full account information for a particular account ID."),
				primitives.NewMetadataModuleStorageEntry(
					"ExtrinsicCount",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"Total extrinsics count for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"BlockWeight",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesPerDispatchClassWeight)),
					"The current weight for the block."),
				primitives.NewMetadataModuleStorageEntry(
					"AllExtrinsicsLen",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"Total length (in bytes) for all extrinsics put together, for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"BlockHash",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
						sc.ToCompact(metadata.PrimitiveTypesU32),
						sc.ToCompact(metadata.TypesFixedSequence32U8)),
					"Map of block numbers to block hashes."),
				primitives.NewMetadataModuleStorageEntry(
					"ExtrinsicData",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
						sc.ToCompact(metadata.PrimitiveTypesU32),
						sc.ToCompact(metadata.TypesSequenceU8)),
					"Extrinsics data for the current block (maps an extrinsic's index to its data)."),
				primitives.NewMetadataModuleStorageEntry(
					"Number",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"The current block number being processed. Set by `execute_block`."),
				primitives.NewMetadataModuleStorageEntry(
					"ParentHash",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesFixedSequence32U8)),
					"Hash of the previous block."),
				primitives.NewMetadataModuleStorageEntry(
					"Digest",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesDigest)),
					"Digest of the current block, also part of the block header."),
				primitives.NewMetadataModuleStorageEntry(
					"Events",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesSystemEventStorage)),
					"Events deposited for the current block.   NOTE: The item is unbound and should therefore never be read on chain."),
				primitives.NewMetadataModuleStorageEntry(
					"EventTopics",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
						sc.ToCompact(metadata.TypesH256),
						sc.ToCompact(metadata.TypesVecBlockNumEventIndex)), "Mapping between a topic (represented by T::Hash) and a vector of indexes  of events in the `<Events<T>>` list."),
				primitives.NewMetadataModuleStorageEntry(
					"EventCount",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"The number of events in the `Events<T>` list."),
				primitives.NewMetadataModuleStorageEntry(
					"LastRuntimeUpgrade",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesLastRuntimeUpgradeInfo)),
					"Stores the `spec_version` and `spec_name` of when the last runtime upgrade happened."),
				primitives.NewMetadataModuleStorageEntry(
					"ExecutionPhase",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesPhase)),
					"The execution phase of the block."),
			},
		}),
		Call: sc.NewOption[sc.Compact](sc.ToCompact(metadata.SystemCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				name,
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SystemCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<System, Runtime>"),
				},
				moduleId,
				"Call.System"),
		),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemEvent)),
		EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				name,
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSystemEvent, "frame_system::Event<Runtime>"),
				},
				moduleId,
				"Events.System"),
		),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"BlockWeights",
				sc.ToCompact(metadata.TypesBlockWeights),
				sc.BytesToSequenceU8(blockWeights.Bytes()),
				"Block & extrinsics weights: base values and limits.",
			),
			primitives.NewMetadataModuleConstant(
				"BlockLength",
				sc.ToCompact(metadata.TypesBlockLength),
				sc.BytesToSequenceU8(blockLength.Bytes()),
				"The maximum length of a block (in bytes).",
			),
			primitives.NewMetadataModuleConstant(
				"BlockHashCount",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(blockHashCount.Bytes()),
				"Maximum number of block number to block hash mappings to keep (oldest pruned first).",
			),
			primitives.NewMetadataModuleConstant(
				"DbWeight",
				sc.ToCompact(metadata.TypesDbWeight),
				sc.BytesToSequenceU8(dbWeight.Bytes()),
				"The weight of runtime database operations the runtime can invoke.",
			),
			primitives.NewMetadataModuleConstant(
				"Version",
				sc.ToCompact(metadata.TypesRuntimeVersion),
				sc.BytesToSequenceU8(version.Bytes()),
				"Get the chain's current version.",
			),
		},
		Error: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesSystemErrors)),
		Index: moduleId,
	}

	expectMetadataModule := primitives.MetadataModule{
		Version:   primitives.ModuleVersion14,
		ModuleV14: moduleV14,
	}

	resultTypes, resultMetadataModule := target.Metadata()

	assert.Equal(t, expectMetadataTypes, resultTypes)
	assert.Equal(t, expectMetadataModule, resultMetadataModule)
}

func testDigest() primitives.Digest {
	digest := primitives.Digest{}

	digest[primitives.DigestTypeSeal] = append(digest[primitives.DigestTypeSeal], primitives.DigestItem{
		Engine:  sc.NewFixedSequence[sc.U8](2, sc.U8(0), sc.U8(1)),
		Payload: sc.BytesToSequenceU8(sc.U64(5).Bytes()),
	})

	return digest
}

func setupModule() module {
	config := NewConfig(blockHashCount, blockWeights, blockLength, dbWeight, version)

	target := New(moduleId, config).(module)

	initMockStorage()
	target.storage.Account = mockStorageAccount
	target.storage.BlockWeight = mockStorageBlockWeight
	target.storage.BlockHash = mockStorageBlockHash
	target.storage.BlockNumber = mockStorageBlockNumber
	target.storage.AllExtrinsicsLen = mockStorageAllExtrinsicsLen
	target.storage.ExtrinsicIndex = mockStorageExtrinsicIndex
	target.storage.ExtrinsicData = mockStorageExtrinsicData
	target.storage.ExtrinsicCount = mockStorageExtrinsicCount
	target.storage.ParentHash = mockStorageParentHash
	target.storage.Digest = mockStorageDigest
	target.storage.Events = mockStorageEvents
	target.storage.EventCount = mockStorageEventCount
	target.storage.EventTopics = mockStorageEventTopics
	target.storage.LastRuntimeUpgrade = mockStorageLastRuntimeUpgrade
	target.storage.ExecutionPhase = mockStorageExecutionPhase

	target.ioStorage = mockIoStorage
	target.trie = mockIoTrie

	return target
}

func initMockStorage() {
	mockStorageAccount = new(mocks.StorageMap[primitives.PublicKey, primitives.AccountInfo])
	mockStorageBlockWeight = new(mocks.StorageValue[primitives.ConsumedWeight])
	mockStorageBlockHash = new(mocks.StorageMap[sc.U64, primitives.Blake2bHash])
	mockStorageBlockNumber = new(mocks.StorageValue[sc.U64])
	mockStorageAllExtrinsicsLen = new(mocks.StorageValue[sc.U32])
	mockStorageExtrinsicIndex = new(mocks.StorageValue[sc.U32])
	mockStorageExtrinsicData = new(mocks.StorageMap[sc.U32, sc.Sequence[sc.U8]])
	mockStorageExtrinsicCount = new(mocks.StorageValue[sc.U32])
	mockStorageParentHash = new(mocks.StorageValue[primitives.Blake2bHash])
	mockStorageDigest = new(mocks.StorageValue[primitives.Digest])
	mockStorageEvents = new(mocks.StorageValue[primitives.EventRecord])
	mockStorageEventCount = new(mocks.StorageValue[sc.U32])
	mockStorageEventTopics = new(mocks.StorageMap[primitives.H256, sc.VaryingData])
	mockStorageLastRuntimeUpgrade = new(mocks.StorageValue[primitives.LastRuntimeUpgradeInfo])
	mockStorageExecutionPhase = new(mocks.StorageValue[primitives.ExtrinsicPhase])

	mockIoStorage = new(mocks.IoStorage)
	mockIoTrie = new(mocks.IoTrie)
}
