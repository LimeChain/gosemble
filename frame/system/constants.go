package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockHashCount sc.U32
	BlockWeights   BlockWeights
	BlockLength    BlockLength
	DbWeight       types.RuntimeDbWeight
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount sc.U32, blockWeights BlockWeights, blockLength BlockLength, dbWeight types.RuntimeDbWeight, version types.RuntimeVersion) *consts {
	return &consts{
		blockHashCount,
		blockWeights,
		blockLength,
		dbWeight,
		version,
	}
}
