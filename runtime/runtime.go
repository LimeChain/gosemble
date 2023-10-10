/*
Targets WebAssembly MVP
*/
package main

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/api/account_nonce"
	apiAura "github.com/LimeChain/gosemble/api/aura"
	blockbuilder "github.com/LimeChain/gosemble/api/block_builder"
	"github.com/LimeChain/gosemble/api/core"
	apiGrandpa "github.com/LimeChain/gosemble/api/grandpa"
	"github.com/LimeChain/gosemble/api/metadata"
	"github.com/LimeChain/gosemble/api/offchain_worker"
	"github.com/LimeChain/gosemble/api/session_keys"
	taggedtransactionqueue "github.com/LimeChain/gosemble/api/tagged_transaction_queue"
	apiTxPayments "github.com/LimeChain/gosemble/api/transaction_payment"
	apiTxPaymentsCall "github.com/LimeChain/gosemble/api/transaction_payment_call"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/frame/balances"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/frame/system"
	sysExtensions "github.com/LimeChain/gosemble/frame/system/extensions"
	tm "github.com/LimeChain/gosemble/frame/testable"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	txExtensions "github.com/LimeChain/gosemble/frame/transaction_payment/extensions"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	AuraMaxAuthorites = 100
)

const (
	BalancesMaxLocks    = 50
	BalancesMaxReserves = 50
)

const (
	TimestampMinimumPeriod = 1 * 1_000 // 1 second
)

// RuntimeVersion contains the version identifiers of the Runtime.
var RuntimeVersion = &primitives.RuntimeVersion{
	SpecName:           sc.Str(constants.SpecName),
	ImplName:           sc.Str(constants.ImplName),
	AuthoringVersion:   sc.U32(constants.AuthoringVersion),
	SpecVersion:        sc.U32(constants.SpecVersion),
	ImplVersion:        sc.U32(constants.ImplVersion),
	TransactionVersion: sc.U32(constants.TransactionVersion),
	StateVersion:       sc.U8(constants.StateVersion),
}

var (
	balancesExistentialDeposit = 1 * constants.Dollar
	BalancesExistentialDeposit = sc.NewU128(balancesExistentialDeposit)
)

var (
	BlockWeights = system.WithSensibleDefaults(constants.MaximumBlockWeight, constants.NormalDispatchRatio)
	BlockLength  = system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)
	DbWeight     = constants.RocksDbWeight
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

func initializeModules() map[sc.U8]types.Module {
	systemModule := system.New(
		SystemIndex,
		system.NewConfig(constants.BlockHashCount, BlockWeights, BlockLength, DbWeight, *RuntimeVersion),
	)

	auraModule := aura.New(
		AuraIndex,
		aura.NewConfig(
			primitives.PublicKeySr25519,
			DbWeight,
			TimestampMinimumPeriod,
			AuraMaxAuthorites,
			false,
			systemModule.StorageDigest,
		),
	)

	timestampModule := timestamp.New(
		TimestampIndex,
		timestamp.NewConfig(auraModule, DbWeight, TimestampMinimumPeriod),
	)

	grandpaModule := grandpa.New(GrandpaIndex)

	balancesModule := balances.New(
		BalancesIndex,
		balances.NewConfig(DbWeight, BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule),
	)

	tpmModule := transaction_payment.New(
		TxPaymentsIndex,
		transaction_payment.NewConfig(OperationalFeeMultiplier, WeightToFee, LengthToFee, BlockWeights),
	)

	testableModule := tm.New(TestableIndex)

	return map[sc.U8]types.Module{
		SystemIndex:     systemModule,
		TimestampIndex:  timestampModule,
		AuraIndex:       auraModule,
		GrandpaIndex:    grandpaModule,
		BalancesIndex:   balancesModule,
		TxPaymentsIndex: tpmModule,
		TestableIndex:   testableModule,
	}
}

func newSignedExtra() primitives.SignedExtra {
	systemModule := modules[SystemIndex].(system.Module)
	balancesModule := modules[BalancesIndex].(balances.Module)
	txPaymentModule := modules[TxPaymentsIndex].(transaction_payment.Module)

	checkMortality := sysExtensions.NewCheckMortality(systemModule)
	checkNonce := sysExtensions.NewCheckNonce(systemModule)
	chargeTxPayment := txExtensions.NewChargeTransactionPayment(systemModule, txPaymentModule, balancesModule)

	extras := []primitives.SignedExtension{
		sysExtensions.NewCheckNonZeroAddress(),
		sysExtensions.NewCheckSpecVersion(systemModule),
		sysExtensions.NewCheckTxVersion(systemModule),
		sysExtensions.NewCheckGenesis(systemModule),
		&checkMortality,
		&checkNonce,
		sysExtensions.NewCheckWeight(systemModule),
		&chargeTxPayment,
	}
	return primitives.NewSignedExtra(extras)
}

