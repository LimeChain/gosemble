package system

import (
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockWeights   types.BlockWeights
	BlockLength    types.BlockLength
	BlockHashCount types.BlockHashCount
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount types.BlockHashCount, blockWeights types.BlockWeights, blockLength types.BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *consts {
	return &consts{
		blockWeights,
		blockLength,
		blockHashCount,
		dbWeight,
		version,
	}
}
