package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	BlockHashCount sc.U32
	Version        types.RuntimeVersion
}

func NewConfig(blockHashCount sc.U32, version types.RuntimeVersion) *Config {
	return &Config{
		blockHashCount,
		version,
	}
}
