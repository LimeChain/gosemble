package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedEventRecordBytes, _ = hex.DecodeString("02010200000008030400")
)

var (
	event1 = Event{
		sc.U8(1),
		sc.U32(2),
		sc.Sequence[sc.U8]{3, 4},
	}

	eventRecord1 = EventRecord{
		Phase:  NewExtrinsicPhaseInitialization(),
		Event:  event1,
		Topics: sc.Sequence[H256]{},
	}
)

func Test_NewEvent(t *testing.T) {
	expectedEvent := sc.NewVaryingData(sc.U8(1), sc.U8(2), sc.FixedSequence[sc.U8]{3, 4})

	assert.Equal(t, expectedEvent, NewEvent(1, 2, sc.FixedSequence[sc.U8]{3, 4}))
}

func Test_EventRecord_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	eventRecord1.Encode(buf)

	assert.Equal(t, expectedEventRecordBytes, buf.Bytes())
}

func Test_EventRecord_Bytes(t *testing.T) {
	assert.Equal(t, expectedEventRecordBytes, eventRecord1.Bytes())
}
