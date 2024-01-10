//go:build benchmarks

package env

//go:wasmimport env ext_benchmarking_current_time_version_1
func ExtBenchmarkingCurrentTimeVersion1() int64

//go:wasmimport env ext_benchmarking_set_whitelist_version_1
func ExtBenchmarkingSetWhitelistVersion1(key int64)

//go:wasmimport env ext_benchmarking_reset_read_write_count_version_1
func ExtBenchmarkingResetReadWriteCountVersion1()

//go:wasmimport env ext_benchmarking_start_db_tracker_version_1
func ExtBenchmarkingStartDbTrackerVersion1()

//go:wasmimport env ext_benchmarking_stop_db_tracker_version_1
func ExtBenchmarkingStopDbTrackerVersion1()

//go:wasmimport env ext_benchmarking_wipe_db_version_1
func ExtBenchmarkingWipeDbVersion1()

//go:wasmimport env ext_benchmarking_commit_db_version_1
func ExtBenchmarkingCommitDbVersion1()

//go:wasmimport env ext_benchmarking_db_read_count_version_1
func ExtBenchmarkingDbReadCountVersion1() int32

//go:wasmimport env ext_benchmarking_db_write_count_version_1
func ExtBenchmarkingDbWriteCountVersion1() int32
