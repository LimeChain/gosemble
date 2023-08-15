/*
Targets WebAssembly MVP
*/
package main

import (
	"math/big"
	"reflect"

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
		sm.NewConfig(constants.BlockHashCount, BlockWeights, BlockLength, *RuntimeVersion))

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

func getInstance[T primitives.Module]() T {
	for _, module := range modules {
		if reflect.TypeOf(module) == reflect.TypeOf(*new(T)) {
			return modules[module.GetIndex()].(T)
		}
	}
	log.Critical("unknown type T for module instance")
	panic("unreachable")
}

func newExecutiveModule() executive.Module {
	return executive.New(
		getInstance[sm.SystemModule](),
		extrinsic.New(modules),
		getInstance[aura.Module](),
	)
}

func newModuleDecoder() types.ModuleDecoder {
	return types.NewModuleDecoder(modules, newSignedExtra())
}

func newSignedExtra() primitives.SignedExtra {
	systeModule := getInstance[sm.SystemModule]()
	balancesModule := getInstance[bm.BalancesModule]()
	txPaymentModule := getInstance[tpm.TransactionPaymentModule]()

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

func runtimeApi() types.RuntimeApi {
	executiveModule := newExecutiveModule()
	decoder := newModuleDecoder()
	runtimeExtrinsic := extrinsic.New(modules)
	auraModule := getInstance[aura.Module]()
	grandpaModule := getInstance[grandpa.Module]()

	sessions := []primitives.Session{
		auraModule,
		grandpaModule,
	}

	apis := []primitives.ApiModule{
		core.New(executiveModule, decoder, RuntimeVersion),
		blockbuilder.New(runtimeExtrinsic, executiveModule, decoder),
		taggedtransactionqueue.New(executiveModule, decoder),
		metadata.New(modules),
		auraModule,
		grandpaModule,
		account_nonce.New(getInstance[sm.SystemModule]()),
		transaction_payment.New(decoder, getInstance[tpm.TransactionPaymentModule]()),
		transaction_payment.NewCallApi(decoder, getInstance[tpm.TransactionPaymentModule]()),
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
		Module(aura.ApiModuleName).(aura.Module).
		SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return runtimeApi().
		Module(aura.ApiModuleName).(aura.Module).
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
		Module(transaction_payment.ApiModuleName).(transaction_payment.Module).
		QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(transaction_payment.ApiModuleName).(transaction_payment.Module).
		QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return runtimeApi().
		Module(transaction_payment.CallApiModuleName).(transaction_payment.TransactionPaymentCallApi).
		QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(transaction_payment.CallApiModuleName).(transaction_payment.TransactionPaymentCallApi).
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
		Module(grandpa.ApiModuleName).(grandpa.Module).
		Authorities()
}

//go:export OffchainWorkerApi_offchain_worker
func OffchainWorkerApiOffchainWorker(dataPtr int32, dataLen int32) int64 {
	runtimeApi().
		Module(offchain_worker.ApiModuleName).(offchain_worker.Module).
		OffchainWorker(dataPtr, dataLen)

	return 0
}
