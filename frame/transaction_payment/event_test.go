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
	expectedEvent := NewEventTransactionFeePaid(moduleId, who.FixedSequence, sc.NewU128(7), sc.NewU128(1))
	expectedEvent.Encode(buffer)

	result := DecodeEvent(moduleId, buffer)

	assert.Equal(t, expectedEvent, result)
}

func Test_DecodeEvent_ModuleIndexError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	expectedEvent := NewEventTransactionFeePaid(moduleId, who.FixedSequence, sc.NewU128(7), sc.NewU128(1))
	expectedEvent.Encode(buffer)

	assert.PanicsWithValue(t, "invalid transaction_payment.Event module", func() {
		DecodeEvent(sc.U8(123), buffer)
	})
}

func Test_DecodeEvent_TypeError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	expectedEvent := types.NewEvent(moduleId, 99, who.FixedSequence, sc.NewU128(7), sc.NewU128(1))
	expectedEvent.Encode(buffer)

	assert.PanicsWithValue(t, "invalid transaction_payment.Event type", func() {
		DecodeEvent(moduleId, buffer)
	})
}
