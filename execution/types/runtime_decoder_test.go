package types

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockModuleOne        *mocks.Module
	mockCallOne          *mocks.Call
	moduleOneIdx         = sc.U8(0)
	moduleIdx            = sc.U8(0)
	functionIdx          = sc.U8(0)
	signedExtrinsicBytes = []byte{
		132,                                                                                               // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		uint8(moduleOneIdx), uint8(functionIdx), // call
	}

	expectedSignedExtrinsicsBytesAfterDecode = []byte{
		0x8d, 0x1, // length
		132,                                                                                               // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		// call
	}

	header = primitives.Header{
		ParentHash: primitives.Blake2bHash{
			FixedSequence: sc.BytesToFixedSequenceU8(parentHash)},
		Number:         5,
		StateRoot:      primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(stateRoot)},
		ExtrinsicsRoot: primitives.H256{FixedSequence: sc.BytesToFixedSequenceU8(extrinsicsRoot)},
		Digest:         primitives.Digest{},
	}

	expectedUncheckedExtrinsicUnsignedAfterDecode = []byte{0x4, 0x4}

	parentHash      = common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355c").ToBytes()
	stateRoot       = common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355b").ToBytes()
	extrinsicsRoot  = common.MustHexToHash("0x3aa96b0149b6ca3688878bdbd19464448624136398e3ce45b9e755d3ab61355a").ToBytes()
	moduleFunctions = map[sc.U8]primitives.Call{}
)

func Test_RuntimeDecoder_New(t *testing.T) {
	target := setupRuntimeDecoder()
	expect := runtimeDecoder{
		modules: []primitives.Module{mockModuleOne},
		extra:   mockSignedExtra,
	}

	assert.Equal(t, expect, target)
}

func Test_RuntimeDecoder_DecodeBlock_ZeroExtrinsicsEmptyBody(t *testing.T) {
	target := setupRuntimeDecoder()

	lenExtrinsics := sc.ToCompact(0).Bytes()
	buff := bytes.NewBuffer(append(header.Bytes(), lenExtrinsics...))

	resultBlock := target.DecodeBlock(buff)

	expectedBlock := NewBlock(header, sc.Sequence[primitives.UncheckedExtrinsic]{})

	assert.Equal(t, resultBlock, expectedBlock)
}

func Test_RuntimeDecoder_DecodeBlock_Module_Not_Exists(t *testing.T) {
	target := setupRuntimeDecoder()

	moduleFunctions[0] = mockCallOne

	idxModuleNotExists := sc.U8(10)

	nonExistentExtrinsicsBytes := []byte{
		132,                                                                                               // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		uint8(idxModuleNotExists), uint8(functionIdx), // call
	}

	buffExtrinsic := bytes.NewBuffer(append(sc.ToCompact(len(nonExistentExtrinsicsBytes)).Bytes(), nonExistentExtrinsicsBytes...))

	lenExtrinsics := sc.ToCompact(1).Bytes()
	decodeBlockBytes := append(lenExtrinsics, buffExtrinsic.Bytes()...)
	decodeBlockBytes = append(header.Bytes(), decodeBlockBytes...)

	mockCallOne.On("Encode", mock.Anything)

	decodeBlockBuff := bytes.NewBuffer(decodeBlockBytes)

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockSignedExtra.On("Decode", mock.Anything).Return()

	assert.PanicsWithValue(t, "module with index ["+strconv.Itoa(int(idxModuleNotExists))+"] not found", func() {
		target.DecodeBlock(decodeBlockBuff)
	})
}

func Test_RuntimeDecoder_DecodeBlock_Function_Not_Exists(t *testing.T) {
	target := setupRuntimeDecoder()

	moduleFunctions[0] = mockCallOne

	idxFunctionNotExists := sc.U8(10)

	nonExistentExtrinsicsBytes := []byte{
		132,                                                                                               // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		uint8(moduleOneIdx), uint8(idxFunctionNotExists), // call
	}

	buffExtrinsic := bytes.NewBuffer(append(sc.ToCompact(len(nonExistentExtrinsicsBytes)).Bytes(), nonExistentExtrinsicsBytes...))

	lenExtrinsics := sc.ToCompact(1).Bytes()
	decodeBlockBytes := append(lenExtrinsics, buffExtrinsic.Bytes()...)
	decodeBlockBytes = append(header.Bytes(), decodeBlockBytes...)

	mockCallOne.On("Encode", mock.Anything)

	decodeBlockBuff := bytes.NewBuffer(decodeBlockBytes)

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockSignedExtra.On("Decode", mock.Anything).Return()
	mockModuleOne.On("Functions").Return(moduleFunctions)

	assert.PanicsWithValue(t, "function index ["+strconv.Itoa(int(idxFunctionNotExists))+"] for module ["+strconv.Itoa(int(moduleOneIdx))+"] not found", func() {
		target.DecodeBlock(decodeBlockBuff)
	})
}

