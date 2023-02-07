package types

import (
	sc "github.com/LimeChain/goscale"
)

type ApplyExtrinsicResult sc.VaryingData

func NewApplyExtrinsicResult(value sc.Encodable) ApplyExtrinsicResult {
	// DispatchOutcome 					= 0 Outcome of dispatching the extrinsic.
	// TransactionValidityError = 1 Possible errors while checking the validity of a transaction.
	return ApplyExtrinsicResult(sc.NewVaryingData(value))
}
