package benchmarking

import (
	"bytes"

	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"

	sc "github.com/LimeChain/goscale"
)

// TODO:
// Execute on a range of components (values per component are set by default to 6)

// Configuration used to setup and run runtime benchmarks.
type BenchmarkConfig struct {
	// The encoded name of the pallet to benchmark.
	// Module sc.Sequence[sc.U8]

	// The encoded name of the benchmark/extrinsic to run.
	Extrinsic sc.Sequence[sc.U8]

	Origin sc.Option[types.RawOrigin]

	// The selected component values to use when running the benchmark.
	// SelectedComponents Vec<(BenchmarkParameter, u32)>

	// Enable an extra benchmark iteration which runs the verification logic for a benchmark.
	// Verify bool

	// Number of times to repeat benchmark within the Wasm environment. (versus in the client)
	InternalRepeats sc.U32
}

func (bc BenchmarkConfig) Encode(buffer *bytes.Buffer) {
	bc.InternalRepeats.Encode(buffer)
	bc.Origin.Encode(buffer)
	bc.Extrinsic.Encode(buffer)
}

func (bc BenchmarkConfig) Bytes() []byte {
	buffer := bytes.Buffer{}
	bc.Encode(&buffer)
	return buffer.Bytes()
}

func DecodeBenchmarkConfig(buffer *bytes.Buffer) (BenchmarkConfig, error) {
	internalRepeats, err := sc.DecodeU32(buffer)
	if err != nil {
		return BenchmarkConfig{}, err
	}

	origin, err := sc.DecodeOptionWith(buffer, func(buffer *bytes.Buffer) (types.RawOrigin, error) {
		return types.DecodeRawOrigin(buffer)
	})
	if err != nil {
		return BenchmarkConfig{}, err
	}

	extrinsic, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return BenchmarkConfig{}, err
	}

	return BenchmarkConfig{
		InternalRepeats: internalRepeats,
		Extrinsic:       extrinsic,
		Origin:          origin,
	}, nil
}

// Result from running benchmarks on a FRAME pallet.
// Contains duration of the function call in nanoseconds along with the benchmark parameters
// used for that benchmark result.
type BenchmarkResult struct {
	// Components Vec<(BenchmarkParameter, sc.U32)>
	ExtrinsicTime sc.U128
	// StorageRootTime sc.U128
	Reads sc.U32
	// RepeatReads sc.U32
	Writes sc.U32
	// RepeatWrites sc.U32
	// ProofSize sc.U32
	// Keys Vec<(Vec<u8>, u32, u32, bool)> // Skip
}

func (br BenchmarkResult) Encode(buffer *bytes.Buffer) {
	br.ExtrinsicTime.Encode(buffer)
	br.Reads.Encode(buffer)
	br.Writes.Encode(buffer)
}

func (br BenchmarkResult) Bytes() []byte {
	buffer := bytes.Buffer{}
	br.Encode(&buffer)
	return buffer.Bytes()
}

func DecodeBenchmarkResult(buffer *bytes.Buffer) (BenchmarkResult, error) {
	extrinsicTime, err := sc.DecodeU128(buffer)
	if err != nil {
		return BenchmarkResult{}, err
	}
	reads, err := sc.DecodeU32(buffer)
	if err != nil {
		return BenchmarkResult{}, err
	}
	writes, err := sc.DecodeU32(buffer)
	if err != nil {
		return BenchmarkResult{}, err
	}
	return BenchmarkResult{
		ExtrinsicTime: extrinsicTime,
		Reads:         reads,
		Writes:        writes,
	}, nil
}

func CurrentTime() int64 {
	return env.ExtBenchmarkingCurrentTimeVersion1()
}

func SetWhitelist(key []byte) {
	keyOffsetSize := utils.NewMemoryTranslator().BytesToOffsetAndSize(key)
	env.ExtBenchmarkingSetWhitelistVersion1(keyOffsetSize)
}

func ResetReadWriteCount() {
	env.ExtBenchmarkingResetReadWriteCountVersion1()
}

func StartDbTracker() {
	env.ExtBenchmarkingStartDbTrackerVersion1()
}

func StopDbTracker() {
	env.ExtBenchmarkingStopDbTrackerVersion1()
}

func WipeDb() {
	env.ExtBenchmarkingWipeDbVersion1()
}

func CommitDb() {
	env.ExtBenchmarkingCommitDbVersion1()
}

func DbReadCount() int32 {
	return env.ExtBenchmarkingDbReadCountVersion1()
}

func DbWriteCount() int32 {
	return env.ExtBenchmarkingDbWriteCountVersion1()
}
