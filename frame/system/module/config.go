package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount sc.U32
	BlockWeights   system.BlockWeights
	BlockLength    system.BlockLength
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount sc.U32, blockWeights system.BlockWeights, blockLength system.BlockLength, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		blockWeights,
		blockLength,
		version,
	}
}
