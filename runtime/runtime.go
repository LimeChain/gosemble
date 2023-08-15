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

type api struct {
	core           core.Module
	blockBuilder   blockbuilder.Module
	taggedTxQueue  taggedtransactionqueue.Module
	metadata       metadata.Module
	aura           aura.Module
	grandpa        grandpa.Module
	accountNonce   account_nonce.Module
	txPayment      transaction_payment.Module
	txPaymentCall  transaction_payment.TransactionPaymentCallApi
	sessionKeys    session_keys.Module
	offchainWorker offchain_worker.Module
}

func (api api) Items() sc.Sequence[primitives.ApiItem] {
	return sc.Sequence[primitives.ApiItem]{
		api.core.Item(),
		api.blockBuilder.Item(),
		api.taggedTxQueue.Item(),
		api.metadata.Item(),
		api.aura.Item(),
		api.grandpa.Item(),
		api.accountNonce.Item(),
		api.txPayment.Item(),
		api.txPaymentCall.Item(),
		api.sessionKeys.Item(),
		api.offchainWorker.Item(),
	}
}

func apiModules() api {
	executiveModule := newExecutiveModule()
	decoder := newModuleDecoder()
	runtimeExtrinsic := extrinsic.New(modules)

	sessions := []primitives.Session{
		modules[AuraIndex].(aura.Module),
		modules[GrandpaIndex].(grandpa.Module),
	}

	apis := api{
		core:           core.New(executiveModule, decoder, RuntimeVersion),
		blockBuilder:   blockbuilder.New(runtimeExtrinsic, executiveModule, decoder),
		taggedTxQueue:  taggedtransactionqueue.New(executiveModule, decoder),
		metadata:       metadata.New(modules),
		aura:           modules[AuraIndex].(aura.Module),
		grandpa:        modules[GrandpaIndex].(grandpa.Module),
		accountNonce:   account_nonce.New(modules[SystemIndex].(sm.SystemModule)),
		txPayment:      transaction_payment.New(decoder, modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)),
		txPaymentCall:  transaction_payment.NewCallApi(decoder, modules[TxPaymentsIndex].(tpm.TransactionPaymentModule)),
		sessionKeys:    session_keys.New(sessions),
		offchainWorker: offchain_worker.New(executiveModule),
	}

	RuntimeVersion.SetApis(apis.Items())

	return apis
}

//go:export Core_version
func CoreVersion(_ int32, _ int32) int64 {
	return apiModules().core.Version()
}

//go:export Core_initialize_block
func CoreInitializeBlock(dataPtr int32, dataLen int32) int64 {
	apiModules().core.InitializeBlock(dataPtr, dataLen)

	return 0
}

//go:export Core_execute_block
func CoreExecuteBlock(dataPtr int32, dataLen int32) int64 {
	apiModules().core.ExecuteBlock(dataPtr, dataLen)

	return 0
}

//go:export BlockBuilder_apply_extrinsic
func BlockBuilderApplyExtrinsic(dataPtr int32, dataLen int32) int64 {
	return apiModules().blockBuilder.ApplyExtrinsic(dataPtr, dataLen)
}

//go:export BlockBuilder_finalize_block
func BlockBuilderFinalizeBlock(_, _ int32) int64 {
	return apiModules().blockBuilder.FinalizeBlock()
}

//go:export BlockBuilder_inherent_extrinsics
func BlockBuilderInherentExtrinsics(dataPtr int32, dataLen int32) int64 {
	return apiModules().blockBuilder.InherentExtrinsics(dataPtr, dataLen)
}

//go:export BlockBuilder_check_inherents
func BlockBuilderCheckInherents(dataPtr int32, dataLen int32) int64 {
	return apiModules().blockBuilder.CheckInherents(dataPtr, dataLen)
}

//go:export TaggedTransactionQueue_validate_transaction
func TaggedTransactionQueueValidateTransaction(dataPtr int32, dataLen int32) int64 {
	return apiModules().taggedTxQueue.ValidateTransaction(dataPtr, dataLen)
}

//go:export AuraApi_slot_duration
func AuraApiSlotDuration(_, _ int32) int64 {
	return apiModules().aura.SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return apiModules().aura.Authorities()
}

//go:export AccountNonceApi_account_nonce
func AccountNonceApiAccountNonce(dataPtr int32, dataLen int32) int64 {
	return apiModules().accountNonce.AccountNonce(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_info
func TransactionPaymentApiQueryInfo(dataPtr int32, dataLen int32) int64 {
	return apiModules().txPayment.QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return apiModules().txPayment.QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return apiModules().txPaymentCall.QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return apiModules().txPaymentCall.QueryCallFeeDetails(dataPtr, dataLen)
}

//go:export Metadata_metadata
func Metadata(_, _ int32) int64 {
	return apiModules().metadata.Metadata()
}

//go:export SessionKeys_generate_session_keys
func SessionKeysGenerateSessionKeys(dataPtr int32, dataLen int32) int64 {
	return apiModules().sessionKeys.GenerateSessionKeys(dataPtr, dataLen)
}

//go:export SessionKeys_decode_session_keys
func SessionKeysDecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	return apiModules().sessionKeys.DecodeSessionKeys(dataPtr, dataLen)
}

//go:export GrandpaApi_grandpa_authorities
func GrandpaApiAuthorities(_, _ int32) int64 {
	return apiModules().grandpa.Authorities()
}

//go:export OffchainWorkerApi_offchain_worker
func OffchainWorkerApiOffchainWorker(dataPtr int32, dataLen int32) int64 {
	apiModules().offchainWorker.OffchainWorker(dataPtr, dataLen)

	return 0
}
