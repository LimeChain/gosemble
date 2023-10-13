package types

import sc "github.com/LimeChain/goscale"

type CheckedExtrinsic interface {
	Apply(validator UnsignedValidator, info *DispatchInfo, length sc.Compact) (DispatchResultWithPostInfo[PostDispatchInfo], TransactionValidityError)
	Function() Call
	Validate(validator UnsignedValidator, source TransactionSource, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError)
}
