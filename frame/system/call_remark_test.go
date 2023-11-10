package system

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

const (
	moduleId = 0
)

var (
	baseWeight = primitives.WeightFromParts(567, 123)
)

func Test_Call_Remark_New(t *testing.T) {
	expected := callRemark{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionRemarkIndex,
		},
	}

	actual := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, expected, actual)
}

func Test_Call_Remark_DecodeArgs_Success(t *testing.T) {
	seq := sc.BytesToSequenceU8([]byte{1, 2, 3})
	buf := bytes.NewBuffer(seq.Bytes())

	call := newCallRemark(moduleId, functionRemarkIndex)
	call, err := call.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(seq), call.Args())
}

func Test_Call_Remark_Encode(t *testing.T) {
	expectedBuf := bytes.NewBuffer([]byte{moduleId, functionRemarkIndex})
	buf := &bytes.Buffer{}

	call := newCallRemark(moduleId, functionRemarkIndex)

	err := call.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuf, buf)
}

func Test_Call_Remark_EncodeWithArgs(t *testing.T) {
	seq := sc.BytesToSequenceU8([]byte{1, 2, 3})

	expectedBuf := bytes.NewBuffer(append([]byte{moduleId, functionRemarkIndex}, seq.Bytes()...))

	buf := bytes.NewBuffer(seq.Bytes())

	call := newCallRemark(moduleId, functionRemarkIndex)
	call, err := call.DecodeArgs(buf)
	assert.Nil(t, err)

	buf.Reset()
	err = call.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuf, buf)
}

func Test_Call_Remark_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionRemarkIndex}

	call := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, expected, call.Bytes())
}

func Test_Call_Remark_ModuleIndex(t *testing.T) {
	testCases := []sc.U8{
		moduleId,
		1,
		2,
		3,
	}

	for _, tc := range testCases {
		call := newCallRemark(tc, functionRemarkIndex)

		assert.Equal(t, tc, call.ModuleIndex())
	}
}

func Test_Call_Remark_FunctionIndex(t *testing.T) {
	testCases := []sc.U8{
		functionRemarkIndex,
		1,
		2,
		3,
	}

	for _, tc := range testCases {
		call := newCallRemark(moduleId, tc)

		assert.Equal(t, tc, call.FunctionIndex())
	}
}

func Test_Call_Remark_BaseWeight_EmptyArgs(t *testing.T) {
	call := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, primitives.WeightFromParts(2_091_000, 0), call.BaseWeight())
}

func Test_Call_Remark_BaseWeight_WithArgs(t *testing.T) {
	seq := sc.BytesToSequenceU8([]byte{1})
	call := newCallRemark(moduleId, functionRemarkIndex)
	call, err := call.DecodeArgs(bytes.NewBuffer(seq.Bytes()))
	assert.Nil(t, err)

	assert.Equal(t, primitives.WeightFromParts(2_091_362, 0), call.BaseWeight())
}

func Test_Call_Remark_WeighData(t *testing.T) {
	call := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, primitives.WeightFromParts(567, 0), call.WeighData(baseWeight))
}

func Test_Call_Remark_ClassifyDispatch(t *testing.T) {
	call := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, primitives.NewDispatchClassNormal(), call.ClassifyDispatch(baseWeight))
}

func Test_Call_Remark_PaysFee(t *testing.T) {
	call := newCallRemark(moduleId, functionRemarkIndex)

	assert.Equal(t, primitives.NewPaysYes(), call.PaysFee(baseWeight))
}

func Test_Call_Remark_Dispatch_Success(t *testing.T) {
	call := newCallRemark(moduleId, functionRemarkIndex)

	result := call.Dispatch(primitives.NewRawOriginRoot(), nil)

	assert.Equal(t, primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, result)
}

func Test_Call_Remark_Dispatch_Fail(t *testing.T) {
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	call := newCallRemark(moduleId, functionRemarkIndex)

	result := call.Dispatch(primitives.NewRawOriginNone(), nil)

	assert.Equal(t, expected, result)
}

func Test_EnsureSignedOrRoot_Root(t *testing.T) {
	r, err := EnsureSignedOrRoot(primitives.NewRawOriginRoot())

	assert.Nil(t, err)
	assert.Equal(t, sc.NewOption[primitives.AccountId[primitives.PublicKey]](nil), r)
}

func Test_EnsureSignedOrRoot_Signed(t *testing.T) {
	slice := make([]sc.U8, 32)
	seq := sc.NewFixedSequence[sc.U8](32, slice...)
	address, err := primitives.NewEd25519PublicKey(seq...)
	assert.Nil(t, err)
	signer := primitives.NewAccountId[primitives.PublicKey](address)

	r, e := EnsureSignedOrRoot(primitives.NewRawOriginSigned(signer))

	assert.Nil(t, e)
	assert.Equal(t, sc.NewOption[primitives.AccountId[primitives.PublicKey]](signer), r)
}

func Test_EnsureSignedOrRoot_BadOrigin(t *testing.T) {
	r, err := EnsureSignedOrRoot(primitives.NewRawOriginNone())

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), err)
	assert.Equal(t, sc.NewOption[primitives.AccountId[primitives.PublicKey]](nil), r)
}
