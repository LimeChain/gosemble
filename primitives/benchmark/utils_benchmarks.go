//go:build benchmarks

package benchmark

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func TimeNow() int64 {
	return env.ExtBenchmarksTimeNowVersion1()
}

func DbWhitelistKey(key []byte) {
	keyOffsetSize := utils.NewMemoryTranslator().BytesToOffsetAndSize(key)
	env.ExtBenchmarksDbWhitelistKeyVersion1(keyOffsetSize)
}

func DbResetTracker() {
	env.ExtBenchmarksDbResetTrackerVersion1()
}

func DbStartTracker() {
	env.ExtBenchmarksDbStartTrackerVersion1()
}

func DbStopTracker() {
	env.ExtBenchmarksDbStopTrackerVersion1()
}

func DbWipe() {
	env.ExtBenchmarksDbWipeVersion1()
}

func DbCommit() {
	env.ExtBenchmarksDbCommitVersion1()
}

func DbReadCount() int32 {
	return env.ExtBenchmarksDbReadCountVersion1()
}

func DbWriteCount() int32 {
	return env.ExtBenchmarksDbWriteCountVersion1()
}
