package system

import (
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockHashCount types.BlockHashCount
	BlockWeights   types.BlockWeights
	BlockLength    types.BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount types.BlockHashCount, blockWeights types.BlockWeights, blockLength types.BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *consts {
	return &consts{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
