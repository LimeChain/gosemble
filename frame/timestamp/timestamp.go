package timestamp

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

var Module = support.ModuleMetadata{
	Index: timestamp.ModuleIndex,
	Functions: map[string]support.FunctionMetadata{
		"set": {Index: timestamp.FunctionSetIndex, Func: Set},
	},
}

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
				ModuleIndex:   Module.Index,
				FunctionIndex: Module.Functions["set"].Index,
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

func Set(now sc.U64) {
	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	didUpdateHash := hashing.Twox128(constants.KeyDidUpdate)

	didUpdate := storage.Exists(append(timestampHash, didUpdateHash...))

	if didUpdate == 1 {
		log.Critical("Timestamp must be updated only once in the block")
	}

	nowHash := hashing.Twox128(constants.KeyNow)
	previousTimestamp := storage.GetDecode(append(timestampHash, nowHash...), sc.DecodeU64)

	if !(previousTimestamp == 0 || now >= previousTimestamp+timestamp.MinimumPeriod) {
		log.Critical("Timestamp must increment by at least <MinimumPeriod> between sequential blocks")
	}

	storage.Set(append(timestampHash, nowHash...), now.Bytes())
	storage.Set(append(timestampHash, didUpdateHash...), sc.Bool(true).Bytes())

	// TODO: Every consensus that uses the timestamp must implement
	// <T::OnTimestampSet as OnTimestampSet<_>>::on_timestamp_set(now)

	// TODO:
	// timestamp module should not depend on the aura module
	aura.OnTimestampSet(now)
}

func OnFinalize() {
	timestampHash := hashing.Twox128(constants.KeyTimestamp)
	didUpdateHash := hashing.Twox128(constants.KeyDidUpdate)

	didUpdate := storage.Get(append(timestampHash, didUpdateHash...))

	if didUpdate.HasValue {
		storage.Clear(append(timestampHash, didUpdateHash...))
	} else {
		log.Critical("Timestamp must be updated once in the block")
	}
}
