package types

import (
	"bytes"
	"errors"
	"sort"

	sc "github.com/LimeChain/goscale"
)

type InherentData struct {
	data map[[8]byte]sc.Sequence[sc.U8]
}

func NewInherentData() *InherentData {
	return &InherentData{
		data: make(map[[8]byte]sc.Sequence[sc.U8]),
	}
}

func (id *InherentData) Encode(buffer *bytes.Buffer) error {
	err := sc.ToCompact(uint64(len(id.data))).Encode(buffer)
	if err != nil {
		return err
	}

	keys := make([][8]byte, 0)
	for k := range id.data {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return string(keys[i][:]) < string(keys[j][:]) })

	for _, k := range keys {
		value := id.data[k]

		buffer.Write(k[:])
		buffer.Write(value.Bytes())
	}

	return nil
}

func (id *InherentData) Bytes() []byte {
	return sc.EncodedBytes(id)
}

func (id InherentData) Get(key [8]byte) sc.Sequence[sc.U8] {
	return id.data[key]
}

func (id *InherentData) Put(key [8]byte, value sc.Encodable) error {
	if id.data[key] != nil {
		return NewInherentErrorInherentDataExists(sc.BytesToFixedSequenceU8(key[:]))
	}

	id.data[key] = sc.BytesToSequenceU8(value.Bytes())

	return nil
}

func (id *InherentData) Clear() {
	id.data = make(map[[8]byte]sc.Sequence[sc.U8])
}

func DecodeInherentData(buffer *bytes.Buffer) (*InherentData, error) {
	result := NewInherentData()
	lenCompact, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	length := lenCompact.ToBigInt().Int64()

	for i := 0; i < int(length); i++ {
		key := [8]byte{}
		len, err := buffer.Read(key[:])
		if err != nil {
			return nil, err
		}
		if len != 8 {
			return nil, errors.New("invalid length")
		}
		value, err := sc.DecodeSequence[sc.U8](buffer)
		if err != nil {
			return nil, err
		}

		result.data[key] = value
	}

	return result, nil
}
