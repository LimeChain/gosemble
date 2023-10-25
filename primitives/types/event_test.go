package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedEventRecordBytes, _ = hex.DecodeString("02010203040500")
)

var (
	event1 = Event{
		sc.U8(1),
		sc.U8(2),
		sc.FixedSequence[sc.U8]{3, 4},
		sc.U8(5),
	}

	event2 = Event{
		sc.U8(6),
		sc.U8(7),
		sc.U8(8),
	}

	eventRecord1 = EventRecord{
		Phase:  NewExtrinsicPhaseInitialization(),
		Event:  event1,
		Topics: sc.Sequence[H256]{},
	}

	eventRecord2 = EventRecord{
		Phase:  NewExtrinsicPhaseInitialization(),
		Topics: sc.Sequence[H256]{},
	}

	eventRecords = sc.Sequence[EventRecord]{eventRecord2, eventRecord2}
)

func Test_NewEvent(t *testing.T) {
	expectedEvent := sc.NewVaryingData(sc.U8(1), sc.U8(2), sc.FixedSequence[sc.U8]{3, 4}, sc.U8(5))

	assert.Equal(t, expectedEvent, NewEvent(1, 2, sc.FixedSequence[sc.U8]{3, 4}, sc.U8(5)))
}

func Test_DecodeEvents(t *testing.T) {
	eventRecordsBytes := []byte{}
	for _, eventRecord := range eventRecords {
		eventRecordsBytes = append(eventRecordsBytes, eventRecord.Bytes()...)
	}
	eventRecordsBytes = append(sc.ToCompact(len(eventRecords)).Bytes(), eventRecordsBytes...)

	assert.Equal(t, eventRecords, DecodeEvents(bytes.NewBuffer(eventRecordsBytes)))
}

func Test_EventRecord_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	eventRecord1.Encode(buf)

	assert.Equal(t, expectedEventRecordBytes, buf.Bytes())
}

func Test_EventRecord_Bytes(t *testing.T) {
	assert.Equal(t, expectedEventRecordBytes, eventRecord1.Bytes())
}

func Test_DecodeEventRecord(t *testing.T) {
	buf := &bytes.Buffer{}

	eventRecord2.Encode(buf)

	assert.Equal(t, eventRecord2, DecodeEventRecord(buf))
}
