package aura

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	MinimumPeriod              sc.U64
	MaxAuthorities             sc.U32
	AllowMultipleBlocksPerSlot bool
	SystemDigest               func() primitives.Digest
}

func NewConfig(minimumPeriod sc.U64, maxAuthorities sc.U32, allowMultipleBlocksPerSlot bool, systemDigest func() primitives.Digest) *Config {
	return &Config{
		MinimumPeriod:              minimumPeriod,
		MaxAuthorities:             maxAuthorities,
		AllowMultipleBlocksPerSlot: allowMultipleBlocksPerSlot,
		SystemDigest:               systemDigest,
	}
}
