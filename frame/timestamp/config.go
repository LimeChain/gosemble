package timestamp

import (
	"github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Config struct {
	OnTimestampSet hooks.OnTimestampSet[goscale.U64]
	DbWeight       primitives.RuntimeDbWeight
	MinimumPeriod  goscale.U64
}

func NewConfig(onTsSet hooks.OnTimestampSet[goscale.U64], dbWeight primitives.RuntimeDbWeight, minimumPeriod goscale.U64) *Config {
	return &Config{
		OnTimestampSet: onTsSet,
		DbWeight:       dbWeight,
		MinimumPeriod:  minimumPeriod,
	}
}
