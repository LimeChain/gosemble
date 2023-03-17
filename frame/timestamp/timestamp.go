package timestamp

import (
	"bytes"

	types2 "github.com/LimeChain/gosemble/execution/types"

	"github.com/LimeChain/gosemble/frame/timestamp/module"

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

	extrinsic := types2.UncheckedExtrinsic{
		Version: types.ExtrinsicFormatVersion,
		Function: types2.Call{
			CallIndex: types.CallIndex{
				ModuleIndex:   module.Module.Index(),
				FunctionIndex: module.Module.Set.Index(),
			},
			Args: []sc.Encodable{sc.ToCompact(uint64(nextTimestamp))},
		},
	}

	return extrinsic.Bytes()
}

func CheckInherent(args []sc.Encodable, inherent types.InherentData) error {
	t := args[0].(sc.U64)

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
	systemNow := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	minimum := systemNow + timestamp.MinimumPeriod
	if t > ts+timestamp.MaxTimestampDriftMillis {
		return types.NewTimestampErrorTooFarInFuture()
	} else if t < minimum {
		return types.NewTimestampErrorTooEarly()
	}

	return nil
}
