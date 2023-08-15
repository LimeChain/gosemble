/*
Targets WebAssembly MVP
*/
package main

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/account_nonce"
	"github.com/LimeChain/gosemble/frame/aura"
	bm "github.com/LimeChain/gosemble/frame/balances/module"
	blockbuilder "github.com/LimeChain/gosemble/frame/block_builder"
	"github.com/LimeChain/gosemble/frame/core"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/frame/metadata"
	"github.com/LimeChain/gosemble/frame/offchain_worker"
	"github.com/LimeChain/gosemble/frame/session_keys"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/system/extensions"
	sm "github.com/LimeChain/gosemble/frame/system/module"
	taggedtransactionqueue "github.com/LimeChain/gosemble/frame/tagged_transaction_queue"
	tm "github.com/LimeChain/gosemble/frame/testable/module"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	tpm "github.com/LimeChain/gosemble/frame/transaction_payment/module"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	AuraMaxAuthorites = 100
)

const (
	BalancesMaxLocks    = 50
	BalancesMaxReserves = 50
)

var (
	balancesExistentialDeposit = 1 * constants.Dollar
	BalancesExistentialDeposit = big.NewInt(0).SetUint64(balancesExistentialDeposit)
)

var (
	BlockWeights = system.WithSensibleDefaults(constants.MaximumBlockWeight, constants.NormalDispatchRatio)
	BlockLength  = system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)
)

var (
	OperationalFeeMultiplier                        = sc.U8(5)
	WeightToFee              primitives.WeightToFee = primitives.IdentityFee{}
	LengthToFee              primitives.WeightToFee = primitives.IdentityFee{}
)

const (
	SystemIndex sc.U8 = iota
	TimestampIndex
	AuraIndex
	GrandpaIndex
	BalancesIndex
	TxPaymentsIndex
	TestableIndex = 255
)

// Modules contains all the modules used by the runtime.
var modules = initializeModules()

func initializeModules() map[sc.U8]primitives.Module {
	systemModule := sm.NewSystemModule(SystemIndex,
		sm.NewConfig(constants.BlockHashCount, BlockWeights, BlockLength, constants.RuntimeVersion))

	auraModule := aura.NewModule(AuraIndex,
		aura.NewConfig(
			primitives.PublicKeySr25519,
			timestamp.MinimumPeriod,
			AuraMaxAuthorites,
			false,
			systemModule.Storage.Digest.Get))

	timestampModule := tsm.NewModule(TimestampIndex,
		tsm.NewConfig(auraModule, timestamp.MinimumPeriod))

	grandpaModule := grandpa.NewModule(GrandpaIndex)

	balancesModule := bm.NewBalancesModule(BalancesIndex,
		bm.NewConfig(BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule))

	tpmModule := tpm.NewTransactionPaymentModule(TxPaymentsIndex, tpm.NewConfig(OperationalFeeMultiplier, WeightToFee, LengthToFee, BlockWeights))
	testableModule := tm.NewTestingModule(TestableIndex)

	return map[sc.U8]primitives.Module{
		SystemIndex:     systemModule,
		TimestampIndex:  timestampModule,
		AuraIndex:       auraModule,
		GrandpaIndex:    grandpaModule,
		BalancesIndex:   balancesModule,
		TxPaymentsIndex: tpmModule,
		TestableIndex:   testableModule,
	}
}

func newExecutiveModule() executive.Module {
	return executive.New(modules[SystemIndex].(sm.SystemModule), extrinsic.New(modules), modules[AuraIndex].(aura.Module))
}

func newModuleDecoder() types.ModuleDecoder {
	return types.NewModuleDecoder(modules, newSignedExtra())
}

func newSignedExtra() primitives.SignedExtra {
	systeModule := modules[SystemIndex].(sm.SystemModule)
	balancesModule := modules[BalancesIndex].(bm.BalancesModule)
	txPaymentModule := modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)

	checkMortality := extensions.NewCheckMortality(systeModule)
	checkNonce := extensions.NewCheckNonce(systeModule)
	chargeTxPayment := transaction_payment.NewChargeTransactionPayment(systeModule, txPaymentModule, balancesModule)

	extras := []primitives.SignedExtension{
		extensions.NewCheckNonZeroAddress(),
		extensions.NewCheckSpecVersion(systeModule),
		extensions.NewCheckTxVersion(systeModule),
		extensions.NewCheckGenesis(systeModule),
		&checkMortality,
		&checkNonce,
		extensions.NewCheckWeight(systeModule),
		&chargeTxPayment,
	}
	return primitives.NewSignedExtra(extras)
}

