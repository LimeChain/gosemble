package system

import sc "github.com/LimeChain/goscale"

// System module errors.
const (
	ErrorInvalidSpecName sc.U8 = iota
	ErrorSpecVersionNeedsToIncrease
	ErrorFailedToExtractRuntimeVersion
	ErrorNonDefaultComposite
	ErrorNonZeroRefCount
	ErrorCallFiltered
)
