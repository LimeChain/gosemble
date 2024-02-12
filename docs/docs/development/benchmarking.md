---
layout: default
permalink: /development/benchmark
---

# Benchmarking ⏱️ 

The goal of benchmarking is to determine how computationally heavy it is to execute a given operation, measured as time, which reflects the operation's computational complexity. The execution time is represented as weight, and **1 second** of computation on a physical machine is equal to **10^12 weight** units. This measure is used to limit the amount of work that can be done in a single block and to charge fees proportionally to the resources consumed by the operation.

## Process 📌

Gosemble includes a CLI that provides a way of executing benchmark tests in a configurable manner, including extrinsics, steps, repeatability, etc. As a result, it automatically generates the weight files. This functionality relies on a set of utility functions provided by both the runtime and the host (Gossamer), allowing to measure the execution time in an isolated manner. It also accounts for database reads and writes of the storage keys hit during execution (some keys are preloaded and thus are excluded from the counts).
Here are the necessary steps to follow:

### 1. Switch the host branch 🔀

Checkout the host branch that contains the necessary functionality for benchmarking:

```bash
cd gossamer
git checkout bench-instance
```

Later, this set of functions provided by the host (Gossamer) is imported by the runtime and used during the benchmarking process.

```bash
"env"."ext_benchmarking_current_time_version_1": [] -> [I64]
"env"."ext_benchmarking_set_whitelist_version_1": [I64] -> []
"env"."ext_benchmarking_reset_read_write_count_version_1": [] -> []
"env"."ext_benchmarking_start_db_tracker_version_1": [] -> []
"env"."ext_benchmarking_stop_db_tracker_version_1": [] -> []
"env"."ext_benchmarking_db_read_count_version_1": [] -> [I32]
"env"."ext_benchmarking_db_write_count_version_1": [] -> [I32]
"env"."ext_benchmarking_wipe_db_version_1": [] -> []
"env"."ext_benchmarking_commit_db_version_1": [] -> []
"env"."ext_benchmarking_store_snapshot_db_version_1": [] -> []
"env"."ext_benchmarking_restore_snapshot_db_version_1": [] -> []
```

### 2. Build the runtime 🏗️

Build the runtime with the benchmarking feature:

```bash
make build-benchmarking
```
or
```bash
make build-docker-benchmarking
```

It exposes additional utility functions exported by the runtime, which allow the execution of benchmark tests in a Wasm environment.

```bash
"Benchmark_dispatch": [I32, I32] -> [I64]
"Benchmark_hook": [I32, I32] -> [I64]
```

### 3. Write benchmarks 📝

It is important to note that benchmark tests should always assess the **worst-case** scenario. The general process of writing a benchmark test includes setting up an initial state, executing an operation, and asserting the final state, which encompasses both success and failure scenarios.

### 3.1. Dispatch calls 📞

* Example benchmark test:
[benchmark_timestamp_set_test.go](https://github.com/LimeChain/gosemble/blob/develop/runtime/benchmark_timestamp_set_test.go)

* Example benchmark test with linear components:
[runtime/benchmark_system_remark_test.go](https://github.com/LimeChain/gosemble/blob/develop/runtime/benchmark_system_remark_test.go)

Extrinsic calls are executed through the `Benchmark_dispatch` runtime function.

### 3.2. System hooks 🪝

* Example benchmark test:
[benchmark_hooks_test.go](https://github.com/LimeChain/gosemble/blob/develop/runtime/benchmark_hooks_test.go)

System hooks are executed through the `Benchmark_hooks` runtime function.

### 3.3. Block overhead 🧊

* Example benchmark test:
[overhead_test.go](https://github.com/LimeChain/gosemble/blob/develop/benchmarking/overhead_test.go)

## 4. Run benchmarks ▶️

Run extrinsic and hook benchmarks with auto-generating weight files (the default):

```bash
make benchmark
```

```bash
make benchmark steps=50 repeat=100
```

Run the overhead benchmarks:

```bash
make benchmark-overhead
```

Run benchmarks without generating weight files:

```bash
GENERATE_WEIGHT_FILES=false make benchmark
GENERATE_WEIGHT_FILES=false make benchmark-overhead
```