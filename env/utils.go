//go:build !benchmarks

package env

//go:wasmimport env ext_benchmarks_time_now_version_1
func ExtBenchmarksTimeNowVersion1() int64 {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_whitelist_key_version_1
func ExtBenchmarksDbWhitelistKeyVersion1(key int64) {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_reset_tracker_version_1
func ExtBenchmarksDbResetTrackerVersion1() {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_start_tracker_version_1
func ExtBenchmarksDbStartTrackerVersion1() {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_stop_tracker_version_1
func ExtBenchmarksDbStopTrackerVersion1() {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_wipe_version_1
func ExtBenchmarksDbWipeVersion1() {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_commit_version_1
func ExtBenchmarksDbCommitVersion1() {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_read_count_version_1
func ExtBenchmarksDbReadCountVersion1() int32 {
	panic("not implemented")
}

//go:wasmimport env ext_benchmarks_db_write_count_version_1
func ExtBenchmarksDbWriteCountVersion1() int32 {
	panic("not implemented")
}
