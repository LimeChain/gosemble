package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Event struct {
	sc.VaryingData
}

func NewEvent(module sc.U8, event sc.U8, values ...sc.Encodable) Event {
	args := []sc.Encodable{module, event}
	args = append(args, values...)
	return Event{sc.NewVaryingData(args...)}
}

type EventRecord struct {
	Phase  ExtrinsicPhase
	Event  Event
	Topics sc.Sequence[H256]
}

func (er EventRecord) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		er.Phase,
		er.Event,
		er.Topics,
	)
}

func (er EventRecord) Bytes() []byte {
	return sc.EncodedBytes(er)
}

func DecodeEventRecord(
	moduleIndex sc.U8,
	decodeEvent func(moduleIndex sc.U8, buffer *bytes.Buffer) (Event, error),
	buffer *bytes.Buffer,
) (EventRecord, error) {
	phase, err := DecodeExtrinsicPhase(buffer)
	if err != nil {
		return EventRecord{}, err
	}

	event, err := decodeEvent(moduleIndex, buffer)
	if err != nil {
		return EventRecord{}, err
	}

	topics, err := sc.DecodeSequence[H256](buffer)
	if err != nil {
		return EventRecord{}, err
	}

	return EventRecord{
		Phase:  phase,
		Event:  event,
		Topics: topics,
	}, nil
}
