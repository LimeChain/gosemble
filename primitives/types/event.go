package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Event = sc.VaryingData

func NewEvent(module sc.U8, event sc.U8, values ...sc.Encodable) Event {
	args := []sc.Encodable{module, event}
	args = append(args, values...)
	return sc.NewVaryingData(args...)
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

func DecodeEventRecord(buffer *bytes.Buffer) (EventRecord, error) {
	phase, err := DecodeExtrinsicPhase(buffer)
	if err != nil {
		return EventRecord{}, nil
	}
	topics, err := sc.DecodeSequence[H256](buffer)
	if err != nil {
		return EventRecord{}, nil
	}
	return EventRecord{
		Phase:  phase,
		Event:  nil, // TODO:
		Topics: topics,
	}, nil
}

// func DecodeEvents(buffer *bytes.Buffer) sc.Sequence[EventRecord] {
// 	compactSize := sc.DecodeCompact(buffer)
// 	size := int(compactSize.ToBigInt().Int64())

// 	sequence := make(sc.Sequence[EventRecord], size)
// 	for i := 0; i < size; i++ {
// 		sequence[i] = DecodeEventRecord(buffer)
// 	}

// 	return sequence
// }
