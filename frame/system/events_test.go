package system

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	targetAccount = constants.OneAddress
)

func Test_System_DecodeEvent_ExtrinsicSuccess(t *testing.T) {
	dispatchInfo := types.DispatchInfo{
		Weight:  baseWeight,
		Class:   types.NewDispatchClassOperational(),
		PaysFee: types.NewPaysNo(),
	}
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventExtrinsicSuccess.Bytes())
	buffer.Write(dispatchInfo.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventExtrinsicSuccess, dispatchInfo),
		result,
	)
}

func Test_System_DecodeEvent_ExtrinsicFailed(t *testing.T) {
	dispatchInfo := types.DispatchInfo{
		Weight:  baseWeight,
		Class:   types.NewDispatchClassOperational(),
		PaysFee: types.NewPaysNo(),
	}
	dispatchError := types.NewDispatchErrorBadOrigin()
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventExtrinsicFailed.Bytes())
	buffer.Write(dispatchError.Bytes())
	buffer.Write(dispatchInfo.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventExtrinsicFailed, dispatchError, dispatchInfo),
		result,
	)
}

func Test_System_DecodeEvent_CodeUpdated(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventCodeUpdated.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventCodeUpdated),
		result,
	)
}

func Test_System_DecodeEvent_NewAccount(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventNewAccount.Bytes())
	buffer.Write(targetAccount.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventNewAccount, targetAccount.FixedSequence),
		result,
	)
}

func Test_System_DecodeEvent_KilledAccount(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventKilledAccount.Bytes())
	buffer.Write(targetAccount.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventKilledAccount, targetAccount.FixedSequence),
		result,
	)
}

func Test_System_DecodeEvent_Remarked(t *testing.T) {
	emptyHash := [32]sc.U8{}
	hash := types.NewH256(emptyHash[:]...)
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventRemarked.Bytes())
	buffer.Write(targetAccount.Bytes())
	buffer.Write(hash.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventRemarked, targetAccount.FixedSequence, hash),
		result,
	)
}

func Test_System_DecodeEvent_InvalidModule_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	assert.PanicsWithValue(t, errInvalidEventModule, func() {
		DecodeEvent(moduleId, buffer)
	})
}

func Test_System_DecodeEvent_InvalidType_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.WriteByte(255)

	assert.PanicsWithValue(t, errInvalidEventType, func() {
		DecodeEvent(moduleId, buffer)
	})
}
