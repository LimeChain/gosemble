package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockHashCount sc.U64
	BlockWeights   types.BlockWeights
	BlockLength    types.BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount sc.U64, blockWeights types.BlockWeights, blockLength types.BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *consts {
	return &consts{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