func Test_RuntimeDecoder_DecodeBlock_Single_Extrinsic(t *testing.T) {
	target := setupRuntimeDecoder()

	moduleFunctions[0] = mockCallOne

	buffExtrinsic := bytes.NewBuffer(append(sc.ToCompact(len(signedExtrinsicBytes)).Bytes(), signedExtrinsicBytes...))

	lenExtrinsics := sc.ToCompact(1).Bytes()
	decodeBlockBytes := append(lenExtrinsics, buffExtrinsic.Bytes()...)
	decodeBlockBytes = append(header.Bytes(), decodeBlockBytes...)

	mockCallOne.On("Encode", mock.Anything)

	decodeBlockBuff := bytes.NewBuffer(decodeBlockBytes)

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockSignedExtra.On("Decode", mock.Anything).Return()
	mockModuleOne.On("Functions").Return(moduleFunctions)
	mockCallOne.On("DecodeArgs", decodeBlockBuff).Return(mockCallOne)
	result := target.DecodeBlock(decodeBlockBuff)

	mockCallOne.On("DecodeArgs", buffExtrinsic).Return(mockCallOne)

	decodedExtr := target.DecodeUncheckedExtrinsic(buffExtrinsic)

	extrinsics := sc.Sequence[primitives.UncheckedExtrinsic]{decodedExtr}

	expectedBlock := NewBlock(header, extrinsics)

	mockSignedExtra.On("Encode", mock.Anything)
	assert.Equal(t, result.Bytes(), expectedBlock.Bytes())

	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleOne.AssertCalled(t, "Functions")
	mockCallOne.AssertCalled(t, "DecodeArgs", decodeBlockBuff)
	mockSignedExtra.AssertCalled(t, "Decode", mock.Anything)
}

func Test_RuntimeDecoder_DecodeBlock_Multiple_Extrinsics(t *testing.T) {
	target := setupRuntimeDecoder()
	moduleFunctions[0] = mockCallOne

	totalExtrinsicsInBlock := 10
	lenSignedExtrinsicBytes := sc.ToCompact(len(signedExtrinsicBytes)).Bytes()

	signedExtrinsicBytes := append(lenSignedExtrinsicBytes, signedExtrinsicBytes...)
	allExtrinsicsBytes := signedExtrinsicBytes
	assert.Equal(t, allExtrinsicsBytes, signedExtrinsicBytes)
	for i := 0; i < totalExtrinsicsInBlock-1; i++ {
		allExtrinsicsBytes = append(signedExtrinsicBytes, allExtrinsicsBytes...)
	}

	lenExtrinsics := sc.ToCompact(totalExtrinsicsInBlock).Bytes()
	decodeBlockBytes := append(lenExtrinsics, allExtrinsicsBytes...)
	decodeBlockBytes = append(header.Bytes(), decodeBlockBytes...)

	mockCallOne.On("Encode", mock.Anything)

	buff := bytes.NewBuffer(decodeBlockBytes)

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockSignedExtra.On("Decode", mock.Anything).Return()
	mockModuleOne.On("Functions").Return(moduleFunctions)
	mockCallOne.On("DecodeArgs", buff).Return(mockCallOne)
	result := target.DecodeBlock(buff)

	bufferSignedExtrinsics := bytes.NewBuffer(signedExtrinsicBytes)

	mockCallOne.On("DecodeArgs", bufferSignedExtrinsics).Return(mockCallOne)

	decodedExtr := target.DecodeUncheckedExtrinsic(bufferSignedExtrinsics)

	extrinsics := sc.Sequence[primitives.UncheckedExtrinsic]{}

	for i := 0; i < totalExtrinsicsInBlock; i++ {
		extrinsics = append(extrinsics, decodedExtr)
	}

	expectedBlock := NewBlock(header, extrinsics)

	mockSignedExtra.On("Encode", mock.Anything)
	assert.Equal(t, result.Bytes(), expectedBlock.Bytes())

	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleOne.AssertCalled(t, "Functions")
	mockCallOne.AssertCalled(t, "DecodeArgs", bufferSignedExtrinsics)
	mockSignedExtra.AssertCalled(t, "Decode", mock.Anything)
}

