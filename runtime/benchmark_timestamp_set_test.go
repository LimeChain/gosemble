package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTimestampSet(b *testing.B) {
	benchmarkInstance(b, TimestampSet)
}

func TimestampSet(b *testing.B, rt *wazero_runtime.Instance, storage *runtime.Storage, metadata *ctypes.Metadata, args ...interface{}) (ctypes.Call, []byte) {
	// Setup the input params
	now := uint64(time.Now().UnixMilli())
	call, err := ctypes.NewCall(metadata, "Timestamp.set", ctypes.NewUCompactFromUInt(now))
	assert.NoError(b, err)

	benchmarkConfig := newExtrinsicCall(b, types.NewRawOriginNone(), call)

	// Setup the state
	(*storage).Put(append(keyTimestampHash, keyTimestampNowHash...), sc.U64(0).Bytes())
	assert.NoError(b, err)

	// Whitelist the keys
	(*storage).DbWhitelistKey(string(append(keyTimestampHash, keyTimestampDidUpdateHash...)))

	// Execute the call
	res, err := rt.Exec("Benchmark_run", benchmarkConfig.Bytes())
	assert.NoError(b, err)

	// Validate the result
	nowStorageValue, err := sc.DecodeU64(bytes.NewBuffer((*storage).Get(append(keyTimestampHash, keyTimestampNowHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.U64(now), nowStorageValue)

	didUpdateStorageValue, err := sc.DecodeBool(bytes.NewBuffer((*storage).Get(append(keyTimestampHash, keyTimestampDidUpdateHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.Bool(true), didUpdateStorageValue)

	return call, res
}