//go:export Core_version
func CoreVersion(_ int32, _ int32) int64 {
	return core.
		New(newExecutiveModule(), newModuleDecoder()).
		Version()
}

//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) int64 {
	core.
		New(newExecutiveModule(), newModuleDecoder()).
		InitializeBlock(dataPtr, dataLen)

	return 0
}

//go:export Core_execute_block
func CoreExecuteBlock(dataPtr int32, dataLen int32) int64 {
	core.
		New(newExecutiveModule(), newModuleDecoder()).
		ExecuteBlock(dataPtr, dataLen)

	return 0
}

//go:export BlockBuilder_apply_extrinsic
func BlockBuilderApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	return blockbuilder.
		New(newExecutiveModule(), newModuleDecoder()).
		ApplyExtrinsic(dataPtr, dataLen)
}

//go:export BlockBuilder_finalize_block
func BlockBuilderFinalizeBlock(_, _ int32) int64 {
	return blockbuilder.
		New(newExecutiveModule(), newModuleDecoder()).
		FinalizeBlock()
}

//go:export BlockBuilder_inherent_extrinsics
func BlockBuilderInherentExtrinsics(dataPtr int32, dataLen int32) int64 {
	return blockbuilder.InherentExtrinsics(dataPtr, dataLen)
}

//go:export BlockBuilder_check_inherents
func BlockBuilderCheckInherents(dataPtr int32, dataLen int32) int64 {
	return blockbuilder.
		New(newExecutiveModule(), newModuleDecoder()).
		CheckInherents(dataPtr, dataLen)
}

//go:export TaggedTransactionQueue_validate_transaction
func TaggedTransactionQueueValidateTransaction(dataPtr int32, dataLen int32) int64 {
	return taggedtransactionqueue.
		New(newExecutiveModule(), newModuleDecoder()).
		ValidateTransaction(dataPtr, dataLen)
}

//go:export AuraApi_slot_duration
func AuraApiSlotDuration(_, _ int32) int64 {
	return modules[AuraIndex].(aura.Module).SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return modules[AuraIndex].(aura.Module).Authorities()
}

//go:export AccountNonceApi_account_nonce
func AccountNonceApiAccountNonce(dataPtr int32, dataLen int32) int64 {
	return account_nonce.New(modules[SystemIndex].(sm.SystemModule)).
		AccountNonce(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_info
func TransactionPaymentApiQueryInfo(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder(), modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)).
		QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder(), modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)).
		QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return transaction_payment.
		New(newModuleDecoder(), modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)).
		QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder(), modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)).
		QueryCallFeeDetails(dataPtr, dataLen)
}

//go:export Metadata_metadata
func Metadata(_, _ int32) int64 {
	return metadata.
		New(modules).
		Metadata()
}

//go:export SessionKeys_generate_session_keys
func SessionKeysGenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	sessions := []primitives.Session{
		modules[AuraIndex].(aura.Module),
		modules[GrandpaIndex].(grandpa.Module),
	}
	return session_keys.New(sessions).GenerateSessionKeys(dataPtr, dataLen)
}

//go:export SessionKeys_decode_session_keys
func SessionKeysDecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	sessions := []primitives.Session{
		modules[AuraIndex].(aura.Module),
		modules[GrandpaIndex].(grandpa.Module),
	}
	return session_keys.New(sessions).DecodeSessionKeys(dataPtr, dataLen)
}

//go:export GrandpaApi_grandpa_authorities
func GrandpaApiAuthorities(_, _ int32) int64 {
	return modules[GrandpaIndex].(grandpa.Module).Authorities()
}

//go:export OffchainWorkerApi_offchain_worker
func OffchainWorkerApiOffchainWorker(dataPtr int32, dataLen int32) int64 {
	offchain_worker.
		New(newExecutiveModule()).
		OffchainWorker(dataPtr, dataLen)

	return 0
}
