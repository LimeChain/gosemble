package metadata

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/constants/transaction_payment"
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

func Metadata() int64 {
	metadata := buildMetadata()
	bMetadata := sc.BytesToSequenceU8(metadata.Bytes())

	return utils.BytesToOffsetAndSize(bMetadata.Bytes())
}

func buildMetadata() primitives.Metadata {
	metadataTypes := append(primitiveTypes(), runtimeTypes()...)

	var modules sc.Sequence[primitives.MetadataModule]

	for _, module := range config.Modules {
		mTypes, mModule := module.Metadata()

		metadataTypes = append(metadataTypes, mTypes...)
		modules = append(modules, mModule)
	}

	extrinsic := primitives.MetadataExtrinsic{
		Type:    sc.ToCompact(metadata.UncheckedExtrinsic),
		Version: types.ExtrinsicFormatVersion,
		SignedExtensions: sc.Sequence[primitives.MetadataSignedExtension]{
			primitives.NewMetadataSignedExtension("CheckNonZeroSender", metadata.CheckNonZeroSender, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("CheckSpecVersion", metadata.CheckSpecVersion, metadata.PrimitiveTypesU32),
			primitives.NewMetadataSignedExtension("CheckTxVersion", metadata.CheckTxVersion, metadata.PrimitiveTypesU32),
			primitives.NewMetadataSignedExtension("CheckGenesis", metadata.CheckGenesis, metadata.TypesH256),
			primitives.NewMetadataSignedExtension("CheckMortality", metadata.CheckMortality, metadata.TypesH256),
			primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("CheckWeight", metadata.CheckWeight, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("ChargeTransactionPayment", metadata.ChargeTransactionPayment, metadata.TypesEmptyTuple),
		},
	}

	runtimeV14Metadata := primitives.RuntimeMetadataV14{
		Types:     metadataTypes,
		Modules:   modules,
		Extrinsic: extrinsic,
		Type:      sc.ToCompact(metadata.Runtime),
	}

	return primitives.NewMetadata(runtimeV14Metadata)
}

// primitiveTypes returns all primitive types
func primitiveTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.PrimitiveTypesBool, "bool", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveBoolean)),
		primitives.NewMetadataType(metadata.PrimitiveTypesChar, "char", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveChar)),
		primitives.NewMetadataType(metadata.PrimitiveTypesString, "string", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveString)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU8, "U8", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU16, "U16", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU32, "U32", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU64, "U64", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU128, "U128", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU256, "U256", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU256)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI8, "I8", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI16, "I16", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI32, "I32", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI64, "I64", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI128, "I128", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI256, "I256", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI256)),
	}
}

func runtimeTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.TypesFixedSequence4U8, "[4]byte", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionFixedSequence(4, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence20U8, "[20]byte", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionFixedSequence(20, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence32U8, "[32]byte", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionFixedSequence(32, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence64U8, "[64]byte", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionFixedSequence(64, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence65U8, "[65]byte", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionFixedSequence(65, sc.ToCompact(metadata.PrimitiveTypesU8))),

		primitives.NewMetadataType(metadata.TypesCompactU32, "CompactU32", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.TypesCompactU64, "CompactU64", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))),
		primitives.NewMetadataType(metadata.TypesCompactU128, "CompactU128", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))),

		primitives.NewMetadataType(metadata.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),
		primitives.NewMetadataType(metadata.TypesVecBlockNumEventIndex, "Vec<BlockNumber, EventIndex>", sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleU32U32))),
		primitives.NewMetadataType(metadata.TypesTupleU32U32, "(U32, U32)", sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.PrimitiveTypesU32), sc.ToCompact(metadata.PrimitiveTypesU32)})),
		primitives.NewMetadataType(metadata.TypesSystemEventStorage, "Vec<Box<EventRecord<T::RuntimeEvent, T::Hash>>>", sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesEventRecord))),
		primitives.NewMetadataTypeWithParams(metadata.TypesEventRecord, "frame_system EventRecord", sc.Sequence[sc.Str]{"frame_system", "EventRecord"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPhase, "phase", "Phase"),       // phase
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesRuntimeEvent, "event", "E"),    // event
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesVecTopics, "topics", "Vec<T>"), // topics
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesRuntimeEvent, "E"),
				primitives.NewMetadataTypeParameter(metadata.TypesH256, "T"),
			}),
		primitives.NewMetadataType(metadata.TypesPhase, "frame_system Phase", sc.Sequence[sc.Str]{"frame_system", "Phase"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"ApplyExtrinsic",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),
					},
					0,
					"Phase.ApplyExtrinsic"),
				primitives.NewMetadataDefinitionVariant(
					"Finalization",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"Phase.Finalization"),
				primitives.NewMetadataDefinitionVariant(
					"Initialization",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"Phase.Init"),
			})),
		primitives.NewMetadataType(metadata.TypesRuntimeEvent, "node_template_runtime RuntimeEvent", sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"System",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSystemEvent, "frame_system::Event<Runtime>"),
					},
					system.ModuleIndex,
					"Events.System"),
				primitives.NewMetadataDefinitionVariant( // TODO:
					"Grandpa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					grandpa.ModuleIndex,
					"Events.Grandpa"),
				primitives.NewMetadataDefinitionVariant(
					"Balances",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesBalancesEvent, "pallet_balances::Event<Runtime>"),
					},
					balances.ModuleIndex,
					"Events.Balances"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionPayment",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTransactionPaymentEvents, "pallet_transaction_payment::Event<Runtime>"),
					},
					transaction_payment.ModuleIndex,
					"Events.TransactionPayment"),
			})),
		primitives.NewMetadataType(metadata.TypesSystemEvent, "frame_system pallet Event", sc.Sequence[sc.Str]{"frame_system", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"ExtrinsicSuccess",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
					},
					0,
					"Event.ExtrinsicSuccess"),
				primitives.NewMetadataDefinitionVariant(
					"ExtrinsicFailed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchError, "dispatch_error", "DispatchError"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchInfo, "dispatch_info", "DispatchInfo"),
					},
					1,
					"Events.ExtrinsicFailed"),
				primitives.NewMetadataDefinitionVariant(
					"CodeUpdated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"Events.CodeUpdated"),
				primitives.NewMetadataDefinitionVariant(
					"NewAccount",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
					},
					3,
					"Events.NewAccount"),
				primitives.NewMetadataDefinitionVariant(
					"KilledAccount",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
					},
					4,
					"Events.KilledAccount"),
				primitives.NewMetadataDefinitionVariant(
					"Remarked",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "sender", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesH256, "hash", "T::Hash"),
					},
					5,
					"Events.Remarked"),
			})),

		primitives.NewMetadataType(metadata.TypesDispatchInfo, "DispatchInfo", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "weight", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchClass, "class", "DispatchClass"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPays, "pays_fee", "Pays"),
			},
		)),
		primitives.NewMetadataType(metadata.TypesDispatchClass, "DispatchClass", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchClass"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Normal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"DispatchClass.Normal"),
				primitives.NewMetadataDefinitionVariant(
					"Operational",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"DispatchClass.Operational"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"DispatchClass.Mandatory"),
			})),
		primitives.NewMetadataType(metadata.TypesPays, "Pays", sc.Sequence[sc.Str]{"frame_support", "dispatch", "Pays"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Yes",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Pays.Yes"),
				primitives.NewMetadataDefinitionVariant(
					"No",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"Pays.No"),
			})),

		primitives.NewMetadataType(metadata.TypesDispatchError, "DispatchError", sc.Sequence[sc.Str]{"sp_runtime", "DispatchError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"DispatchError.Other"),
				primitives.NewMetadataDefinitionVariant(
					"Cannotlookup",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"DispatchError.Cannotlookup"),
				primitives.NewMetadataDefinitionVariant(
					"BadOrigin",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"DispatchError.BadOrigin"),
				primitives.NewMetadataDefinitionVariant(
					"Module",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesModuleError, "ModuleError"),
					},
					3,
					"DispatchError.Module"),
				primitives.NewMetadataDefinitionVariant(
					"ConsumerRemaining",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					4,
					"DispatchError.ConsumerRemaining"),
				primitives.NewMetadataDefinitionVariant(
					"NoProviders",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					5,
					"DispatchError.NoProviders"),
				primitives.NewMetadataDefinitionVariant(
					"TooManyConsumers",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					6,
					"DispatchError.TooManyConsumers"),
				primitives.NewMetadataDefinitionVariant(
					"Token",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTokenError, "TokenError"),
					},
					7,
					"DispatchError.Token"),
				primitives.NewMetadataDefinitionVariant(
					"Arithmetic",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesArithmeticError, "ArithmeticError"),
					},
					8,
					"DispatchError.Arithmetic"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionalError",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTransactionalError, "TransactionalError"),
					},
					9,
					"DispatchError.TransactionalError"),
				primitives.NewMetadataDefinitionVariant(
					"Exhausted",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					10,
					"DispatchError.Exhausted"),
				primitives.NewMetadataDefinitionVariant(
					"Corruption",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					11,
					"DispatchError.Corruption"),
				primitives.NewMetadataDefinitionVariant(
					"Unavailable",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					12,
					"DispatchError.Unavailable"),
			})),
		primitives.NewMetadataType(metadata.TypesModuleError, "ModuleError", sc.Sequence[sc.Str]{"sp_runtime", "ModuleError"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU8, "index", "u8"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesFixedSequence4U8, "error", "[u8; MAX_MODULE_ERROR_ENCODED_SIZE]"),
			})),
		primitives.NewMetadataType(metadata.TypesTokenError, "TokenError", sc.Sequence[sc.Str]{"sp_runtime", "TokenError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"NoFunds",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"TokenError.NoFunds"),
				primitives.NewMetadataDefinitionVariant(
					"WouldDie",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"TokenError.WouldDie"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"TokenError.BelowMinimum"),
				primitives.NewMetadataDefinitionVariant(
					"CannotCreate",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					3,
					"TokenError.CannotCreate"),
				primitives.NewMetadataDefinitionVariant(
					"UnknownAsset",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					4,
					"TokenError.UnknownAsset"),
				primitives.NewMetadataDefinitionVariant(
					"Frozen",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					5,
					"TokenError.Frozen"),
				primitives.NewMetadataDefinitionVariant(
					"Unsupported",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					6,
					"TokenError.Unsupported"),
			})),
		primitives.NewMetadataType(metadata.TypesArithmeticError, "ArithmeticError", sc.Sequence[sc.Str]{"sp_arithmetic", "ArithmeticError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Underflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"ArithmeticError.Underflow"),
				primitives.NewMetadataDefinitionVariant(
					"Overflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"ArithmeticError.Overflow"),
				primitives.NewMetadataDefinitionVariant(
					"DivisionByZero",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"ArithmeticError.DivisionByZero"),
			})),
		primitives.NewMetadataType(metadata.TypesTransactionalError, "TransactionalError", sc.Sequence[sc.Str]{"sp_runtime", "TransactionalError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"LimitReached",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"TransactionalError.LimitReached"),
				primitives.NewMetadataDefinitionVariant(
					"NoLayer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"TransactionalError.NoLayer"),
			})),

		primitives.NewMetadataType(metadata.TypesBalancesEvent, "pallet_balances pallet Event", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Endowed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free_balance", "T::Balance"),
					},
					0,
					"Event.Endowed"),
				primitives.NewMetadataDefinitionVariant(
					"DustLost",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					1,
					"Events.DustLost"),
				primitives.NewMetadataDefinitionVariant(
					"Transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					2,
					"Events.Transfer"),
				primitives.NewMetadataDefinitionVariant(
					"BalanceSet",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "T::Balance"),
					},
					3,
					"Events.BalanceSet"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					4,
					"Events.Reserved"),
				primitives.NewMetadataDefinitionVariant(
					"Unreserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					5,
					"Events.Unreserved"),
				primitives.NewMetadataDefinitionVariant(
					"ReserveRepatriated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesBalanceStatus, "destination_status", "Status"),
					},
					6,
					"Events.ReserveRepatriated"),
				primitives.NewMetadataDefinitionVariant(
					"Deposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					7,
					"Event.Deposit"),
				primitives.NewMetadataDefinitionVariant(
					"Withdraw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					8,
					"Event.Withdraw"),
				primitives.NewMetadataDefinitionVariant(
					"Slashed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					9,
					"Event.Slashed"),
			},
		)),

		primitives.NewMetadataType(metadata.TypesBalanceStatus, "BalanceStatus", sc.Sequence[sc.Str]{"frame_support", "traits", "tokens", "misc", "BalanceStatus"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Free",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"BalanceStatus.Free"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"BalanceStatus.Reserved"),
			})),

		primitives.NewMetadataType(metadata.TypesVecTopics, "Vec<Topics>", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesH256))),

		primitives.NewMetadataType(metadata.TypesLastRuntimeUpgradeInfo, "LastRuntimeUpgradeInfo", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU32),
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString),
			})),

		primitives.NewMetadataType(metadata.TypesSystemErrors, "frame_system pallet Error", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"InvalidSpecName",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"The name of specification does not match between the current runtime and the new runtime."),
				primitives.NewMetadataDefinitionVariant(
					"SpecVersionNeedsToIncrease",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"The specification version is not allowed to decrease between the current runtime and the new runtime."),
				primitives.NewMetadataDefinitionVariant(
					"FailedToExtractRuntimeVersion",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"Failed to extract the runtime version from the new runtime.  Either calling `Core_version` or decoding `RuntimeVersion` failed."),
				primitives.NewMetadataDefinitionVariant(
					"NonDefaultComposite",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					3,
					"Suicide called when the account has non-default composite data."),
				primitives.NewMetadataDefinitionVariant(
					"NonZeroRefCount",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					4,
					"There is a non-zero reference count preventing the account from being purged."),
				primitives.NewMetadataDefinitionVariant(
					"CallFiltered",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					5,
					"The origin filter prevent the call to be dispatched."),
			})),

		primitives.NewMetadataType(metadata.TypesBlockLength,
			"frame_system limits BlockLength",
			sc.Sequence[sc.Str]{"frame_system", "limits", "BlockLength"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassU32, "max", "PerDispatchClass<u32>"), // max
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassU32, "PerDispatchClass[U32]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "mandatory", "T"),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.PrimitiveTypesU32, "T"),
		),

		primitives.NewMetadataType(metadata.TypesDbWeight, "sp_weights RuntimeDbWeight", sc.Sequence[sc.Str]{"sp_weights", "RuntimeDbWeight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // read
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // write
			})),

		primitives.NewMetadataType(metadata.TypesRuntimeVersion, "sp_version RuntimeVersion", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString), // spec_name
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString), // impl_name
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // authoring_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // spec_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // impl_version
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeApis),     // apis
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32),    // transaction_version
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),     // state_version
			})),

		primitives.NewMetadataTypeWithParam(metadata.TypesRuntimeApis, "ApisVec = sp_std::borrow::Cow<'static, [(ApiId, u32)]>;", sc.Sequence[sc.Str]{"Cow"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeVecApis),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesRuntimeVecApis, "T"),
		),

		primitives.NewMetadataType(
			metadata.TypesRuntimeVecApis,
			"[(ApiId, u32)]",
			sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleApiIdU32))),

		primitives.NewMetadataType(
			metadata.TypesTupleApiIdU32,
			"(ApiId, u32)",
			sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionTuple(
				sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesApiId), sc.ToCompact(metadata.PrimitiveTypesU32)})),

		primitives.NewMetadataType(
			metadata.TypesApiId,
			"ApiId",
			sc.Sequence[sc.Str]{},
			primitives.NewMetadataTypeDefinitionFixedSequence(8, sc.ToCompact(metadata.PrimitiveTypesU8))),

		primitives.NewMetadataTypeWithParams(
			metadata.TypesAuraStorageAuthorites, "BoundedVec<T::AuthorityId, T::MaxAuthorities>", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequencePubKeys),
				}), sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAuthorityId, "T"),
				primitives.NewMetadataEmptyTypeParameter("S"),
			}),
		primitives.NewMetadataType(metadata.TypesAuthorityId, "sp_consensus_aura sr25519 app_sr25519 Public", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesSr25519PubKey)})),

		primitives.NewMetadataType(metadata.TypesSr25519PubKey, "sp_core sr25519 Public", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		primitives.NewMetadataType(metadata.TypesSequencePubKeys, "sequence pub keys", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesAuthorityId))),

		primitives.NewMetadataType(metadata.TypesAuraSlot, "sp_consensus_slots Slot", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64),
			})),

		primitives.NewMetadataTypeWithParams(metadata.TypesBalancesErrors, "pallet_balances pallet Error", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"VestingBalance",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Vesting balance too high to send value"),
				primitives.NewMetadataDefinitionVariant(
					"LiquidityRestrictions",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"Account liquidity restrictions prevent withdrawal"),
				primitives.NewMetadataDefinitionVariant(
					"InsufficientBalance",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					2,
					"Balance too low to send value."),
				primitives.NewMetadataDefinitionVariant(
					"ExistentialDeposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					3,
					"Value too low to create account due to existential deposit"),
				primitives.NewMetadataDefinitionVariant(
					"KeepAlive",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					4,
					"Transfer/payment would kill account"),
				primitives.NewMetadataDefinitionVariant(
					"ExistingVestingSchedule",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					5,
					"A vesting schedule already exists for this account"),
				primitives.NewMetadataDefinitionVariant(
					"DeadAccount",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					6,
					"Beneficiary account must pre-exist"),
				primitives.NewMetadataDefinitionVariant(
					"TooManyReserves",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					7,
					"Number of named reserves exceed MaxReserves"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),

		primitives.NewMetadataType(metadata.TypesFixedU128, "FixedU128", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
			})),

		primitives.NewMetadataType(metadata.TypesTransactionPaymentReleases, "Releases", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"V1Ancient",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Original version of the pallet."),
				primitives.NewMetadataDefinitionVariant(
					"V2",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"One that bumps the usage to FixedU128 from FixedI128."),
			})),

		primitives.NewMetadataTypeWithParam(metadata.TypesTransactionPaymentEvents, "pallet_transaction_payment pallet Event", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"TransactionFeePaid",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),     // who
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // actual_fee
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // fee
					},
					0,
					"Event.TransactionFeePaid"),
			}), primitives.NewMetadataEmptyTypeParameter("T")),

		primitives.NewMetadataType(metadata.TypesAddress32, "Address32", sc.Sequence[sc.Str]{"sp_core", "crypto", "AccountId32"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]")},
		)),

		primitives.NewMetadataType(metadata.TypesAccountData, "AccountData", sc.Sequence[sc.Str]{"pallet_balances", "AccountData"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "Balance"),        // Free
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "Balance"),    // Reserved
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "misc_frozen", "Balance"), // MiscFrozen
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "fee_frozen", "Balance"),  // FeeFrozen
			},
		)),
		primitives.NewMetadataType(metadata.TypesAccountInfo, "AccountInfo", sc.Sequence[sc.Str]{"frame_system", "AccountInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "nonce", "Index"),          // Nonce
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "consumers", "RefCount"),   // Consumers
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "providers", "RefCount"),   // Providers
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "sufficients", "RefCount"), // Sufficients
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAccountData, "data", "AccountData"),      // Data
			},
		)),
		primitives.NewMetadataType(metadata.TypesWeight, "Weight", sc.Sequence[sc.Str]{"sp_weights", "weight_v2", "Weight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "ref_time", "u64"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "proof_size", "u64"),
			},
		)),
		primitives.NewMetadataTypeWithParam(metadata.TypesOptionWeight, "Option<Weight>", sc.Sequence[sc.Str]{"Option"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Option<Weight>(nil)"),
				primitives.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),
					},
					1,
					"Option<Weight>(value)"),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight, "T"),
		),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeight, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "mandatory", "T"),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight, "T"),
		),
		primitives.NewMetadataType(metadata.TypesWeightPerClass, "WeightPerClass", sc.Sequence[sc.Str]{"frame_system", "limits", "WeightsPerClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_extrinsic", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_extrinsic", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "max_total", "Option<Weight>"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesOptionWeight, "reserved", "Option<Weight>"),
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeightsPerClass, "PerDispatchClass<WeightPerClass>", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "normal", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "operational", "T"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeightPerClass, "mandatory", "T"),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesWeightPerClass, "T")),
		primitives.NewMetadataType(metadata.TypesBlockWeights, "BlockWeights", sc.Sequence[sc.Str]{"frame_system", "limits", "BlockWeights"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "base_block", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "max_block", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPerDispatchClassWeightsPerClass, "per_class", "PerDispatchClass<WeightPerClass>"),
			})),
		primitives.NewMetadataType(metadata.TypesSequenceU8, "Sequence[U8]", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesDigestItem, "DigestItem", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "DigestItem"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"PreRuntime",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					6,
					"DigestItem.PreRuntime"),
				primitives.NewMetadataDefinitionVariant(
					"Consensus",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					4,
					"DigestItem.Consensus"),
				primitives.NewMetadataDefinitionVariant(
					"Seal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					5,
					"DigestItem.Seal"),
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					0,
					"DigestItem.Seal"),
				primitives.NewMetadataDefinitionVariant(
					"RuntimeEnvironmentUpdated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					8,
					"DigestItem.RuntimeEnvironmentUpdated"),
			},
		)),
		primitives.NewMetadataType(metadata.TypesDigest, "sp_runtime generic digest Digest", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "Digest"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesSliceDigestItem, "logs", "Vec<DigestItem>"),
				})),
		primitives.NewMetadataType(metadata.TypesSliceDigestItem, "Vec<DigestItem>", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesDigestItem))),
		primitives.NewMetadataTypeWithParams(metadata.TypesMultiAddress, "MultiAddress", sc.Sequence[sc.Str]{"sp_runtime", "multiaddress", "MultiAddress"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Id",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesAddress32, "AccountId"),
					},
					primitives.MultiAddressId,
					"MultiAddress.Id"),
				primitives.NewMetadataDefinitionVariant( // TODO: Check if this has sto be an empty tuple
					"Index",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU32, "AccountIndex"),
					},
					primitives.MultiAddressIndex,
					"MultiAddress.Index"),
				primitives.NewMetadataDefinitionVariant(
					"Raw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.MultiAddressRaw,
					"MultiAddress.Raw"),
				primitives.NewMetadataDefinitionVariant(
					"Address32",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]"),
					},
					primitives.MultiAddress32,
					"MultiAddress.Address32"),
				primitives.NewMetadataDefinitionVariant(
					"Address20",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence20U8, "[u8; 20]"),
					},
					primitives.MultiAddress20,
					"MultiAddress.Address20"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAddress32, "AccountId"),
				primitives.NewMetadataTypeParameter(metadata.TypesEmptyTuple, "AccountIndex"),
			}),
		primitives.NewMetadataType(metadata.CheckNonZeroSender, "CheckNonZeroSender", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_non_zero_sender", "CheckNonZeroSender"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckSpecVersion, "CheckSpecVersion", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_spec_version", "CheckSpecVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckTxVersion, "CheckTxVersion", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_tx_version", "CheckTxVersion"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckGenesis, "CheckGenesis", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_genesis", "CheckGenesis"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckMortality, "CheckMortality", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_mortality", "CheckMortality"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesEra, "Era")})),
		primitives.NewMetadataType(metadata.TypesEra, "Era", sc.Sequence[sc.Str]{"sp_runtime", "generic", "era", "Era"}, primitives.NewMetadataTypeDefinitionVariant(primitives.EraTypeDefinition())),
		primitives.NewMetadataType(metadata.CheckNonce, "CheckNonce", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_nonce", "CheckNonce"}, primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.CheckWeight, "CheckWeight", sc.Sequence[sc.Str]{"frame_system", "extensions", "check_weight", "CheckWeight"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataTypeWithParam(metadata.ChargeTransactionPayment, "ChargeTransactionPayment", sc.Sequence[sc.Str]{"pallet_transaction_payment", "ChargeTransactionPayment"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU128, "BalanceOf<T>"),
			}),
			primitives.NewMetadataEmptyTypeParameter("T"),
		),
		primitives.NewMetadataType(metadata.SignedExtra, "SignedExtra", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{
				sc.ToCompact(metadata.CheckNonZeroSender),
				sc.ToCompact(metadata.CheckSpecVersion),
				sc.ToCompact(metadata.CheckTxVersion),
				sc.ToCompact(metadata.CheckGenesis),
				sc.ToCompact(metadata.CheckMortality),
				sc.ToCompact(metadata.CheckNonce),
				sc.ToCompact(metadata.CheckWeight),
				sc.ToCompact(metadata.ChargeTransactionPayment),
			})),
		primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress, "Address"),
				primitives.NewMetadataTypeParameter(metadata.Call, "Call"),
				primitives.NewMetadataTypeParameter(metadata.MultiSignature, "Signature"),
				primitives.NewMetadataTypeParameter(metadata.SignedExtra, "Extra"),
			},
		),

		primitives.NewMetadataType(metadata.Call, "RuntimeCall", sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"System",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SystemCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<System, Runtime>"),
					},
					system.ModuleIndex,
					"Call.System"),
				primitives.NewMetadataDefinitionVariant(
					"Timestamp",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TimestampCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Timestamp, Runtime>"),
					},
					timestamp.ModuleIndex,
					"Call.Timestamp"),
				primitives.NewMetadataDefinitionVariant(
					"Grandpa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TimestampCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Grandpa, Runtime>"),
					},
					grandpa.ModuleIndex,
					"Call.Grandpa"),
				primitives.NewMetadataDefinitionVariant(
					"Balances",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.BalancesCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Balances, Runtime>"),
					},
					balances.ModuleIndex,
					"Call.Balances"),
			})),
		primitives.NewMetadataType(metadata.SignatureEd25519, "SignatureEd25519", sc.Sequence[sc.Str]{"sp_core", "ed25519", "Signature"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
		)),
		primitives.NewMetadataType(metadata.SignatureSr25519, "SignatureSr25519", sc.Sequence[sc.Str]{"sp_core", "sr25519", "Signature"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
		)),
		primitives.NewMetadataType(metadata.SignatureEcdsa, "SignatureEcdsa", sc.Sequence[sc.Str]{"sp_core", "ecdsa", "Signature"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence65U8, "[u8; 65]")},
		)),
		primitives.NewMetadataType(metadata.MultiSignature, "MultiSignature", sc.Sequence[sc.Str]{"sp_runtime", "MultiSignature"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Ed25519",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SignatureEd25519, "ed25519::Signature"),
					},
					primitives.MultiSignatureEd25519,
					"MultiSignature.Ed25519"),
				primitives.NewMetadataDefinitionVariant(
					"Sr25519",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SignatureSr25519, "sr25519::Signature"),
					},
					primitives.MultiSignatureSr25519,
					"MultiSignature.Sr25519"),
				primitives.NewMetadataDefinitionVariant(
					"Ecdsa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.SignatureEcdsa, "ecdsa::Signature"),
					},
					primitives.MultiSignatureEcdsa,
					"MultiSignature.Ecdsa"),
			})),
		primitives.NewMetadataType(metadata.TypesEmptyTuple, "EmptyTuple", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{},
		)),
		primitives.NewMetadataType(metadata.Runtime, "Runtime", sc.Sequence[sc.Str]{}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
	}
}
