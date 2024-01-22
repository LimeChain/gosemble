package benchmarking

import (
	"bytes"
	"testing"
	"time"

	"github.com/ChainSafe/gossamer/lib/common"
	wazero_runtime "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/lib/trie"
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var (
	keyTimestampHash, _    = common.Twox128Hash([]byte("Timestamp"))
	keyTimestampNowHash, _ = common.Twox128Hash([]byte("Now"))
)

func TestInstance(t *testing.T) {
	b := &testing.B{}
	runtime := wazero_runtime.NewBenchInstanceWithTrie(b, "../build/runtime.wasm", trie.NewEmptyTrie())
	defer runtime.Stop()

	instance, err := newBenchmarkingInstance(runtime, 1)
	assert.NoError(t, err)

	(*instance.Storage()).Put(append(keyTimestampHash, keyTimestampNowHash...), sc.U64(0).Bytes())

	nowStorageValue, err := sc.DecodeU64(bytes.NewBuffer((*instance.Storage()).Get(append(keyTimestampHash, keyTimestampNowHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.U64(0), nowStorageValue)

	now := uint64(time.Now().UnixMilli())

	benchmarkResult, err := instance.ExecuteExtrinsic(
		"Timestamp.set",
		sc.NewOption[primitives.RawOrigin](primitives.NewRawOriginNone()),
		nil,
		ctypes.NewUCompactFromUInt(now),
	)
	assert.NoError(t, err)
	assert.NotZero(t, benchmarkResult.ExtrinsicTime)
	assert.NotZero(t, benchmarkResult.Reads)
	assert.NotZero(t, benchmarkResult.Writes)

	nowStorageValue, err = sc.DecodeU64(bytes.NewBuffer((*instance.Storage()).Get(append(keyTimestampHash, keyTimestampNowHash...))))
	assert.NoError(b, err)
	assert.Equal(b, sc.U64(now), nowStorageValue)
}
