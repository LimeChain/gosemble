package benchmarking

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/benchmarking"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type Module struct {
	systemModule  system.Module
	transactional support.Transactional[primitives.PostDispatchInfo]
	decoder       types.RuntimeDecoder
	memUtils      utils.WasmMemoryTranslator
	hashing       io.Hashing
	logger        log.Logger
}

func New(systemModule system.Module, decoder types.RuntimeDecoder, logger log.Logger) Module {
	return Module{
		systemModule:  systemModule,
		decoder:       decoder,
		transactional: support.NewTransactional[primitives.PostDispatchInfo](logger),
		memUtils:      utils.NewMemoryTranslator(),
		hashing:       io.NewHashing(),
		logger:        logger,
	}
}

// TODO:
// Implement DbCommit, DbWipe once the state implementation
// in Gossamer supports caching and nested transactions.
// https://github.com/ChainSafe/gossamer/discussions/3646

func (m Module) Run(dataPtr int32, dataLen int32) int64 {
	data := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	benchmarkConfig, err := benchmarking.DecodeBenchmarkConfig(buffer)
	if err != nil {
		m.logger.Critical(err.Error())
	}

	opaqueExtrinsic := sc.SequenceU8ToBytes(benchmarkConfig.Extrinsic)
	extrinsic, err := m.decoder.DecodeUncheckedExtrinsic(bytes.NewBuffer(opaqueExtrinsic))
	if err != nil {
		m.logger.Critical(err.Error())
	}

	function := extrinsic.Function()
	args := function.Args()
	origin, accountId := m.originAndMaybeAccount(benchmarkConfig)

	measuredDurations := []int64{}

	benchmarking.StoreSnapshotDb()

	// Always do at least one internal repeat.
	repeats := int(benchmarkConfig.InternalRepeats)
	if repeats < 1 {
		repeats = 1
	}
	for i := 1; i <= repeats; i++ {
		// The dispatch call is executed in a transactional context,
		// allowing to rollback and reset the state after each iteration.
		// as an alternative of providing before hook.

		benchmarking.RestoreSnapshotDb()

		// Does nothing, for now
		benchmarking.WipeDb()

		// Set up the externalities environment for the setup we want to
		// benchmark.

		// Sets the block number to 1 to allow emitting events
		m.systemModule.StorageBlockNumberSet(1)

		// Commit the externalities to the database, flushing the DB cache.
		// This will enable worst case scenario for reading from the database.
		// Does nothing, for now
		benchmarking.CommitDb()

		// Whitelist known storage keys.
		m.whitelistWellKnownKeys()

		// Whitelist the signer account key.
		keyStorageAccount := m.accountStorageKeyFrom(accountId.Value)
		benchmarking.SetWhitelist(keyStorageAccount)

		// Reset the read/write counter so we don't count
		// operations in the setup process.
		benchmarking.ResetReadWriteCount()

		benchmarking.StartDbTracker()

		var start, end int64

		start = benchmarking.CurrentTime()
		_, err := m.transactional.WithStorageLayer(
			func() (primitives.PostDispatchInfo, error) {
				return function.Dispatch(origin, args)
			},
		)
		end = benchmarking.CurrentTime()
		if err != nil {
			m.logger.Critical(err.Error())
		}

		// Calculate the diff caused by the benchmark.
		measuredDurations = append(measuredDurations, end-start)

		benchmarking.StopDbTracker()

		// Commit the changes to get proper write count.
		// Does nothing, for now
		benchmarking.CommitDb()
	}

	// Calculate the average time.
	extrinsicTime := calculateAverageTime(measuredDurations)

	benchmarkResult := benchmarking.BenchmarkResult{
		ExtrinsicTime: sc.NewU128(extrinsicTime),
		Reads:         sc.U32(benchmarking.DbReadCount()),
		Writes:        sc.U32(benchmarking.DbWriteCount()),
	}.Bytes()

	return m.memUtils.BytesToOffsetAndSize(benchmarkResult)
}

func calculateAverageTime(durations []int64) int64 {
	var sum int64
	for _, duration := range durations {
		sum += duration
	}
	return sum / int64(len(durations))
}

func (m Module) originAndMaybeAccount(benchmarkConfig benchmarking.BenchmarkConfig) (primitives.RawOrigin, sc.Option[primitives.AccountId]) {
	origin := benchmarkConfig.Origin

	if origin.IsSignedOrigin() {
		id, err := origin.AsSigned()
		if err != nil {
			m.logger.Critical(err.Error())
		}
		return origin, sc.NewOption[primitives.AccountId](id)
	}

	return origin, sc.NewOption[primitives.AccountId](nil)
}

func (m Module) whitelistWellKnownKeys() {
	keySystemHash := m.hashing.Twox128([]byte("System"))
	keyBlockWeight := m.hashing.Twox128([]byte("BlockWeight"))
	keyExecutionPhaseHash := m.hashing.Twox128([]byte("ExecutionPhase"))
	keyEventCountHash := m.hashing.Twox128([]byte("EventCount"))
	keyEventsHash := m.hashing.Twox128([]byte("Events"))
	keyNumberHash := m.hashing.Twox128([]byte("Number"))
	keyTotalIssuanceHash := m.hashing.Twox128([]byte("TotalIssuance"))

	benchmarking.SetWhitelist(append(keySystemHash, keyTotalIssuanceHash...))
	benchmarking.SetWhitelist(append(keySystemHash, keyBlockWeight...))
	benchmarking.SetWhitelist(append(keySystemHash, keyNumberHash...))
	benchmarking.SetWhitelist(append(keySystemHash, keyExecutionPhaseHash...))
	benchmarking.SetWhitelist(append(keySystemHash, keyEventCountHash...))
	benchmarking.SetWhitelist(append(keySystemHash, keyEventsHash...))

	benchmarking.SetWhitelist([]byte(":transaction_level:"))
	benchmarking.SetWhitelist([]byte(":extrinsic_index"))
	benchmarking.SetWhitelist([]byte(":intrablock_entropy"))
}

func (m Module) accountStorageKeyFrom(address primitives.AccountId) []byte {
	addressBytes := address.FixedSequence.Bytes()
	keySystemHash := m.hashing.Twox128([]byte("System"))
	keyAccountHash := m.hashing.Twox128([]byte("Account"))
	addressHash := m.hashing.Blake128(addressBytes)
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, addressHash...)
	keyStorageAccount = append(keyStorageAccount, addressBytes...)
	return keyStorageAccount
}