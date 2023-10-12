package types

import (
	sc "github.com/LimeChain/goscale"
)

type Block interface {
	sc.Encodable

	Header() Header
	Extrinsics() sc.Sequence[UncheckedExtrinsic]
}
