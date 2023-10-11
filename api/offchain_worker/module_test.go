package offchain_worker

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr = int32(0)
	dataLen = int32(5)
)

var (
	mockExecutive   *mocks.Executive
	mockMemoryUtils *mocks.MemoryTranslator
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName[:], apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_OffchainWorker(t *testing.T) {
	target := setup()

	parentHash := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355c").ToBytes()
	stateRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b").ToBytes()
	extrinsicsRoot := common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a").ToBytes()
	header := primitives.Header{
		ParentHash: primitives.Blake2bHash{
			FixedSequence: sc.BytesToFixedSequenceU8(parentHash)},
		Number:         5,
		StateRoot:      primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(stateRoot)},
		ExtrinsicsRoot: primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(extrinsicsRoot)},
		Digest:         primitives.Digest{},
	}

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(header.Bytes())
	mockExecutive.On("OffchainWorker", header).Return()

	target.OffchainWorker(dataPtr, dataLen)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockExecutive.AssertCalled(t, "OffchainWorker", header)
}

func setup() Module {
	mockExecutive = new(mocks.Executive)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New(mockExecutive)
	target.memUtils = mockMemoryUtils

	return target
}
