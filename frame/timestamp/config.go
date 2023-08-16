package timestamp

import (
	"github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
)

type Config struct {
	OnTimestampSet hooks.OnTimestampSet[goscale.U64]
	MinimumPeriod  goscale.U64
}

func NewConfig(onTsSet hooks.OnTimestampSet[goscale.U64], minimumPeriod goscale.U64) *Config {
	return &Config{
		OnTimestampSet: onTsSet,
		MinimumPeriod:  minimumPeriod,
	}
}
