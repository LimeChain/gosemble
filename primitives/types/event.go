package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Event = sc.VaryingData

type EventRecord struct {
	Phase  ExtrinsicPhase
	Event  Event
	Topics sc.Sequence[H256]
}

func (er EventRecord) Encode(buffer *bytes.Buffer) {
	er.Phase.Encode(buffer)
	er.Event.Encode(buffer)
	er.Topics.Encode(buffer)
}

func (er EventRecord) Bytes() []byte {
	return sc.EncodedBytes(er)
}
