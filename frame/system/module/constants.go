package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockHashCount sc.U32
	BlockWeights   system.BlockWeights
	BlockLength    system.BlockLength
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount sc.U32, blockWeights system.BlockWeights, blockLength system.BlockLength, version types.RuntimeVersion) *consts {
	return &consts{
		blockHashCount,
		blockWeights,
		blockLength,
		version,
	}
}
