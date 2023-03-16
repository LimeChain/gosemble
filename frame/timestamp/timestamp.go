package timestamp

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

func CreateInherent(inherent types.InherentData) []byte {
	inherentData := inherent.Data[timestamp.InherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.
	buffer.Reset()

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	nowHash := hashing.Twox128(constants.KeyNow)

	nextTimestamp := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64) + timestamp.MinimumPeriod

	if ts > nextTimestamp {
		nextTimestamp = ts
	}

	extrinsic := types.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   Module.Index(),
				FunctionIndex: Module.Set.Index(),
			},
			Args: sc.BytesToSequenceU8(sc.ToCompact(uint64(nextTimestamp)).Bytes()),
		},
	}

	return extrinsic.Bytes()
}

func CheckInherent(call types.Call, inherent types.InherentData) error {
	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(call.Args))
	compactTimestamp := sc.DecodeCompact(buffer)
	t := sc.U64(compactTimestamp.ToBigInt().Uint64())
	buffer.Reset()

	inherentData := inherent.Data[timestamp.InherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.
	buffer.Reset()

	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	nowHash := hashing.Twox128(constants.KeyNow)
	systemNow := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	minimum := systemNow + timestamp.MinimumPeriod
	if t > ts+timestamp.MaxTimestampDriftMillis {
		return types.NewTimestampErrorTooFarInFuture()
	} else if t < minimum {
		return types.NewTimestampErrorTooEarly()
	}

	return nil
}
