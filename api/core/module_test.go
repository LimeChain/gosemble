package core

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(2)

	version = &primitives.RuntimeVersion{
		SpecName:           "test-spec-name",
		ImplName:           "test-impl-name",
		AuthoringVersion:   1,
		SpecVersion:        2,
		ImplVersion:        3,
		TransactionVersion: 4,
		StateVersion:       5,
	}
)

var (
	mockExecutive      *mocks.Executive
	mockRuntimeDecoder *mocks.RuntimeDecoder
	mockMemoryUtils    *mocks.MemoryTranslator
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

func Test_Module_Version(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("BytesToOffsetAndSize", version.Bytes()).Return(ptrAndSize)

	result := target.Version()

	assert.Equal(t, ptrAndSize, result)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", version.Bytes())
}

func Test_Module_InitializeBlock(t *testing.T) {
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
	mockExecutive.On("InitializeBlock", header).Return()

	target.InitializeBlock(dataPtr, dataLen)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockExecutive.AssertCalled(t, "InitializeBlock", header)
}

func Test_Module_ExecuteBlock(t *testing.T) {
	target := setup()
	bBlock := []byte{1, 2, 3}
	block := types.Block{
		Header: primitives.Header{
			Number: 2,
		}}
	buffer := bytes.NewBuffer(bBlock)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bBlock)
	mockRuntimeDecoder.On("DecodeBlock", buffer).Return(block)
	mockExecutive.On("ExecuteBlock", block).Return()

	target.ExecuteBlock(dataPtr, dataLen)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertCalled(t, "DecodeBlock", buffer)
	mockExecutive.AssertCalled(t, "ExecuteBlock", block)
}

func setup() Module {
	mockExecutive = new(mocks.Executive)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New(mockExecutive, mockRuntimeDecoder, version)
	target.memUtils = mockMemoryUtils

	return target
}
