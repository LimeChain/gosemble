package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	bc "github.com/LimeChain/gosemble/constants/balances"
	system_constants "github.com/LimeChain/gosemble/constants/system"
	tsc "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/execution/types"
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
	// TODO: This should go though all the pallets and call their ValidateUnsigned method
	noUnsignedValidatorError := primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	// TODO: Add more modules
	switch call.CallIndex.ModuleIndex {
	case system_constants.ModuleIndex:
		switch call.CallIndex.FunctionIndex {
		case system_constants.FunctionRemarkIndex:
			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}

	case tsc.ModuleIndex:
		switch call.CallIndex.FunctionIndex {
		case tsc.FunctionSetIndex:
			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}
	case bc.ModuleIndex:
		switch call.CallIndex.FunctionIndex {
		case bc.FunctionTransferIndex,
			bc.FunctionSetBalanceIndex,
			bc.FunctionForceTransferIndex,
			bc.FunctionTransferKeepAliveIndex,
			bc.FunctionTransferAllIndex,
			bc.FunctionForceFreeIndex:

			ok = primitives.DefaultValidTransaction()
		default:
			err = noUnsignedValidatorError
		}
	default:
		log.Critical("no module found")
	}

	return ok, err
}
