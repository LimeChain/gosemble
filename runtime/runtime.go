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
	am "github.com/LimeChain/gosemble/frame/aura/module"
	bm "github.com/LimeChain/gosemble/frame/balances/module"
	blockbuilder "github.com/LimeChain/gosemble/frame/block_builder"
	"github.com/LimeChain/gosemble/frame/core"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/frame/metadata"
	"github.com/LimeChain/gosemble/frame/offchain_worker"
	"github.com/LimeChain/gosemble/frame/session_keys"
	sm "github.com/LimeChain/gosemble/frame/system/module"
	taggedtransactionqueue "github.com/LimeChain/gosemble/frame/tagged_transaction_queue"
	tm "github.com/LimeChain/gosemble/frame/testable/module"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	tpm "github.com/LimeChain/gosemble/frame/transaction_payment/module"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

// TODO:
// remove the _start export and find a way to call it from the runtime to initialize the memory.
// TinyGo requires to have a main function to compile to Wasm.
func main() {}

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
		sm.NewConfig(constants.BlockHashCount, constants.RuntimeVersion))

	auraModule := am.NewModule(AuraIndex,
		am.NewConfig(timestamp.MinimumPeriod, AuraMaxAuthorites, false))

	timestampModule := tsm.NewModule(TimestampIndex,
		tsm.NewConfig(auraModule, timestamp.MinimumPeriod))

	grandpaModule := grandpa.NewModule(GrandpaIndex)

	balancesModule := bm.NewBalancesModule(BalancesIndex,
		bm.NewConfig(BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule))

	tpmModule := tpm.NewTransactionPaymentModule()
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
	return executive.New(modules[SystemIndex].(sm.SystemModule), extrinsic.New(modules))
}

func newModuleDecoder() types.ModuleDecoder {
	return types.NewModuleDecoder(modules)
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
	return aura.SlotDuration()
}

//go:export AuraApi_authorities
func AuraApiAuthorities(_, _ int32) int64 {
	return aura.Authorities()
}

//go:export AccountNonceApi_account_nonce
func AccountNonceApiAccountNonce(dataPtr int32, dataLen int32) int64 {
	return account_nonce.AccountNonce(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_info
func TransactionPaymentApiQueryInfo(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder()).
		QueryInfo(dataPtr, dataLen)
}

//go:export TransactionPaymentApi_query_fee_details
func TransactionPaymentApiQueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder()).
		QueryFeeDetails(dataPtr, dataLen)
}

//go:export TransactionPaymentCallApi_query_call_info
func TransactionPaymentCallApiQueryCallInfo(dataPtr int32, dataLan int32) int64 {
	return transaction_payment.
		New(newModuleDecoder()).
		QueryCallInfo(dataPtr, dataLan)
}

//go:export TransactionPaymentCallApi_query_call_fee_details
func TransactionPaymentCallApiQueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	return transaction_payment.
		New(newModuleDecoder()).
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
	return session_keys.GenerateSessionKeys(dataPtr, dataLen)
}

//go:export SessionKeys_decode_session_keys
func SessionKeysDecodeSessionKeys(dataPtr int32, dataLen int32) int64 {
	return session_keys.DecodeSessionKeys(dataPtr, dataLen)
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