func runtimeApi() types.RuntimeApi {
	extra := newSignedExtra()
	decoder := types.NewRuntimeDecoder(modules, extra)
	runtimeExtrinsic := extrinsic.New(modules, extra)
	auraModule := modules[AuraIndex].(aura.Module)
	grandpaModule := modules[GrandpaIndex].(grandpa.Module)
	txPaymentsModule := modules[TxPaymentsIndex].(transaction_payment.Module)

	executiveModule := executive.New(
		modules[SystemIndex].(system.Module),
		runtimeExtrinsic,
		hooks.DefaultOnRuntimeUpgrade{},
	)

	sessions := []primitives.Session{
		auraModule,
		grandpaModule,
	}

	apis := []primitives.ApiModule{
		core.New(executiveModule, decoder, RuntimeVersion),
		blockbuilder.New(runtimeExtrinsic, executiveModule, decoder),
		taggedtransactionqueue.New(executiveModule, decoder),
		metadata.New(runtimeExtrinsic),
		apiAura.New(auraModule),
		apiGrandpa.New(grandpaModule),
		account_nonce.New(modules[SystemIndex].(system.Module)),
		apiTxPayments.New(decoder, txPaymentsModule),
		apiTxPaymentsCall.NewCallApi(decoder, txPaymentsModule),
		session_keys.New(sessions),
		offchain_worker.New(executiveModule),
	}

	runtimeApi := types.NewRuntimeApi(apis)

	RuntimeVersion.SetApis(runtimeApi.Items())

	return runtimeApi
}

//go:export Core_version
func CoreVersion(_ int32, _ int32) int64 {
	return runtimeApi().
		Module(core.ApiModuleName).(core.Core).
		Version()
}

//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) int64 {
	runtimeApi().
		Module(core.ApiModuleName).(core.Core).
		InitializeBlock(dataPtr, dataLen)

	return 0
}

//go:export Core_execute_block
func CoreExecuteBlock(dataPtr int32, dataLen int32) int64 {
	runtimeApi().Module(core.ApiModuleName).(core.Core).
		ExecuteBlock(dataPtr, dataLen)

	return 0
}

//go:export BlockBuilder_apply_extrinsic
func BlockBuilderApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(blockbuilder.ApiModuleName).(blockbuilder.BlockBuilder).
		ApplyExtrinsic(dataPtr, dataLen)
}

//go:export BlockBuilder_finalize_block
func BlockBuilderFinalizeBlock(_, _ int32) int64 {
	return runtimeApi().
		Module(blockbuilder.ApiModuleName).(blockbuilder.BlockBuilder).
		FinalizeBlock()
}

//go:export BlockBuilder_inherent_extrinsics
func BlockBuilderInherentExtrinsics(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(blockbuilder.ApiModuleName).(blockbuilder.BlockBuilder).
		InherentExtrinsics(dataPtr, dataLen)
}

//go:export BlockBuilder_check_inherents
func BlockBuilderCheckInherents(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(blockbuilder.ApiModuleName).(blockbuilder.BlockBuilder).
		CheckInherents(dataPtr, dataLen)
}

//go:export TaggedTransactionQueue_validate_transaction
func TaggedTransactionQueueValidateTransaction(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(taggedtransactionqueue.ApiModuleName).(taggedtransactionqueue.TaggedTransactionQueue).
		ValidateTransaction(dataPtr, dataLen)
}

//go:export AuraApi_slot_duration
func AuraApiSlotDuration(_, _ int32) int64 {
	return runtimeApi().
		Module(apiAura.ApiModuleName).(apiAura.Module).
		SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return runtimeApi().
		Module(apiAura.ApiModuleName).(apiAura.Module).
		Authorities()
}

//go:export AccountNonceApi_account_nonce
func AccountNonceApiAccountNonce(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(account_nonce.ApiModuleName).(account_nonce.Module).
		AccountNonce(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_info
func TransactionPaymentApiQueryInfo(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPayments.ApiModuleName).(apiTxPayments.Module).
		QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPayments.ApiModuleName).(apiTxPayments.Module).
		QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return runtimeApi().
		Module(apiTxPaymentsCall.ApiModuleName).(apiTxPaymentsCall.Module).
		QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPaymentsCall.ApiModuleName).(apiTxPaymentsCall.Module).
		QueryCallFeeDetails(dataPtr, dataLen)
}

//go:export Metadata_metadata
func Metadata(_, _ int32) int64 {
	return runtimeApi().
		Module(metadata.ApiModuleName).(metadata.Module).
		Metadata()
}

//go:export SessionKeys_generate_session_keys
func SessionKeysGenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(session_keys.ApiModuleName).(session_keys.Module).
		GenerateSessionKeys(dataPtr, dataLen)
}

//go:export SessionKeys_decode_session_keys
func SessionKeysDecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(session_keys.ApiModuleName).(session_keys.Module).
		DecodeSessionKeys(dataPtr, dataLen)
}

//go:export GrandpaApi_grandpa_authorities
func GrandpaApiAuthorities(_, _ int32) int64 {
	return runtimeApi().
		Module(apiGrandpa.ApiModuleName).(apiGrandpa.Module).
		Authorities()
}

//go:export OffchainWorkerApi_offchain_worker
func OffchainWorkerApiOffchainWorker(dataPtr int32, dataLen int32) int64 {
	runtimeApi().
		Module(offchain_worker.ApiModuleName).(offchain_worker.Module).
		OffchainWorker(dataPtr, dataLen)

	return 0
}
