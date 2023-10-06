package types

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckedExtrinsic interface {
	Apply(validator UnsignedValidator, info *primitives.DispatchInfo, length sc.Compact) (primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo], primitives.TransactionValidityError)
	Extra() primitives.SignedExtra
	Function() primitives.Call
	Signed() sc.Option[primitives.Address32]
	Validate(validator UnsignedValidator, source primitives.TransactionSource, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError)
}
