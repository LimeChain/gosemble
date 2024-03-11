package types

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedEventRecordBytes, _      = hex.DecodeString("02010200000008030400")
	eventRecordInvalidPhaseBytes, _  = hex.DecodeString("03010200000008030400")
	eventRecordInvalidTipicsBytes, _ = hex.DecodeString("020102000000080304")
)

var (
	event1 = Event{sc.VaryingData{
		sc.U8(1),
		sc.U32(2),
		sc.Sequence[sc.U8]{3, 4},
	}}

	decodeEvent1 = func(moduleIndex sc.U8, buffer *bytes.Buffer) (Event, error) {
		sc.DecodeU8(buffer)
		sc.DecodeU32(buffer)
		sc.DecodeSequence[sc.U8](buffer)
		return event1, nil
	}

	eventRecord1 = EventRecord{
		Phase:  NewExtrinsicPhaseInitialization(),
		Event:  event1,
		Topics: sc.Sequence[H256]{},
	}
)

func Test_NewEvent(t *testing.T) {
	expectedEvent := Event{sc.NewVaryingData(sc.U8(1), sc.U8(2), sc.FixedSequence[sc.U8]{3, 4})}

	assert.Equal(t, expectedEvent, NewEvent(1, 2, sc.FixedSequence[sc.U8]{3, 4}))
}

func Test_EventRecord_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	err := eventRecord1.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedEventRecordBytes, buf.Bytes())
}

func Test_EventRecord_Bytes(t *testing.T) {
	assert.Equal(t, expectedEventRecordBytes, eventRecord1.Bytes())
}

func Test_DecodeEventRecord(t *testing.T) {
	buf := bytes.NewBuffer(expectedEventRecordBytes)

	decodedEventRecord, err := DecodeEventRecord(1, decodeEvent1, buf)

	assert.NoError(t, err)
	assert.Equal(t, eventRecord1, decodedEventRecord)
}

func Test_DecodeEventRecord_Fails_To_Decode_Phase(t *testing.T) {
	buf := bytes.NewBuffer(eventRecordInvalidPhaseBytes)

	_, err := DecodeEventRecord(1, decodeEvent1, buf)

	assert.Equal(t, newTypeError("ExtrinsicPhase"), err)
}

func Test_DecodeEventRecord_Fails_To_Decode_Event(t *testing.T) {
	buf := bytes.NewBuffer(expectedEventRecordBytes)

	expectedError := newTypeError("Event")

	_, err := DecodeEventRecord(1, func(moduleIndex sc.U8, buffer *bytes.Buffer) (Event, error) {
		return Event{}, expectedError
	}, buf)

	assert.Equal(t, expectedError, err)
}

func Test_DecodeEventRecord_Fails_To_Decode_Topics(t *testing.T) {
	buf := bytes.NewBuffer(eventRecordInvalidTipicsBytes)

	_, err := DecodeEventRecord(1, decodeEvent1, buf)

	assert.Equal(t, io.EOF, err)
}
