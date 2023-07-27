package config

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/timestamp"
	am "github.com/LimeChain/gosemble/frame/aura/module"
	bm "github.com/LimeChain/gosemble/frame/balances/module"
	gm "github.com/LimeChain/gosemble/frame/grandpa/module"
	sm "github.com/LimeChain/gosemble/frame/system/module"
	tm "github.com/LimeChain/gosemble/frame/testable/module"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	tpm "github.com/LimeChain/gosemble/frame/transaction_payment/module"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Modules contains all the modules used by the runtime.
var Modules = initializeModules()

const (
	SystemIndex sc.U8 = iota
	TimestampIndex
	AuraIndex
	GrandpaIndex
	BalancesIndex
	TxPaymentsIndex
	TestableIndex = 255
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

func initializeModules() map[sc.U8]types.Module {
	systemModule := sm.NewSystemModule(SystemIndex,
		sm.NewConfig(constants.BlockHashCount, constants.RuntimeVersion))

	auraModule := am.NewModule(AuraIndex,
		am.NewConfig(timestamp.MinimumPeriod, AuraMaxAuthorites, false))

	timestampModule := tsm.NewModule(TimestampIndex,
		tsm.NewConfig(auraModule, timestamp.MinimumPeriod))

	grandpaModule := gm.NewGrandpaModule()

	balancesModule := bm.NewBalancesModule(BalancesIndex,
		bm.NewConfig(BalancesMaxLocks, BalancesMaxReserves, BalancesExistentialDeposit, systemModule))

	tpmModule := tpm.NewTransactionPaymentModule()
	testableModule := tm.NewTestingModule()

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
