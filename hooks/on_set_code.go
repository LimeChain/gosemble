package hooks

import (
	sc "github.com/LimeChain/goscale"
)

// OnSetCode does something when we should be setting the code.
type OnSetCode interface {
	// SetCode sets the code to the given blob.
	SetCode(sc.Sequence[sc.U8]) error
}
