//go:build !benchmarks

package benchmark

func TimeNow() int64 {
	panic("not implemented")
}

func DbWhitelistKey(key []byte) {
	panic("not implemented")
}

func DbResetTracker() {
	panic("not implemented")
}

func DbStartTracker() {
	panic("not implemented")
}

func DbStopTracker() {
	panic("not implemented")
}

func DbWipe() {
	panic("not implemented")
}

func DbCommit() {
	panic("not implemented")
}

func DbReadCount() int32 {
	panic("not implemented")
}

func DbWriteCount() int32 {
	panic("not implemented")
}