func Test_RuntimeDecoder_DecodeUncheckedExtrinsic_Unsigned(t *testing.T) {
	target := setupRuntimeDecoder()
	moduleFunctions[0] = mockCallOne

	unsignedExtrBytes := append(moduleIdx.Bytes(), functionIdx.Bytes()...)
	unsignedExtrBytes = append(sc.U8(ExtrinsicFormatVersion).Bytes(), unsignedExtrBytes...)
	unsignedExtrBytes = append(sc.ToCompact(len(unsignedExtrBytes)).Bytes(), unsignedExtrBytes...)

	buff := bytes.NewBuffer(unsignedExtrBytes)

	mockModuleOne.On("GetIndex").Return(moduleIdx)
	mockModuleOne.On("Functions").Return(moduleFunctions)
	mockCallOne.On("DecodeArgs", buff).Return(mockCallOne)

	result := target.DecodeUncheckedExtrinsic(buff)
	assert.Equal(t, result.IsSigned(), false)

	mockCallOne.On("Encode", mock.Anything)
	assert.Equal(t, result.Bytes(), expectedUncheckedExtrinsicUnsignedAfterDecode)

	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleOne.AssertCalled(t, "Functions")
	mockCallOne.AssertCalled(t, "DecodeArgs", buff)
}

func Test_RuntimeDecoder_DecodeUncheckedExtrinsic_Signed(t *testing.T) {
	target := setupRuntimeDecoder()

	mockSignedExtra.On("Decode", mock.Anything).Return()

	// Append the length of the bytes to decode as compact
	buff := bytes.NewBuffer(append(sc.ToCompact(len(signedExtrinsicBytes)).Bytes(), signedExtrinsicBytes...))

	moduleFunctions := map[sc.U8]primitives.Call{}
	moduleFunctions[0] = mockCallOne

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockModuleOne.On("Functions").Return(moduleFunctions)
	mockCallOne.On("DecodeArgs", buff).Return(mockCallOne)

	result := target.DecodeUncheckedExtrinsic(buff)

	mockCallOne.On("Encode", mock.Anything)
	mockSignedExtra.On("Encode", mock.Anything)

	assert.Equal(t, result.Bytes(), expectedSignedExtrinsicsBytesAfterDecode)
	assert.Equal(t, result.IsSigned(), true)

	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleOne.AssertCalled(t, "Functions")
	mockCallOne.AssertCalled(t, "DecodeArgs", buff)
}

func Test_RuntimeDecoder_DecodeCall(t *testing.T) {
	target := setupRuntimeDecoder()

	args := sc.NewVaryingData(sc.U8(1), sc.U8(2), sc.U8(3))

	callBytes := []byte{
		uint8(moduleOneIdx), uint8(functionIdx),
	}

	callBytes = append(callBytes, args.Bytes()...)

	buf := bytes.NewBuffer(callBytes)

	moduleFunctions := map[sc.U8]primitives.Call{}
	moduleFunctions[0] = mockCallOne

	mockModuleOne.On("GetIndex").Return(moduleIdx)
	mockModuleOne.On("Functions").Return(moduleFunctions)
	mockCallOne.On("DecodeArgs", buf).Run(func(args mock.Arguments) {
		buf := args.Get(0).(*bytes.Buffer)
		// reading 3 bytes for the 3 arguments
		buf.ReadByte()
		buf.ReadByte()
		buf.ReadByte()
	}).Return(mockCallOne)

	target.DecodeCall(buf)

	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleOne.AssertCalled(t, "Functions")
	mockCallOne.AssertCalled(t, "DecodeArgs", buf)
}

func setupRuntimeDecoder() RuntimeDecoder {
	mockModuleOne = new(mocks.Module)

	mockCallOne = new(mocks.Call)

	mockSignedExtra = new(mocks.SignedExtra)

	apis := []primitives.Module{mockModuleOne}

	return NewRuntimeDecoder(apis, mockSignedExtra)
}
