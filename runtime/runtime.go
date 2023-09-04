/*
Targets WebAssembly MVP
*/
package main

import (
	"math/big"
	"reflect"

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
	"github.com/LimeChain/gosemble/primitives/log"
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
	BalancesExistentialDeposit = big.NewInt(0).SetUint64(balancesExistentialDeposit)
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

type BlockNumberType = sc.U128

// Modules contains all the modules used by the runtime.
var modules = initializeModules()

func initializeModules() map[sc.U8]types.Module[BlockNumberType] {
	systemModule := system.New[BlockNumberType](
		SystemIndex,
		system.NewConfig(constants.BlockHashCount, BlockWeights, BlockLength, DbWeight, *RuntimeVersion),
	)

	auraModule := aura.New[BlockNumberType](
		AuraIndex,
		aura.NewConfig(
			primitives.PublicKeySr25519,
			DbWeight,
			TimestampMinimumPeriod,
			AuraMaxAuthorites,
			false,
			systemModule.Storage.Digest.Get,
		),
	)

	timestampModule := timestamp.New[BlockNumberType](
		TimestampIndex,
		timestamp.NewConfig(auraModule, DbWeight, TimestampMinimumPeriod),
	)

	grandpaModule := grandpa.New[BlockNumberType](GrandpaIndex)

	balancesModule := balances.New[BlockNumberType](
		BalancesIndex,
		balances.NewConfig(DbWeight, BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule),
	)

	tpmModule := transaction_payment.New[BlockNumberType](
		TxPaymentsIndex,
		transaction_payment.NewConfig(OperationalFeeMultiplier, WeightToFee, LengthToFee, BlockWeights),
	)

	testableModule := tm.New[BlockNumberType](TestableIndex)

	return map[sc.U8]types.Module[BlockNumberType]{
		SystemIndex:     systemModule,
		TimestampIndex:  timestampModule,
		AuraIndex:       auraModule,
		GrandpaIndex:    grandpaModule,
		BalancesIndex:   balancesModule,
		TxPaymentsIndex: tpmModule,
		TestableIndex:   testableModule,
	}
}

func getInstance[T types.Module[BlockNumberType]]() T {
	for _, module := range modules {
		if reflect.TypeOf(module) == reflect.TypeOf(*new(T)) {
			return modules[module.GetIndex()].(T)
		}
	}
	log.Critical("unknown type T for module instance")
	panic("unreachable")
}

func newSignedExtra() primitives.SignedExtra {
	systemModule := getInstance[system.Module[BlockNumberType]]()
	balancesModule := getInstance[balances.Module[BlockNumberType]]()
	txPaymentModule := getInstance[transaction_payment.Module[BlockNumberType]]()

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
	decoder := types.NewModuleDecoder[BlockNumberType](modules, extra)
	runtimeExtrinsic := extrinsic.New[BlockNumberType](modules, extra)
	auraModule := getInstance[aura.Module[BlockNumberType]]()
	grandpaModule := getInstance[grandpa.Module[BlockNumberType]]()
	txPaymentsModule := getInstance[transaction_payment.Module[BlockNumberType]]()

	executiveModule := executive.New[BlockNumberType](
		getInstance[system.Module[BlockNumberType]](),
		runtimeExtrinsic,
		hooks.DefaultOnRuntimeUpgrade{},
	)

	sessions := []primitives.Session{
		auraModule,
		grandpaModule,
	}

	apis := []primitives.ApiModule{
		core.New[BlockNumberType](executiveModule, decoder, RuntimeVersion),
		blockbuilder.New[BlockNumberType](runtimeExtrinsic, executiveModule, decoder),
		taggedtransactionqueue.New[BlockNumberType](executiveModule, decoder),
		metadata.New[BlockNumberType](runtimeExtrinsic),
		apiAura.New[BlockNumberType](auraModule),
		apiGrandpa.New[BlockNumberType](grandpaModule),
		account_nonce.New[BlockNumberType](getInstance[system.Module[BlockNumberType]]()),
		apiTxPayments.New[BlockNumberType](decoder, txPaymentsModule),
		apiTxPaymentsCall.NewCallApi[BlockNumberType](decoder, txPaymentsModule),
		session_keys.New(sessions),
		offchain_worker.New[BlockNumberType](executiveModule),
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
		Module(apiAura.ApiModuleName).(apiAura.Module[BlockNumberType]).
		SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return runtimeApi().
		Module(apiAura.ApiModuleName).(apiAura.Module[BlockNumberType]).
		Authorities()
}

//go:export AccountNonceApi_account_nonce
func AccountNonceApiAccountNonce(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(account_nonce.ApiModuleName).(account_nonce.Module[BlockNumberType]).
		AccountNonce(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_info
func TransactionPaymentApiQueryInfo(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPayments.ApiModuleName).(apiTxPayments.Module[BlockNumberType]).
		QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPayments.ApiModuleName).(apiTxPayments.Module[BlockNumberType]).
		QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return runtimeApi().
		Module(apiTxPaymentsCall.ApiModuleName).(apiTxPaymentsCall.Module[BlockNumberType]).
		QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(apiTxPaymentsCall.ApiModuleName).(apiTxPaymentsCall.Module[BlockNumberType]).
		QueryCallFeeDetails(dataPtr, dataLen)
}

//go:export Metadata_metadata
func Metadata(_, _ int32) int64 {
	return runtimeApi().
		Module(metadata.ApiModuleName).(metadata.Module[BlockNumberType]).
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
		Module(apiGrandpa.ApiModuleName).(apiGrandpa.Module[BlockNumberType]).
		Authorities()
}

//go:export OffchainWorkerApi_offchain_worker
func OffchainWorkerApiOffchainWorker(dataPtr int32, dataLen int32) int64 {
	runtimeApi().
		Module(offchain_worker.ApiModuleName).(offchain_worker.Module[BlockNumberType]).
		OffchainWorker(dataPtr, dataLen)

	return 0
}
