package benchmarking

import (
	"fmt"
	"reflect"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	invalidTransactionExhaustsResourcesErr      = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionExhaustsResources())
	dispatchOutcome, _                          = primitives.NewDispatchOutcome(nil)
	applyExtrinsicResultOutcome, _              = primitives.NewApplyExtrinsicResult(dispatchOutcome)
	applyExtrinsicResultExhaustsResourcesErr, _ = primitives.NewApplyExtrinsicResult(invalidTransactionExhaustsResourcesErr.(primitives.TransactionValidityError))
)

type BlockBuilder struct {
	instance     *Instance
	inherentData []byte
	extrinsics   [][]byte
}

func NewBlockBuilder(instance *Instance, inherentData []byte) BlockBuilder {
	return BlockBuilder{
		instance:     instance,
		inherentData: inherentData,
	}
}

// StartSimulation begins a block simulations.
// Create a snapshot of the DB.
func (bb *BlockBuilder) StartSimulation(blockNumber uint) error {
	(*bb.instance.storage).DbStoreSnapshot()
	bb.extrinsics = nil

	header := gossamertypes.NewHeader(parentHash, common.Hash{}, common.Hash{}, blockNumber, gossamertypes.NewDigest())
	encodedHeader, err := scale.Marshal(*header)
	if err != nil {
		return err
	}

	_, err = bb.instance.runtime.Exec("Core_initialize_block", encodedHeader)
	if err != nil {
		return err
	}

	return nil
}

// ApplyInherentExtrinsics converts the inherent data to extrinsics and adds them to the simulated block.
func (bb *BlockBuilder) ApplyInherentExtrinsics() error {
	inherentExt, err := bb.instance.runtime.InherentExtrinsics(bb.inherentData)
	if err != nil {
		return err
	}

	var exts [][]byte
	err = scale.Unmarshal(inherentExt, &exts)
	if err != nil {
		return err
	}

	// Apply the Timestamp Extrinsic
	_, err = bb.instance.runtime.Exec("BlockBuilder_apply_extrinsic", inherentExt[1:])
	if err != nil {
		return err
	}

	bb.extrinsics = append(bb.extrinsics, exts...)

	return nil
}

// AddExtrinsic adds an extrinsic to the simulated block.
// Returns true of the applied extrinsic has reached the ExhaustsResource to the simulated block.
func (bb *BlockBuilder) AddExtrinsic(extrinsic []byte) (bool, error) {
	result, err := bb.instance.runtime.Exec("BlockBuilder_apply_extrinsic", extrinsic)
	if err != nil {
		return false, err
	}

	if reflect.DeepEqual(applyExtrinsicResultOutcome.Bytes(), result) {
		extrinsic := extrinsic[2:] // Exclude len of Extrinsic as it is added later on
		bb.extrinsics = append(bb.extrinsics, extrinsic)
	} else if reflect.DeepEqual(applyExtrinsicResultExhaustsResourcesErr.Bytes(), result) {
		fmt.Println("Reached Exhausts Resources")
		return true, nil
	} else {
		return false, fmt.Errorf("invalid error during block extrinsic weight execution [Apply Extrinsic] - [%v]", result)
	}

	return false, nil
}

// FinishSimulation finalizes the block and returns it.
// Restores the DB to the initial snapshot from StartSimulation
func (bb *BlockBuilder) FinishSimulation() (gossamertypes.Block, error) {
	fmt.Println(fmt.Sprintf("Extrinsics per block [%d]", len(bb.extrinsics)))

	// Finalize the Block to get the result Header
	res, err := bb.instance.runtime.Exec("BlockBuilder_finalize_block", []byte{})
	if err != nil {
		return gossamertypes.Block{}, err
	}

	resultHeader := gossamertypes.NewEmptyHeader()
	err = scale.Unmarshal(res, resultHeader)
	if err != nil {
		return gossamertypes.Block{}, err
	}

	// Restore to Initial Genesis Block
	(*bb.instance.storage).DbRestoreSnapshot()

	return gossamertypes.Block{
		Header: *resultHeader,
		Body:   gossamertypes.BytesArrayToExtrinsics(bb.extrinsics),
	}, nil
}
