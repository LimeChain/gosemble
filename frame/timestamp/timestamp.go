package timestamp

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

const (
	MinimumPeriod = 1 * 1000 // 1 second
	ModuleIndex   = 3
	FunctionIndex = 0
)

var (
	InherentIdentifier = [8]byte{'t', 'i', 'm', 's', 't', 'a', 'p', '0'}
)

func CreateInherent(inherent types.InherentData) []byte {
	inherentData := inherent.Data[InherentIdentifier]

	if inherentData == nil {
		panic("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	timestamp := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.
	buffer.Reset()

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	nowHash := hashing.Twox128(constants.KeyNow)

	nextTimestamp := timestamp

	nowBytes := storage.Get(append(timestampHash, nowHash...))
	if len(nowBytes) > 1 {
		buffer.Write(nowBytes)
		nowTimestamp := sc.DecodeU64(buffer)
		buffer.Reset()

		nextTimestamp = nowTimestamp + MinimumPeriod
	}

	extrinsic := types.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   ModuleIndex,
				FunctionIndex: FunctionIndex,
			},
			Args: sc.BytesToSequenceU8(nextTimestamp.Bytes()),
		},
	}

	return extrinsic.Bytes()
}
