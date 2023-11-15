package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type UnsignedValidatorForChecked struct {
	runtimeExtrinsic RuntimeExtrinsic
}

func NewUnsignedValidatorForChecked(extrinsic RuntimeExtrinsic) primitives.UnsignedValidator {
	return UnsignedValidatorForChecked{
		runtimeExtrinsic: extrinsic,
	}
}

// PreDispatch validates the dispatch call before execution.
// Inherent call is accepted for being dispatched
func (v UnsignedValidatorForChecked) PreDispatch(call primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	module, ok := v.runtimeExtrinsic.Module(call.ModuleIndex())
	if !ok {
		return sc.Empty{}, nil
	}

	return module.PreDispatch(call)
}

// ValidateUnsigned returns the validity of the dispatch call.
// Inherent call is not validated as unsigned
func (v UnsignedValidatorForChecked) ValidateUnsigned(txSource primitives.TransactionSource, call primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	module, ok := v.runtimeExtrinsic.Module(call.ModuleIndex())
	if !ok {
		// unknownTransactionNoUnsignedValidator, err := primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
		// if err != nil {
		// 	log.Critical(err.Error())
		// }
		// return primitives.ValidTransaction{}, unknownTransactionNoUnsignedValidator
		// todo check all references of this function and check for unexpectedError everywhere its called
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
	}

	return module.ValidateUnsigned(txSource, call)
}
