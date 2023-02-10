package types

import "github.com/LimeChain/goscale"

type IsFatalError interface {
	error
	IsFatalError() goscale.Bool
}
