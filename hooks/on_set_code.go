package hooks

import (
	sc "github.com/LimeChain/goscale"
)

// Do something when we should be setting the code.
type OnSetCode interface {
	// Set the code to the given blob.
	SetCode(sc.Sequence[sc.U8]) error
}
