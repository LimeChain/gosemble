package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	balances "github.com/LimeChain/gosemble/frame/balances/module"
	system "github.com/LimeChain/gosemble/frame/system/module"
	timestamp "github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
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
func (v UnsignedValidatorForChecked) PreDispatch(call *types.Call) (ok sc.Empty, err primitives.TransactionValidityError) {
	_, err = v.ValidateUnsigned(primitives.NewTransactionSourceInBlock(), call) // .map(|_| ()).map_err(Into::into)
	return ok, err
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
// Inherent call is not validated as unsigned
func (v UnsignedValidatorForChecked) ValidateUnsigned(_source primitives.TransactionSource, call *types.Call) (ok primitives.ValidTransaction, err primitives.TransactionValidityError) {
	noUnsignedValidatorError := primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	// TODO: Add more modules
	switch call.CallIndex.ModuleIndex {
	case system.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case system.Module.Remark.Index():
			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}

	case timestamp.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case timestamp.Module.Set.Index():
			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}
	case balances.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case balances.Module.Transfer.Index(),
			balances.Module.SetBalance.Index(),
			balances.Module.ForceTransfer.Index(),
			balances.Module.TransferKeepAlive.Index(),
			balances.Module.TransferAll.Index(),
			balances.Module.ForceFree.Index():

			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}
	default:
		log.Critical("no module found")
	}

	return ok, err
}
