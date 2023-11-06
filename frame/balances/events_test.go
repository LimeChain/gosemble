package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_DecodeEvent_Endowed(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventEndowed.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventEndowed, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_DustLost(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDustLost.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDustLost, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Transfer(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventTransfer.Bytes())
	buffer.Write(fromAddress32.Bytes())
	buffer.Write(toAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, _ := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventTransfer, fromAddress32.FixedSequence, toAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_BalanceSet(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventBalanceSet.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(newFree.Bytes())
	buffer.Write(newReserved.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventBalanceSet, targetAddress32.FixedSequence, newFree, newReserved),
		result,
	)
}

func Test_Balances_DecodeEvent_Reserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserved.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventReserved, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Unreserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventUnreserved.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventUnreserved, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_ReserveRepatriated(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserveRepatriated.Bytes())
	buffer.Write(fromAddress32.Bytes())
	buffer.Write(toAddress32.Bytes())
	buffer.Write(targetValue.Bytes())
	buffer.Write(types.BalanceStatusFree.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(
			sc.U8(moduleId),
			EventReserveRepatriated,
			fromAddress32.FixedSequence,
			toAddress32.FixedSequence,
			targetValue, types.BalanceStatusFree),
		result,
	)
}

func Test_Balances_DecodeEvent_Deposit(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDeposit.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDeposit, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Withdraw(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventWithdraw.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventWithdraw, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Slashed(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventSlashed.Bytes())
	buffer.Write(targetAddress32.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventSlashed, targetAddress32.FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_InvalidModule_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	assert.PanicsWithValue(t, errInvalidEventModule, func() {
		DecodeEvent(moduleId, buffer)
	})
}

func Test_Balances_DecodeEvent_InvalidType_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.WriteByte(255)

	assert.PanicsWithValue(t, errInvalidEventType, func() {
		DecodeEvent(moduleId, buffer)
	})
}
