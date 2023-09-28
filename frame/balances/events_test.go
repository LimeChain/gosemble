package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_DecodeEvent_Endowed(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventEndowed.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventEndowed, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_DustLost(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDustLost.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDustLost, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Transfer(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventTransfer.Bytes())
	buffer.Write(fromAddress.AsAddress32().Bytes())
	buffer.Write(toAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventTransfer, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_BalanceSet(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventBalanceSet.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(newFree.Bytes())
	buffer.Write(newReserved.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventBalanceSet, targetAddress.AsAddress32().FixedSequence, newFree, newReserved),
		result,
	)
}

func Test_Balances_DecodeEvent_Reserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserved.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventReserved, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Unreserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventUnreserved.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventUnreserved, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_ReserveRepatriated(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserveRepatriated.Bytes())
	buffer.Write(fromAddress.AsAddress32().Bytes())
	buffer.Write(toAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())
	buffer.Write(primitives.BalanceStatusFree.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(
			sc.U8(moduleId),
			EventReserveRepatriated,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			targetValue, primitives.BalanceStatusFree),
		result,
	)
}

func Test_Balances_DecodeEvent_Deposit(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDeposit.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDeposit, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Withdraw(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventWithdraw.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventWithdraw, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Slashed(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventSlashed.Bytes())
	buffer.Write(targetAddress.AsAddress32().Bytes())
	buffer.Write(targetValue.Bytes())

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventSlashed, targetAddress.AsAddress32().FixedSequence, targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	assert.PanicsWithValue(t, errInvalidEventModule, func() {
		DecodeEvent(moduleId, buffer)
	})
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
