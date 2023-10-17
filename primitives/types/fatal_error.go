package types

import "github.com/LimeChain/goscale"

type FatalError interface {
	goscale.Encodable
	IsFatal() goscale.Bool
	Error() string
}
