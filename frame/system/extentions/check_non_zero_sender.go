package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonZeroAddress types.Address32

func (who CheckNonZeroAddress) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	// TODO:
	// Not sure when this is possible.
	// Checks signed transactions but will fail
	// before this check if the address is all zeros.
	for _, v := range who.Bytes() {
		if v != 0 {
			ok = types.DefaultValidTransaction()
			return ok, err
		}
	}

	err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadSignerError))

	return ok, err
}
