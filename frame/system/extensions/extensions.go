package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TODO:
// we need to have a way for configuring any additional
// checks that need to be performed
//
// For example:
// CheckNonZeroSender
// CheckSpecVersion
// CheckTxVersion
// CheckGenesis
// CheckEra
// CheckNonce
// CheckWeight
// ChargeAssetTxPayment
//
// currently those checks are explicit, but
// depending on the configuration
// we could use reflection instead

type Extra types.SignedExtra

func (e Extra) AdditionalSigned() (ok types.AdditionalSigned, err types.TransactionValidityError) {
	ok = types.AdditionalSigned{} // FormatVersion: types.ExtrinsicFormatVersion

	specVersion, err := CheckSpecVersion{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.SpecVersion = specVersion

	transactionVersion, err := CheckTxVersion{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.TransactionVersion = transactionVersion

	genesishash, err := CheckGenesis{}.AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.GenesisHash = genesishash

	blockHash, err := CheckMortality(e.Era).AdditionalSigned()
	if err != nil {
		return ok, err
	}
	ok.BlockHash = blockHash

	return ok, err
}

// Information on a transaction's validity and, if valid, on how it relates to other transactions.
func (e Extra) Validate(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	valid := types.DefaultValidTransaction()

	ok, err = CheckNonZeroAddress(*who).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	// TODO: CheckSpecVersion<Runtime>
	// TODO: CheckTxVersion<Runtime>
	// TODO: CheckGenesis<Runtime>

	ok, err = CheckMortality(e.Era).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	ok, err = CheckNonce(e.Nonce).Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	ok, err = CheckWeight{}.Validate(who, call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	// TODO: ChargeAssetTxPayment<Runtime>

	return valid, err
}

func (e Extra) ValidateUnsigned(call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	valid := types.DefaultValidTransaction()

	ok, err = CheckWeight{}.ValidateUnsigned(call, info, length)
	if err != nil {
		return ok, err
	}
	valid = valid.CombineWith(ok)

	return valid, err
}

// Do any pre-flight stuff for a signed transaction.
//
// Make sure to perform the same checks as in [`Validate`].
func (e Extra) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = CheckNonZeroAddress(*who).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	// TODO: CheckSpecVersion<Runtime>
	// TODO: CheckTxVersion<Runtime>
	// TODO: CheckGenesis<Runtime>

	_, err = CheckMortality(e.Era).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	_, err = CheckNonce(e.Nonce).PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	_, err = CheckWeight{}.PreDispatch(who, call, info, length)
	if err != nil {
		return ok, err
	}

	// TODO: ChargeAssetTxPayment<Runtime>

	return ok, err
}

func (e Extra) PreDispatchUnsigned(call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = CheckWeight{}.PreDispatchUnsigned(call, info, length)
	return ok, err
}

func (e Extra) PostDispatch(pre sc.Option[types.Pre], info *types.DispatchInfo, postInfo *types.PostDispatchInfo, length sc.Compact, result *types.DispatchResult) (ok types.Pre, err types.TransactionValidityError) {
	_, err = CheckWeight{}.PostDispatch(pre, info, postInfo, length, result)
	return ok, err
}
