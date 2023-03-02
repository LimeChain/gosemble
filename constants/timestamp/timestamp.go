package timestamp

import sc "github.com/LimeChain/goscale"

const (
	ModuleIndex      = sc.U8(3)
	FunctionSetIndex = 0
)

const (
	MaxTimestampDriftMillis = 30 * 1_000 // 30 Seconds
	MinimumPeriod           = 1 * 1_000  // 1 second
)

var (
	InherentIdentifier = [8]byte{'t', 'i', 'm', 's', 't', 'a', 'p', '0'}
)
