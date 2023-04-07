package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type BalancesModule struct {
	functions map[sc.U8]primitives.Call
}

func NewBalancesModule() BalancesModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[balances.FunctionTransferIndex] = dispatchables.NewTransferCall(nil)
	functions[balances.FunctionSetBalanceIndex] = dispatchables.NewSetBalanceCall(nil)
	functions[balances.FunctionForceTransferIndex] = dispatchables.NewForceTransferCall(nil)
	functions[balances.FunctionTransferKeepAliveIndex] = dispatchables.NewTransferKeepAliveCall(nil)
	functions[balances.FunctionTransferAllIndex] = dispatchables.NewTransferAllCall(nil)
	functions[balances.FunctionForceFreeIndex] = dispatchables.NewForceFreeCall(nil)

	return BalancesModule{
		functions: functions,
	}
}

func (bm BalancesModule) Functions() map[sc.U8]primitives.Call {
	return bm.functions
}

func (bm BalancesModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (bm BalancesModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (bm BalancesModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return bm.metadataTypes(), primitives.MetadataModule{
		Name:      "Balances",
		Storage:   sc.Option[primitives.MetadataModuleStorage]{}, // TODO:
		Call:      sc.NewOption[sc.Compact](sc.ToCompact(metadata.BalancesCalls)),
		Event:     sc.NewOption[sc.Compact](nil),                    // TODO:
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{}, // TODO:
		Error:     sc.NewOption[sc.Compact](nil),                    // TODO:
		Index:     balances.ModuleIndex,
	}
}

func (bm BalancesModule) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.BalancesCalls, "Balances calls", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU128),
					},
					balances.FunctionTransferIndex,
					"Transfer some liquid free balance to another account."),
				primitives.NewMetadataDefinitionVariant(
					"set_balance",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU128),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU128),
					},
					balances.FunctionSetBalanceIndex,
					"Set the balances of a given account."),
				primitives.NewMetadataDefinitionVariant(
					"force_transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU128),
					},
					balances.FunctionForceTransferIndex,
					"Exactly as `transfer`, except the origin must be root and the source account may be specified."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_keep_alive",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU128),
					},
					balances.FunctionTransferKeepAliveIndex,
					"Same as the [`transfer`] call, but with a check that the transfer will not kill the origin account."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_all",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesBool),
					},
					balances.FunctionTransferAllIndex,
					"Transfer the entire transferable balance from the caller account."),
				primitives.NewMetadataDefinitionVariant(
					"force_unreserve",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesMultiAddress),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					balances.FunctionForceFreeIndex,
					"Unreserve some balance from a user by force."),
			})),
	}
}
