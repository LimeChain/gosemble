package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
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
func (v UnsignedValidatorForChecked) PreDispatch(call *primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	module, ok := types.Modules[(*call).ModuleIndex()]
	if !ok {
		return sc.Empty{}, nil
	}

	return module.PreDispatch(*call)
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
// Inherent call is not validated as unsigned
func (v UnsignedValidatorForChecked) ValidateUnsigned(_source primitives.TransactionSource, call *primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	module, ok := types.Modules[(*call).ModuleIndex()]
	if !ok {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	}

	return module.ValidateUnsigned(_source, *call)
}
