package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_Balances_DecodeEvent_Endowed(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventEndowed.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventEndowed, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_DustLost(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDustLost.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventDustLost, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_Transfer(t *testing.T) {
	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressAccountId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventTransfer.Bytes())
	buffer.Write(fromAddressId.Bytes())
	buffer.Write(toAddressAccountId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, _ := DecodeEvent(moduleId, buffer)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventTransfer, fromAddressId, toAddressAccountId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_BalanceSet(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventBalanceSet.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(newFree.Bytes())
	buffer.Write(newReserved.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventBalanceSet, targetAddressId, newFree, newReserved)},
		result,
	)
}

func Test_Balances_DecodeEvent_Reserved(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserved.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventReserved, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_Unreserved(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventUnreserved.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventUnreserved, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_ReserveRepatriated(t *testing.T) {
	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)
	toAddressAccountId, err := toAddress.AsAccountId()
	assert.Nil(t, err)
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserveRepatriated.Bytes())
	buffer.Write(fromAddressId.Bytes())
	buffer.Write(toAddressAccountId.Bytes())
	buffer.Write(targetValue.Bytes())
	buffer.Write(types.BalanceStatusFree.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(
			sc.U8(moduleId),
			EventReserveRepatriated,
			fromAddressId,
			toAddressAccountId,
			targetValue, types.BalanceStatusFree)},
		result,
	)
}

func Test_Balances_DecodeEvent_Deposit(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDeposit.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventDeposit, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_Withdraw(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventWithdraw.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventWithdraw, targetAddressId, targetValue)},
		result,
	)
}

func Test_Balances_DecodeEvent_Slashed(t *testing.T) {
	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventSlashed.Bytes())
	buffer.Write(targetAddressId.Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent(moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		primitives.Event{sc.NewVaryingData(sc.U8(moduleId), EventSlashed, targetAddressId, targetValue)},
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
