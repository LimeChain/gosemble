package benchmarking

import "flag"

// cmd flags and other options related to benchmarking
var Config = initBenchmarkingConfig()

type benchmarkingConfig struct {
	Steps, Repeat, HeapPages, DbCache      int
	WasmRuntime, GC, TinyGoVersion, Target string
	GenerateWeightFiles                    bool
}

func initBenchmarkingConfig() benchmarkingConfig {
	cfg := &benchmarkingConfig{}
	cfg.WasmRuntime = "../build/runtime.wasm"
	flag.IntVar(&cfg.Steps, "steps", 50, "Select how many samples we should take across the variable components.")
	flag.IntVar(&cfg.Repeat, "repeat", 20, "Select how many repetitions of this benchmark should run from within the wasm.")
	flag.IntVar(&cfg.HeapPages, "heap-pages", 4096, "Cache heap allocation pages.")
	flag.IntVar(&cfg.DbCache, "db-cache", 1024, "Limit the memory the database cache can use.")
	flag.StringVar(&cfg.GC, "gc", "", "GC flag used for building the runtime.")
	flag.StringVar(&cfg.TinyGoVersion, "tinygoversion", "", "TinyGO version used for building the runtime.")
	flag.StringVar(&cfg.Target, "target", "", "Target used for building the runtime.")
	flag.BoolVar(&cfg.GenerateWeightFiles, "generate-weight-files", true, "Whether to generate weight files.")
	flag.Parse()
	return *cfg
}
