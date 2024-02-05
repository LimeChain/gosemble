package core

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	errPanic = errors.New("panic")
)

var (
	mockExecutive      *mocks.Executive
	mockRuntimeDecoder *mocks.RuntimeDecoder
	mockMemoryUtils    *mocks.MemoryTranslator
	mdGenerator        = primitives.NewMetadataTypeGenerator()
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName, apiVersion)

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
		Digest:         primitives.NewDigest(sc.Sequence[primitives.DigestItem]{}),
	}

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(header.Bytes())
	mockExecutive.On("InitializeBlock", header).Return(nil)

	target.InitializeBlock(dataPtr, dataLen)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockExecutive.AssertCalled(t, "InitializeBlock", header)
}

func Test_Module_InitializeBlock_DecodeHeader_Panics(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return([]byte{})

	assert.PanicsWithValue(t,
		io.EOF.Error(),
		func() { target.InitializeBlock(dataPtr, dataLen) },
	)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockExecutive.AssertNotCalled(t, "InitializeBlock", mock.Anything)

}

func Test_Module_InitializeBlock_InitializeBlock_Panics(t *testing.T) {
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
		Digest:         primitives.NewDigest(sc.Sequence[primitives.DigestItem]{}),
	}

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(header.Bytes())
	mockExecutive.On("InitializeBlock", header).Return(errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.InitializeBlock(dataPtr, dataLen) },
	)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockExecutive.AssertCalled(t, "InitializeBlock", header)
}

func Test_Module_ExecuteBlock(t *testing.T) {
	target := setup()
	bBlock := []byte{1, 2, 3}
	block := types.NewBlock(primitives.Header{Number: 2}, sc.Sequence[primitives.UncheckedExtrinsic]{})
	buffer := bytes.NewBuffer(bBlock)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bBlock)
	mockRuntimeDecoder.On("DecodeBlock", buffer).Return(block, nil)
	mockExecutive.On("ExecuteBlock", block).Return(nil)

	target.ExecuteBlock(dataPtr, dataLen)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ExecuteBlock", block)
}

func Test_Module_ExecuteBlock_DecodeBlock_Panics(t *testing.T) {
	target := setup()
	bBlock := []byte{1, 2, 3}
	block := types.NewBlock(primitives.Header{Number: 2}, sc.Sequence[primitives.UncheckedExtrinsic]{})
	buffer := bytes.NewBuffer(bBlock)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bBlock)
	mockRuntimeDecoder.On("DecodeBlock", buffer).Return(block, errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.ExecuteBlock(dataPtr, dataLen) },
	)
}

func Test_Module_ExecuteBlock_ExecuteBlock_Panics(t *testing.T) {
	target := setup()
	bBlock := []byte{1, 2, 3}
	block := types.NewBlock(primitives.Header{Number: 2}, sc.Sequence[primitives.UncheckedExtrinsic]{})
	buffer := bytes.NewBuffer(bBlock)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(bBlock)
	mockRuntimeDecoder.On("DecodeBlock", buffer).Return(block, nil)
	mockExecutive.On("ExecuteBlock", block).Return(errPanic)

	assert.PanicsWithValue(t,
		errPanic.Error(),
		func() { target.ExecuteBlock(dataPtr, dataLen) },
	)

	mockMemoryUtils.AssertCalled(t, "GetWasmMemorySlice", dataPtr, dataLen)
	mockRuntimeDecoder.AssertExpectations(t)
	mockExecutive.AssertCalled(t, "ExecuteBlock", block)
}

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	blockId, _ := target.mdGenerator.GetId("block")

	expect := primitives.RuntimeApiMetadata{
		Name: ApiModuleName,
		Methods: sc.Sequence[primitives.RuntimeApiMethodMetadata]{
			primitives.RuntimeApiMethodMetadata{
				Name:   "version",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
				Output: sc.ToCompact(metadata.TypesRuntimeVersion),
				Docs:   sc.Sequence[sc.Str]{" Returns the version of the runtime."},
			},
			primitives.RuntimeApiMethodMetadata{
				Name: "execute_block",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "block",
						Type: sc.ToCompact(blockId),
					},
				},
				Output: sc.ToCompact(metadata.TypesEmptyTuple),
				Docs:   sc.Sequence[sc.Str]{" Execute the given block."},
			},
			primitives.RuntimeApiMethodMetadata{
				Name: "initialize_block",
				Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
					primitives.RuntimeApiMethodParamMetadata{
						Name: "header",
						Type: sc.ToCompact(metadata.Header),
					},
				},
				Output: sc.ToCompact(metadata.TypesEmptyTuple),
				Docs:   sc.Sequence[sc.Str]{" Initialize a block with the given header."},
			},
		},
		Docs: sc.Sequence[sc.Str]{" The `Core` runtime api that every Substrate runtime needs to implement."},
	}

	assert.Equal(t, expect, target.Metadata())
}

func setup() Module {
	mockExecutive = new(mocks.Executive)
	mockRuntimeDecoder = new(mocks.RuntimeDecoder)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New(mockExecutive, mockRuntimeDecoder, version, mdGenerator, log.NewLogger())
	target.memUtils = mockMemoryUtils

	return target
}
