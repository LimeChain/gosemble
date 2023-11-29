package transaction_payment

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_DecodeEvent(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	expectedEvent := NewEventTransactionFeePaid(moduleId, who, sc.NewU128(7), sc.NewU128(1))
	err := expectedEvent.Encode(buffer)
	assert.NoError(t, err)

	result, err := DecodeEvent(moduleId, buffer)
	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, result)
}

func Test_DecodeEvent_ModuleIndexError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	expectedEvent := NewEventTransactionFeePaid(moduleId, who, sc.NewU128(7), sc.NewU128(1))
	err := expectedEvent.Encode(buffer)
	assert.NoError(t, err)

	assert.PanicsWithValue(t, "invalid transaction_payment.Event module", func() {
		DecodeEvent(sc.U8(123), buffer)
	})
}

func Test_DecodeEvent_TypeError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	expectedEvent := types.NewEvent(moduleId, 99, who, sc.NewU128(7), sc.NewU128(1))

	err := expectedEvent.Encode(buffer)
	assert.NoError(t, err)

	assert.PanicsWithValue(t, "invalid transaction_payment.Event type", func() {
		DecodeEvent(moduleId, buffer)
	})
}
