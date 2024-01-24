package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/runtime"
	wazero "github.com/ChainSafe/gossamer/lib/runtime/wazero"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

// TODO: move once CLI is merged
type BenchmarkOverheadParams struct {
	Warmup int
	Repeat int
}

func BenchmarkBlockExecutionWeight(t *testing.B) {
	// TODO: Extract in benchmarking CLI logic once it is merged
	config := BenchmarkOverheadParams{
		Warmup: 10,
		Repeat: 100,
	}
	rt, storage := newBenchmarkingRuntime(t)

	// Generate Genesis Block
	generateGenesisBlock(t, rt)

	// Builds the Block
	block := buildBlock(t, rt, false)

	// Measure the block
	result := measureBlock(t, rt, storage, config, block)

	// TODO: Process result once benchmarking CLI is merged
	fmt.Println(result)
}

func measureBlock(t *testing.B, rt *wazero.Instance, storage *runtime.Storage, config BenchmarkOverheadParams, block gossamertypes.Block) []int {
	encodedBlock, err := scale.Marshal(block)
	assert.Nil(t, err)

	// Store the DB Snapshot with only Genesis Block
	(*storage).DbStoreSnapshot()

	fmt.Println(fmt.Sprintf("Running [%d] warmups", config.Warmup))
	for i := 0; i < config.Warmup; i++ {
		(*storage).DbRestoreSnapshot()
		_, err := rt.Exec("Core_execute_block", encodedBlock)
		assert.NoError(t, err)
	}

	var result []int
	fmt.Println(fmt.Sprintf("Executing block [%d] times", config.Repeat))
	for i := 0; i < config.Repeat; i++ {
		(*storage).DbRestoreSnapshot()

		start := time.Now()
		_, err := rt.Exec("Core_execute_block", encodedBlock)
		assert.NoError(t, err)
		end := time.Now()
		result = append(result, end.Nanosecond()-start.Nanosecond())
	}

	return result
}

func buildBlock(t assert.TestingT, rt *wazero.Instance, withExtrinsic bool) gossamertypes.Block {
	metadata := runtimeMetadata(t, rt)

	idata := gossamertypes.NewInherentData()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(dateTime.UnixMilli()))
	assert.NoError(t, err)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	expectedExtrinsicBytes := timestampExtrinsicBytes(t, metadata, uint64(dateTime.UnixMilli()))

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)
	assert.NotNil(t, inherentExt)

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{inherentExt[0]})

	totalInherents, err := sc.DecodeCompact[sc.U128](buffer)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	actualExtrinsic := inherentExt[1:]
	assert.Equal(t, expectedExtrinsicBytes, actualExtrinsic)

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	assert.Nil(t, err)

	parentHash := common.BytesToHash(types.Blake2bHash69().Bytes())
	storageRoot := common.MustHexToHash("0x488445c43e6bfb9f01d25df40b5e1a6b18fd0f45dd89e7e0d5b84d43bd2285eb")

	header := gossamertypes.NewHeader(parentHash, storageRoot, extrinsicsRoot, uint(blockNumber), gossamertypes.NewDigest())

	return gossamertypes.Block{
		Header: *header,
		Body:   gossamertypes.BytesArrayToExtrinsics(exts),
	}
}

func generateGenesisBlock(t assert.TestingT, rt *wazero.Instance) {
	genesisConfig := []byte("{\"system\":{},\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]},\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]},\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1000000000000000000]]},\"transactionPayment\":{\"multiplier\":\"2\"}}")

	// Generate genesis block
	res, err := rt.Exec("GenesisBuilder_build_config", sc.BytesToSequenceU8(genesisConfig).Bytes())
	assert.NoError(t, err)
	assert.Equal(t, []byte{0}, res)
}
