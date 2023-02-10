package types

import (
	"bytes"
	"errors"
	"fmt"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/timestamp"
)

type InherentData struct {
	Data map[[8]byte]sc.Sequence[sc.U8]
}

func NewInherentData() *InherentData {
	return &InherentData{
		Data: make(map[[8]byte]sc.Sequence[sc.U8]),
	}
}

func (id *InherentData) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(uint64(len(id.Data))).Encode(buffer)

	for k, v := range id.Data {
		buffer.Write(k[:])
		buffer.Write(v.Bytes())
	}
}

func (id *InherentData) Bytes() []byte {
	buf := &bytes.Buffer{}

	sc.Compact(sc.NewU128FromUint64(uint64(len(id.Data)))).Encode(buf)

	for k, v := range id.Data {
		buf.Write(k[:])
		buf.Write(v.Bytes())
	}

	return buf.Bytes()
}

func (id *InherentData) Put(key [8]byte, value sc.Encodable) error {
	if id.Data[key] != nil {
		return errors.New(fmt.Sprintf("InherentDataExists - [%v]", key))
	}

	id.Data[key] = sc.BytesToSequenceU8(value.Bytes())

	return nil
}

func (id *InherentData) Clear() {
	id.Data = make(map[[8]byte]sc.Sequence[sc.U8])
}

func (id *InherentData) CheckExtrinsics(block Block) CheckInherentsResult {

	result := NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false

		call := extrinsic.Function

		switch call.CallIndex.ModuleIndex {
		case timestamp.ModuleIndex:
			if call.CallIndex.FunctionIndex == timestamp.FunctionIndex {
				isInherent = true
				err := timestamp.CheckInherent(call, *id)
				if err != nil {
					err := result.PutError(timestamp.InherentIdentifier, sc.Str(err.Error()))
					if err != nil {
						panic(err)
					}

					if result.FatalError {
						return result
					}
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	return result
}

func DecodeInherentData(buffer *bytes.Buffer) (*InherentData, error) {
	result := NewInherentData()
	length := sc.DecodeCompact(buffer).ToBigInt().Int64()

	for i := 0; i < int(length); i++ {
		key := [8]byte{}
		len, err := buffer.Read(key[:])
		if err != nil {
			return nil, err
		}
		if len != 8 {
			return nil, errors.New("invalid length")
		}
		value := sc.DecodeSequence[sc.U8](buffer)

		result.Data[key] = value
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
