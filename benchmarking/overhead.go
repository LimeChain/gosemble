package benchmarking

import (
	"fmt"
	"math"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/pkg/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

var (
	signer = signature.TestKeyringPairAlice
)

type OverheadConfig struct {
	Warmup         int
	Repeat         int
	MaxExtPerBlock int
}

func benchBlock(b *testing.B, instance *Instance, config OverheadConfig) {
	err := instance.BuildGenesisConfig()
	if err != nil {
		b.Fatal(err)
	}

	// Build the block
	block, err := buildBlock(instance, false, config.MaxExtPerBlock)
	if err != nil {
		b.Fatalf("failed to build empty block: [%v]", err)
	}

	// Measure the block
	results := measureBlock(b, instance, config, block)

	// Create the stats
	stats, err := NewOverheadStats(results)
	if err != nil {
		panic(err)
	}
	fmt.Println("result stats")
	fmt.Println(stats)
}

func benchExtrinsic(b *testing.B, instance *Instance, config OverheadConfig) {
	err := instance.BuildGenesisConfig()
	if err != nil {
		b.Fatal(err)
	}

	// Build an Empty block
	baseEmptyBlock, err := buildBlock(instance, false, config.MaxExtPerBlock)
	if err != nil {
		b.Fatalf("failed to build empty block: [%v]", err)
	}

	// Measure the Empty block
	baseResults := measureBlock(b, instance, config, baseEmptyBlock)

	// Create the stats
	baseStats, err := NewOverheadStats(baseResults)
	if err != nil {
		panic(err)
	}

	// Build a Block with maximum possible Extrinsics
	block, err := buildBlock(instance, true, config.MaxExtPerBlock)
	if err != nil {
		b.Fatalf("failed to build block with extrinsics: [%v]", err)
	}
	totalExtrinsics := len(block.Body)

	// Measure the block
	extrinsicResults := measureBlock(b, instance, config, block)

	for i, extrinsicResult := range extrinsicResults {
		extrinsicResult = math.Max(extrinsicResult-baseStats.Mean, 0)

		extrinsicResult = math.Ceil(extrinsicResult / float64(totalExtrinsics))

		extrinsicResults[i] = extrinsicResult
	}

	stats, err := NewOverheadStats(extrinsicResults)
	if err != nil {
		panic(err)
	}
	fmt.Println("result stats")
	fmt.Println(stats)
}

func buildBlock(instance *Instance, hasExtrinsics bool, maxExtrinsicsPerBlock int) (gossamertypes.Block, error) {
	inherentData, err := timestampInherentData(dateTime)
	if err != nil {
		return gossamertypes.Block{}, fmt.Errorf("failed to create inherent data: [%v]", err)
	}

	blockBuilder := NewBlockBuilder(instance, inherentData)

	err = blockBuilder.StartSimulation(blockNumber)
	if err != nil {
		return gossamertypes.Block{}, fmt.Errorf("failed to start simulation: [%v]", err)
	}

	err = blockBuilder.ApplyInherentExtrinsics()
	if err != nil {
		return gossamertypes.Block{}, fmt.Errorf("failed to create inherent extrinsics: [%v]", err)
	}

	if hasExtrinsics {
		fmt.Println("Building block, this takes some time...")

		for i := 0; i < maxExtrinsicsPerBlock; i++ {
			signatureOptions := ctypes.SignatureOptions{
				BlockHash:          ctypes.Hash(parentHash),
				Era:                ctypes.ExtrinsicEra{IsImmortalEra: true},
				GenesisHash:        ctypes.Hash(parentHash),
				Nonce:              ctypes.NewUCompactFromUInt(uint64(i)),
				SpecVersion:        ctypes.U32(instance.version.SpecVersion),
				Tip:                ctypes.NewUCompactFromUInt(0),
				TransactionVersion: ctypes.U32(instance.version.TransactionVersion),
			}

			extrinsic, err := instance.newSignedExtrinsic(signer, signatureOptions, "System.remark", []byte{})
			if err != nil {
				return gossamertypes.Block{}, fmt.Errorf("failed to sign extrinsic: [%v]", err)
			}

			limitReached, err := blockBuilder.AddExtrinsic(extrinsic)
			if err != nil {
				return gossamertypes.Block{}, fmt.Errorf("failed to add extrinsic: [%v]", err)
			}

			if limitReached {
				break
			}
		}
	}

	block, err := blockBuilder.FinishSimulation()
	if err != nil {
		return gossamertypes.Block{}, fmt.Errorf("failed to finish simulation: [%v]", err)
	}

	return block, nil
}

func measureBlock(t *testing.B, instance *Instance, config OverheadConfig, block gossamertypes.Block) []float64 {
	encodedBlock, err := scale.Marshal(block)
	assert.Nil(t, err)

	// Store the DB Snapshot with only Genesis Block
	(*instance.storage).DbStoreSnapshot()

	fmt.Println(fmt.Sprintf("Running [%d] warmups", config.Warmup))
	for i := 0; i < config.Warmup; i++ {
		(*instance.storage).DbRestoreSnapshot()
		_, err := instance.runtime.Exec("Core_execute_block", encodedBlock)
		assert.NoError(t, err)
	}

	var results []float64
	fmt.Println(fmt.Sprintf("Executing block [%d] times", config.Repeat))
	for i := 0; i < config.Repeat; i++ {
		(*instance.storage).DbRestoreSnapshot()

		start := time.Now().UnixNano()
		_, err := instance.runtime.Exec("Core_execute_block", encodedBlock)
		assert.NoError(t, err)
		end := time.Now().UnixNano()
		results = append(results, float64(end-start))
	}

	return results
}
