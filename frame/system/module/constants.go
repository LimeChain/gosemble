package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type consts struct {
	BlockHashCount sc.U32
	Version        types.RuntimeVersion
}

func newConstants(blockHashCount sc.U32, version types.RuntimeVersion) *consts {
	return &consts{
		blockHashCount,
		version,
	}
}
