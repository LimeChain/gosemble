package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Extra types.SignedExtra

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (e Extra) Validate(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	valid := types.DefaultValidTransaction()

	ok, err = CheckNonZeroAddress(*who).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	// TODO: CheckSpecVersion<Runtime>
	// TODO: CheckTxVersion<Runtime>
	// TODO: CheckGenesis<Runtime>

	ok, err = CheckEra(e.Era).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = CheckNonce(e.Nonce).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	ok, err = CheckWeight(e.Weight).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	// TODO: ChargeAssetTxPayment<Runtime>

	return valid, err
}

func (e Extra) ValidateUnsigned(call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	valid := types.DefaultValidTransaction()

	ok, err = CheckWeight(e.Weight).ValidateUnsigned(call, info, length)
	if err != nil {
		return ok, err
	}
	valid.CombineWith(ok)

	return valid, err
}

// Do any pre-flight stuff for a signed transaction.
//
// Make sure to perform the same checks as in [`Validate`].
func PreDispatch(e types.SignedExtra, who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	// TODO:
	return ok, err
}

func PreDispatchUnsigned(call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	// TODO:
	// ValidateUnsigned(call, info, length)
	return ok, err
}

func PostDispatch(pre sc.Option[types.Pre], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, length sc.Compact, result *types.DispatchResult) (ok types.Pre, err types.TransactionValidityError) {
	// TODO:

	switch pre.HasValue {
	case true:
	case false:
	}

	return ok, err
}
