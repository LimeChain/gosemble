package types

import sc "github.com/LimeChain/goscale"

type FatalError interface {
	sc.Encodable
	IsFatal() sc.Bool
	Error() string
}
