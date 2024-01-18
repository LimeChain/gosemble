package main

import (
	"bytes"
	"testing"
	"time"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/benchmarking"
	benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func BenchmarkTimestampSet(b *testing.B) {
	linear1, err := benchmarking.NewLinear(0, 1000)
	assert.NoError(b, err)

	linear2, err := benchmarking.NewLinear(0, 1000)
	assert.NoError(b, err)

	benchmarking.Run(b, "timestamp_set", func(i *benchmarking.Instance) *benchmarkingtypes.BenchmarkResult {
		// arrange
		(*i.Storage()).Put(append(keyTimestampHash, keyTimestampNowHash...), sc.U64(0).Bytes())
		(*i.Storage()).DbWhitelistKey(string(append(keyTimestampHash, keyTimestampDidUpdate...)))

		now := uint64(time.Now().UnixMilli())
		now += uint64(linear1.Value())
		now += uint64(linear2.Value())

		// act
		benchmarkResult, err := i.ExecuteExtrinsic(
			"Timestamp.set",
			sc.NewOption[primitives.RawOrigin](primitives.NewRawOriginNone()),
			nil,
			ctypes.NewUCompactFromUInt(now),
		)

		// assert
		assert.NoError(b, err)

		nowStorageValue, err := sc.DecodeU64(bytes.NewBuffer((*i.Storage()).Get(append(keyTimestampHash, keyTimestampNowHash...))))
		assert.NoError(b, err)
		assert.Equal(b, sc.U64(now), nowStorageValue)

		didUpdateStorageValue, err := sc.DecodeBool(bytes.NewBuffer((*i.Storage()).Get(append(keyTimestampHash, keyTimestampDidUpdate...))))
		assert.NoError(b, err)
		assert.Equal(b, sc.Bool(true), didUpdateStorageValue)

		return benchmarkResult
	}, linear1, linear2)
}
