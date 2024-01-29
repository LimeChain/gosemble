package types

import sc "github.com/LimeChain/goscale"

type BlockHashCount struct {
	sc.U32
}

func (rv BlockHashCount) Docs() string {
	return "Maximum number of block number to block hash mappings to keep (oldest pruned first)."
}
