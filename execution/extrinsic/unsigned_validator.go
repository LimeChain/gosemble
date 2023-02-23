package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type UnsignedValidatorForChecked struct{}

// Validate the call right before dispatch.
//
// This method should be used to prevent transactions already in the pool
// (i.e. passing [`validate_unsigned`](Self::validate_unsigned)) from being included in blocks
// in case they became invalid since being added to the pool.
//
// By default it's a good idea to call [`validate_unsigned`](Self::validate_unsigned) from
// within this function again to make sure we never include an invalid transaction. Otherwise
// the implementation of the call or this method will need to provide proper validation to
// ensure that the transaction is valid.
//
// Changes made to storage *WILL* be persisted if the call returns `Ok`.
func (v UnsignedValidatorForChecked) PreDispatch(call *types.Call) (ok sc.Empty, err types.TransactionValidityError) {
	_, err = v.ValidateUnsigned(types.NewTransactionSource(types.InBlock), call) // .map(|_| ()).map_err(Into::into)
	return ok, err
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
// Inherent call is not validated as unsigned
func (v UnsignedValidatorForChecked) ValidateUnsigned(_source types.TransactionSource, call *types.Call) (ok types.ValidTransaction, err types.TransactionValidityError) {
	noUnsignedValidatorError := types.NewTransactionValidityError(types.NewUnknownTransaction(types.NoUnsignedValidatorError))
	// TODO: Add more modules
	switch call.CallIndex.ModuleIndex {
	case system.Module.Index:
		switch call.CallIndex.FunctionIndex {
		case system.Module.Functions["remark"].Index:
			ok = types.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}

	case timestamp.Module.Index:
		switch call.CallIndex.FunctionIndex {
		case timestamp.Module.Functions["set"].Index:
			ok = types.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}

	default:
		log.Critical("no module found")
	}

	return ok, err
}
