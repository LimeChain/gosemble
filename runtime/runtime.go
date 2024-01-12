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
	genesisbuilder "github.com/LimeChain/gosemble/api/genesis_builder"
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
	AuraMaxAuthorities = 100
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
	BalancesExistentialDeposit = sc.NewU128(1 * constants.Dollar)
)

var (
	DbWeight = constants.RocksDbWeight
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

var (
	logger = log.NewLogger()
	// Modules contains all the modules used by the runtime.
	mdGenerator = primitives.NewMetadataTypeGenerator()
	modules     = initializeModules()
)

func initializeModules() []primitives.Module {
	blockWeights, err := system.WithSensibleDefaults(constants.MaximumBlockWeight, constants.NormalDispatchRatio)
	if err != nil {
		logger.Critical(err.Error())
	}

	blockLength, err := system.MaxWithNormalRatio(constants.FiveMbPerBlockPerExtrinsic, constants.NormalDispatchRatio)
	if err != nil {
		logger.Critical(err.Error())
	}

	systemModule := system.New(
		SystemIndex,
		system.NewConfig(constants.BlockHashCount, blockWeights, blockLength, DbWeight, *RuntimeVersion),
		mdGenerator,
		logger,
	)

	auraModule := aura.New(
		AuraIndex,
		aura.NewConfig(
			primitives.PublicKeySr25519,
			DbWeight,
			TimestampMinimumPeriod,
			AuraMaxAuthorities,
			false,
			systemModule.StorageDigest,
		),
		mdGenerator,
	)

	timestampModule := timestamp.New(
		TimestampIndex,
		timestamp.NewConfig(auraModule, DbWeight, TimestampMinimumPeriod),
		mdGenerator,
	)

	grandpaModule := grandpa.New(GrandpaIndex, logger, mdGenerator)

	balancesModule := balances.New(
		BalancesIndex,
		balances.NewConfig(DbWeight, BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule),
		logger,
		mdGenerator,
	)

	tpmModule := transaction_payment.New(
		TxPaymentsIndex,
		transaction_payment.NewConfig(OperationalFeeMultiplier, WeightToFee, LengthToFee, blockWeights),
		mdGenerator,
	)

	testableModule := tm.New(TestableIndex, mdGenerator)

	//addr := fmt.Sprintf("%p", &mdGenerator)
	//log.NewLogger().Info("Addr in initializeModules: " + addr)

	return []primitives.Module{
		systemModule,
		timestampModule,
		auraModule,
		grandpaModule,
		balancesModule,
		tpmModule,
		testableModule,
	}
}

func newSignedExtra() primitives.SignedExtra {
	systemModule := primitives.MustGetModule(SystemIndex, modules).(system.Module)
	balancesModule := primitives.MustGetModule(BalancesIndex, modules).(balances.Module)
	txPaymentModule := primitives.MustGetModule(TxPaymentsIndex, modules).(transaction_payment.Module)

	extras := []primitives.SignedExtension{
		sysExtensions.NewCheckNonZeroAddress(),
		sysExtensions.NewCheckSpecVersion(systemModule),
		sysExtensions.NewCheckTxVersion(systemModule),
		sysExtensions.NewCheckGenesis(systemModule),
		sysExtensions.NewCheckMortality(systemModule),
		sysExtensions.NewCheckNonce(systemModule),
		sysExtensions.NewCheckWeight(systemModule),
		txExtensions.NewChargeTransactionPayment(systemModule, txPaymentModule, balancesModule),
	}

	//addr := fmt.Sprintf("%p", &mdGenerator)
	//log.NewLogger().Info("Addr in signedExtra: " + addr)

	return primitives.NewSignedExtra(extras, mdGenerator)
}

func runtimeApi() types.RuntimeApi {
	//addr := fmt.Sprintf("%p", &mdGenerator)
	//log.NewLogger().Info("Addr in runtimeApi: " + addr)
	extra := newSignedExtra()
	decoder := types.NewRuntimeDecoder(modules, extra, logger)
	runtimeExtrinsic := extrinsic.New(modules, extra, mdGenerator, logger)
	systemModule := primitives.MustGetModule(SystemIndex, modules).(system.Module)
	auraModule := primitives.MustGetModule(AuraIndex, modules).(aura.Module)
	grandpaModule := primitives.MustGetModule(GrandpaIndex, modules).(grandpa.Module)
	txPaymentsModule := primitives.MustGetModule(TxPaymentsIndex, modules).(transaction_payment.Module)

	executiveModule := executive.New(
		systemModule,
		runtimeExtrinsic,
		hooks.DefaultOnRuntimeUpgrade{},
		logger,
	)

	sessions := []primitives.Session{
		auraModule,
		grandpaModule,
	}

	coreApi := core.New(executiveModule, decoder, RuntimeVersion, logger)
	blockBuilderApi := blockbuilder.New(runtimeExtrinsic, executiveModule, decoder, logger)
	taggedTxQueueApi := taggedtransactionqueue.New(executiveModule, decoder, logger)
	auraApi := apiAura.New(auraModule, logger)
	grandpaApi := apiGrandpa.New(grandpaModule, logger)
	accountNonceApi := account_nonce.New(systemModule, logger)
	txPaymentsApi := apiTxPayments.New(decoder, txPaymentsModule, logger)
	txPaymentsCallApi := apiTxPaymentsCall.New(decoder, txPaymentsModule, logger)
	sessionKeysApi := session_keys.New(sessions, logger)
	offchainWorkerApi := offchain_worker.New(executiveModule, logger)
	genesisBuilderApi := genesisbuilder.New(modules, logger)

	metadataApi := metadata.New(
		runtimeExtrinsic,
		[]primitives.RuntimeApiModule{
			coreApi,
			blockBuilderApi,
			taggedTxQueueApi,
			auraApi,
			grandpaApi,
			accountNonceApi,
			txPaymentsApi,
			txPaymentsCallApi,
			sessionKeysApi,
			offchainWorkerApi,
		},
		logger,
		mdGenerator,
	)

	apis := []primitives.ApiModule{
		coreApi,
		blockBuilderApi,
		taggedTxQueueApi,
		metadataApi,
		auraApi,
		grandpaApi,
		accountNonceApi,
		txPaymentsApi,
		txPaymentsCallApi,
		sessionKeysApi,
		offchainWorkerApi,
		genesisBuilderApi,
	}

	runtimeApi := types.NewRuntimeApi(apis, logger)

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

//go:export Metadata_metadata_at_version
func MetadataAtVersion(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(metadata.ApiModuleName).(metadata.Module).
		MetadataAtVersion(dataPtr, dataLen)
}

//go:export Metadata_metadata_versions
func MetadataVersions(_, _ int32) int64 {
	return runtimeApi().
		Module(metadata.ApiModuleName).(metadata.Module).
		MetadataVersions()
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

//go:export GenesisBuilder_create_default_config
func GenesisBuilderCreateDefaultConfig(_, _ int32) int64 {
	return runtimeApi().
		Module(genesisbuilder.ApiModuleName).(genesisbuilder.Module).
		CreateDefaultConfig()
}

//go:export GenesisBuilder_build_config
func GenesisBuilderBuildConfig(dataPtr int32, dataLen int32) int64 {
	return runtimeApi().
		Module(genesisbuilder.ApiModuleName).(genesisbuilder.Module).
		BuildConfig(dataPtr, dataLen)
}
