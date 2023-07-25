package module

import (
	sc "github.com/LimeChain/goscale"
)

type Config struct {
	MinimumPeriod              sc.U64
	MaxAuthorities             sc.U32
	AllowMultipleBlocksPerSlot bool
}

func NewConfig(minimumPeriod sc.U64, maxAuthorities sc.U32, allowMultipleBlocksPerSlot bool) *Config {
	return &Config{
		MinimumPeriod:              minimumPeriod,
		MaxAuthorities:             maxAuthorities,
		AllowMultipleBlocksPerSlot: allowMultipleBlocksPerSlot,
	}
}
