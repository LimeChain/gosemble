package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount sc.U32
	BlockWeights   BlockWeights
	BlockLength    BlockLength
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount sc.U32, blockWeights BlockWeights, blockLength BlockLength, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		blockWeights,
		blockLength,
		version,
	}
}
