package system

import (
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

var ZeroAddress = types.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

type CheckNonZeroAddress types.Address32

func (a CheckNonZeroAddress) AdditionalSigned() (ok sc.Empty, err types.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (who CheckNonZeroAddress) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	// TODO:
	// Not sure when this is possible.
	// Checks signed transactions but will fail
	// before this check if the address is all zeros.
	if !reflect.DeepEqual(who, ZeroAddress) {
		ok = types.DefaultValidTransaction()
		return ok, err
	}

	err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.BadSignerError))

	return ok, err
}

func (a CheckNonZeroAddress) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = a.Validate(who, call, info, length)
	return ok, err
}
