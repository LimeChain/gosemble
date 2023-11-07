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
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventEndowed.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventEndowed, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_DustLost(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDustLost.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDustLost, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Transfer(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventTransfer.Bytes())
	buffer.Write(fromAddress.AsAccountId().Bytes())
	buffer.Write(toAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, _ := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventTransfer, fromAddress.AsAccountId(), toAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_BalanceSet(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventBalanceSet.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(newFree.Bytes())
	buffer.Write(newReserved.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventBalanceSet, targetAddress.AsAccountId(), newFree, newReserved),
		result,
	)
}

func Test_Balances_DecodeEvent_Reserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserved.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventReserved, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Unreserved(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventUnreserved.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventUnreserved, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_ReserveRepatriated(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventReserveRepatriated.Bytes())
	buffer.Write(fromAddress.AsAccountId().Bytes())
	buffer.Write(toAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())
	buffer.Write(types.BalanceStatusFree.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(
			sc.U8(moduleId),
			EventReserveRepatriated,
			fromAddress.AsAccountId(),
			toAddress.AsAccountId(),
			targetValue, types.BalanceStatusFree),
		result,
	)
}

func Test_Balances_DecodeEvent_Deposit(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventDeposit.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventDeposit, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Withdraw(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventWithdraw.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventWithdraw, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_Slashed(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.Write(EventSlashed.Bytes())
	buffer.Write(targetAddress.AsAccountId().Bytes())
	buffer.Write(targetValue.Bytes())

	result, err := DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	assert.Nil(t, err)

	assert.Equal(t,
		sc.NewVaryingData(sc.U8(moduleId), EventSlashed, targetAddress.AsAccountId(), targetValue),
		result,
	)
}

func Test_Balances_DecodeEvent_InvalidModule_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	assert.PanicsWithValue(t, errInvalidEventModule, func() {
		DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	})
}

func Test_Balances_DecodeEvent_InvalidType_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(moduleId)
	buffer.WriteByte(255)

	assert.PanicsWithValue(t, errInvalidEventType, func() {
		DecodeEvent[primitives.Ed25519Signer](moduleId, buffer)
	})
}
