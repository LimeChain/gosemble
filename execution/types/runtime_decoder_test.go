package types

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockModuleOne               *mocks.Module
	mockModuleTwo               *mocks.Module
	bytesExtrinsicFormatVersion = sc.U8(ExtrinsicFormatVersion).Bytes()
	mockCallOne                 *mocks.Call
	mockCallTwo                 *mocks.Call
	moduleOneIdx                = sc.U8(0)
	moduleTwoIdx                = sc.U8(1)
	mSignedExtra                *mocks.SignedExtra
	mockExtrinsicSignature      *mocks.ExtrinsicSignature
)

func RuntimeDecoder_New(t *testing.T) {
	target := setupRuntimeDecoder()
	expect := runtimeDecoder{
		modules: []primitives.Module{mockModuleOne, mockModuleTwo},
		extra:   mSignedExtra,
	}

	assert.Equal(t, expect, target)
}

func RuntimeDecoder_DecodeBlock_EmptyBody(t *testing.T) {
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

	//blockHeader := primitives.Header{
	//	ParentHash: testHash,
	//	1,
	//	StateRoot: testHash,
	//	ExtrinsicsRoot: testHash,
	//	Digest: types.NewDigest()
	//}

	lenExtrinsics := sc.ToCompact(0).Bytes()
	buff := bytes.NewBuffer(append(header.Bytes(), lenExtrinsics...))

	//ext := NewUncheckedExtrinsic(version, extSignature, function, rd.extra)

	//var exts [][]byte
	//err = scale.Unmarshal(inherentExt, &exts)
	//assert.Nil(t, err)

	//e1 := types.NewExtrinsic([]byte{0x01, 0x02})
	//e2 := types.NewExtrinsic([]byte{0x01, 0x03})
	//
	//// types.BytesArrayToExtrinsics()
	//
	//extrinsics := make([]types.Extrinsic, 2)
	//extrinsics = append(extrinsics, e1)
	//extrinsics = append(extrinsics, e2)
	//
	//blockBody := types.NewBody(extrinsics)
	//
	//block := types.Block{
	//	Header: *blockHeader,
	//	Body:   *blockBody,
	//}
	//
	//bytesBlock, err := block.Encode()
	//assert.NoError(t, err)
	//
	//buffer := bytes.NewBuffer(bytesBlock)

	target := setupRuntimeDecoder()

	resultBlock := target.DecodeBlock(buff)

	expectedBlock := NewBlock(header, sc.Sequence[primitives.UncheckedExtrinsic]{})

	assert.Equal(t, resultBlock, expectedBlock)

	//runtimeDecoder.AssertCalled(t, "Metadata")

}

func RuntimeDecoder_DecodeBlock(t *testing.T) {
	target := setupRuntimeDecoder()

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

	ext1 := NewUnsignedUncheckedExtrinsic(mockCallOne)
	//ext2 := NewUnsignedUncheckedExtrinsic(mockCallTwo)

	extrinsics := sc.Sequence[primitives.UncheckedExtrinsic]{ext1}

	lenExtrinsics := sc.ToCompact(1).Bytes()
	bytesBuff := append(header.Bytes(), lenExtrinsics...)

	mockCallOne.On("Encode", mock.Anything)
	//mockCallTwo.On("Encode", mock.Anything)

	bytesBuff = append(bytesBuff, extrinsics.Bytes()...)

	buff := bytes.NewBuffer(bytesBuff)

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	//mockModuleTwo.On("GetIndex").Return(moduleTwoIdx)

	resultBlock := target.DecodeBlock(buff)

	expectedBlock := NewBlock(header, extrinsics)

	assert.Equal(t, resultBlock, expectedBlock)

	//runtimeDecoder.AssertCalled(t, "Metadata")

}

func RuntimeDecoder_DecodeUncheckedExtrinsic(t *testing.T) {
	target := setupRuntimeDecoder()

	extr := NewUnsignedUncheckedExtrinsic(mockCallOne)

	mockCallOne.On("Encode", mock.Anything)

	mockCallOne.On("Encode", mock.Anything)

	extrBytes := extr.Bytes()

	moduleIdx := sc.U8(0)
	functionIdx := sc.U8(0)
	//
	extrBytes = append(extrBytes, moduleIdx.Bytes()...)
	extrBytes = append(extrBytes, functionIdx.Bytes()...)
	//arg := sc.NewU128(5)
	//extrBytes = append(extrBytes, arg.Bytes()...)

	buff := bytes.NewBuffer(extrBytes)

	//mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	//mockModuleTwo.On("GetIndex").Return(moduleTwoIdx)

	functions := map[sc.U8]primitives.Call{}
	functions[0] = mockCallOne

	mockModuleOne.On("GetIndex").Return(moduleIdx)
	mockModuleOne.On("Functions").Return(functions)
	mockCallOne.On("DecodeArgs", buff).Return(mockCallOne)

	result := target.DecodeUncheckedExtrinsic(buff)

	assert.Equal(t, result, extr)
}

func Test_RuntimeDecoder_DecodeUncheckedExtrinsic(t *testing.T) {
	target := setupRuntimeDecoder()

	extrinsicSignature = sc.NewOption[primitives.ExtrinsicSignature](
		primitives.ExtrinsicSignature{
			Signer:    signer,
			Signature: signatureEd25519,
			Extra:     mSignedExtra,
		},
	)

	uxt := NewUncheckedExtrinsic(ExtrinsicFormatVersion, extrinsicSignature, mockCallOne, mSignedExtra).(uncheckedExtrinsic)

	mockCallOne.On("Encode", mock.Anything).Return()
	mockExtrinsicSignature.On("Encode", mock.Anything)
	mSignedExtra.On("Encode", mock.Anything)

	buff := bytes.NewBuffer(uxt.Bytes())

	functions := map[sc.U8]primitives.Call{}
	functions[0] = mockCallOne

	mockModuleOne.On("GetIndex").Return(moduleOneIdx)
	mockModuleOne.On("Functions").Return(functions)
	mockCallOne.On("DecodeArgs", buff).Return(mockCallOne)
	result := target.DecodeUncheckedExtrinsic(buff)

	assert.Equal(t, result, uxt)
}

func setupRuntimeDecoder() RuntimeDecoder {
	mockModuleOne = new(mocks.Module)

	mockCallOne = new(mocks.Call)

	mSignedExtra = new(mocks.SignedExtra)

	mockExtrinsicSignature = new(mocks.ExtrinsicSignature)

	apis := []primitives.Module{mockModuleOne}

	return NewRuntimeDecoder(apis, mSignedExtra)
}
