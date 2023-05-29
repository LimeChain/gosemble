package metadata

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/constants/system"
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
			primitives.NewMetadataSignedExtension("CheckGenesis", metadata.CheckGenesis, metadata.TypesFixedSequence32U8),
			primitives.NewMetadataSignedExtension("CheckMortality", metadata.CheckMortality, metadata.TypesFixedSequence32U8),
			primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple),
			primitives.NewMetadataSignedExtension("CheckWeight", metadata.CheckWeight, metadata.TypesEmptyTuple),
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
		primitives.NewMetadataType(metadata.PrimitiveTypesBool, "bool", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveBoolean)),
		primitives.NewMetadataType(metadata.PrimitiveTypesChar, "char", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveChar)),
		primitives.NewMetadataType(metadata.PrimitiveTypesString, "string", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveString)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU8, "U8", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU16, "U16", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU32, "U32", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU64, "U64", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU128, "U128", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesU256, "U256", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveU256)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI8, "I8", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI8)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI16, "I16", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI16)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI32, "I32", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI32)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI64, "I64", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI64)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI128, "I128", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI128)),
		primitives.NewMetadataType(metadata.PrimitiveTypesI256, "I256", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveI256)),
	}
}

func runtimeTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.TypesFixedSequence4U8, "[4]byte", primitives.NewMetadataTypeDefinitionFixedSequence(4, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence20U8, "[20]byte", primitives.NewMetadataTypeDefinitionFixedSequence(20, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence32U8, "[32]byte", primitives.NewMetadataTypeDefinitionFixedSequence(32, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence64U8, "[64]byte", primitives.NewMetadataTypeDefinitionFixedSequence(64, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence65U8, "[65]byte", primitives.NewMetadataTypeDefinitionFixedSequence(65, sc.ToCompact(metadata.PrimitiveTypesU8))),

		primitives.NewMetadataType(metadata.TypesCompactU32, "CompactU32", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.TypesCompactU64, "CompactU64", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))),
		primitives.NewMetadataType(metadata.TypesCompactU128, "CompactU128", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))),

		primitives.NewMetadataType(metadata.TypesH256, "primitives H256",
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),
		primitives.NewMetadataType(metadata.TypesVecBlockNumEventIndex, "Vec<BlockNumber, EventIndex>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleU32U32))),
		primitives.NewMetadataType(metadata.TypesTupleU32U32, "(U32, U32)",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.PrimitiveTypesU32), sc.ToCompact(metadata.PrimitiveTypesU32)})),
		primitives.NewMetadataType(metadata.TypesSystemEventStorage, "Vec<Box<EventRecord<T::RuntimeEvent, T::Hash>>>",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesEventRecord))),
		primitives.NewMetadataType(metadata.TypesEventRecord, "frame_system EventRecord",
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesPhase),        // phase
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeEvent), // event
				primitives.NewMetadataTypeDefinitionField(metadata.TypesVecTopics),    // topics
			})),
		primitives.NewMetadataType(metadata.TypesPhase, "frame_system Phase", primitives.NewMetadataTypeDefinitionVariant(
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
		primitives.NewMetadataType(metadata.TypesRuntimeEvent, "node_template_runtime RuntimeEvent", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"System",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesSystemEvent),
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
						primitives.NewMetadataTypeDefinitionField(metadata.TypesBalancesEvent),
					},
					balances.ModuleIndex,
					"Events.Balances"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionPayment",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesTransactionPaymentEvents),
					},
					transaction_payment.ModuleIndex,
					"Events.TransactionPayment"),
			})),
		primitives.NewMetadataType(metadata.TypesSystemEvent, "frame_system pallet Event", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"ExtrinsicSuccess",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesDispatchInfo),
					},
					0,
					"Event.ExtrinsicSuccess"),
				primitives.NewMetadataDefinitionVariant(
					"ExtrinsicFailed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesDispatchError),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesDispatchInfo),
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
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
					},
					3,
					"Events.NewAccount"),
				primitives.NewMetadataDefinitionVariant(
					"KilledAccount",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
					},
					4,
					"Events.KilledAccount"),
				primitives.NewMetadataDefinitionVariant(
					"Remarked",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesH256),
					},
					5,
					"Events.Remarked"),
			})),

		primitives.NewMetadataType(metadata.TypesDispatchInfo, "DispatchInfo", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),
				primitives.NewMetadataTypeDefinitionField(metadata.TypesDispatchClass),
				primitives.NewMetadataTypeDefinitionField(metadata.TypesPays),
			},
		)),
		primitives.NewMetadataType(metadata.TypesDispatchClass, "DispatchClass", primitives.NewMetadataTypeDefinitionVariant(
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
		primitives.NewMetadataType(metadata.TypesPays, "Pays", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Yes",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					0,
					"Pays.Yes"),
				primitives.NewMetadataDefinitionVariant(
					"Operational",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					1,
					"Pays.No"),
			})),

		primitives.NewMetadataType(metadata.TypesDispatchError, "DispatchError", primitives.NewMetadataTypeDefinitionVariant(
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
						primitives.NewMetadataTypeDefinitionField(metadata.TypesModuleError),
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
						primitives.NewMetadataTypeDefinitionField(metadata.TypesTokenError),
					},
					7,
					"DispatchError.Token"),
				primitives.NewMetadataDefinitionVariant(
					"Arithmetic",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesArithmeticError),
					},
					8,
					"DispatchError.Arithmetic"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionalError",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesTransactionalError),
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
		primitives.NewMetadataType(metadata.TypesModuleError, "ModuleError", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence4U8),
			})),
		primitives.NewMetadataType(metadata.TypesTokenError, "TokenError", primitives.NewMetadataTypeDefinitionVariant(
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
		primitives.NewMetadataType(metadata.TypesArithmeticError, "ArithmeticError", primitives.NewMetadataTypeDefinitionVariant(
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
		primitives.NewMetadataType(metadata.TypesTransactionalError, "TransactionalError", primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataType(metadata.TypesBalancesEvent, "pallet_balances pallet Event", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Endowed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					0,
					"Event.Endowed"),
				primitives.NewMetadataDefinitionVariant(
					"DustLost",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					1,
					"Events.DustLost"),
				primitives.NewMetadataDefinitionVariant(
					"Transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					2,
					"Events.Transfer"),
				primitives.NewMetadataDefinitionVariant(
					"BalanceSet",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					3,
					"Events.BalanceSet"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					4,
					"Events.Reserved"),
				primitives.NewMetadataDefinitionVariant(
					"Unreserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					5,
					"Events.Unreserved"),
				primitives.NewMetadataDefinitionVariant(
					"ReserveRepatriated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
						primitives.NewMetadataTypeDefinitionField(metadata.TypesBalanceStatus),
					},
					6,
					"Events.ReserveRepatriated"),
				primitives.NewMetadataDefinitionVariant(
					"Deposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					7,
					"Event.Deposit"),
				primitives.NewMetadataDefinitionVariant(
					"Withdraw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					8,
					"Event.Withdraw"),
				primitives.NewMetadataDefinitionVariant(
					"Slashed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
						primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
					},
					9,
					"Event.Slashed"),
			},
		)),

		primitives.NewMetadataType(metadata.TypesBalanceStatus, "BalanceStatus", primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataType(metadata.TypesVecTopics, "Vec<Topics>", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesH256))),

		primitives.NewMetadataType(metadata.TypesLastRuntimeUpgradeInfo, "LastRuntimeUpgradeInfo", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU32),
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesString),
			})),

		primitives.NewMetadataType(metadata.TypesSystemErrors, "frame_system pallet Error", primitives.NewMetadataTypeDefinitionVariant(
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
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesPerDispatchClassU32), // max
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassU32, "PerDispatchClass[U32]", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Normal
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Operational
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Mandatory
			},
		),
			primitives.NewMetadataTypeParameter(metadata.PrimitiveTypesU32),
		),

		primitives.NewMetadataType(metadata.TypesDbWeight, "sp_weights RuntimeDbWeight", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // read
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64), // write
			})),

		primitives.NewMetadataType(metadata.TypesRuntimeVersion, "sp_version RuntimeVersion", primitives.NewMetadataTypeDefinitionComposite(
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

		primitives.NewMetadataTypeWithParam(metadata.TypesRuntimeApis, "ApisVec = sp_std::borrow::Cow<'static, [(ApiId, u32)]>;", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesRuntimeVecApis),
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesRuntimeVecApis),
		),

		primitives.NewMetadataType(
			metadata.TypesRuntimeVecApis,
			"[(ApiId, u32)]",
			primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleApiIdU32))),

		primitives.NewMetadataType(
			metadata.TypesTupleApiIdU32,
			"(ApiId, u32)",
			primitives.NewMetadataTypeDefinitionTuple(
				sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesApiId), sc.ToCompact(metadata.PrimitiveTypesU32)})),

		primitives.NewMetadataType(
			metadata.TypesApiId,
			"ApiId",
			primitives.NewMetadataTypeDefinitionFixedSequence(8, sc.ToCompact(metadata.PrimitiveTypesU8))),

		primitives.NewMetadataTypeWithParams(
			metadata.TypesAuraStorageAuthorites, "BoundedVec<T::AuthorityId, T::MaxAuthorities>", primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequencePubKeys),
				}), sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAuthorityId),
				primitives.NewMetadataEmptyTypeParameter("S"),
			}),
		primitives.NewMetadataType(metadata.TypesAuthorityId, "sp_consensus_aura sr25519 app_sr25519 Public", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesSr25519PubKey)})),

		primitives.NewMetadataType(metadata.TypesSr25519PubKey, "sp_core sr25519 Public", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		primitives.NewMetadataType(metadata.TypesSequencePubKeys, "sequence pub keys", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesAuthorityId))),

		primitives.NewMetadataType(metadata.TypesAuraSlot, "sp_consensus_slots Slot", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64),
			})),

		primitives.NewMetadataTypeWithParams(metadata.TypesBalancesErrors, "pallet_balances pallet Error", primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataType(metadata.TypesFixedU128, "FixedU128", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128),
			})),

		primitives.NewMetadataType(metadata.TypesTransactionPaymentReleases, "Releases", primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataTypeWithParam(metadata.TypesTransactionPaymentEvents, "pallet_transaction_payment pallet Event", primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataType(metadata.TypesAddress32, "Address32", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)},
		)),

		primitives.NewMetadataType(metadata.TypesAccountData, "AccountData", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // Free
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // Reserved
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // MiscFrozen
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU128), // FeeFrozen
			},
		)),
		primitives.NewMetadataType(metadata.TypesAccountInfo, "AccountInfo", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Nonce
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Consumers
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Providers
				primitives.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU32), // Sufficients
				primitives.NewMetadataTypeDefinitionField(metadata.TypesAccountData),  // Data
			},
		)),
		primitives.NewMetadataType(metadata.TypesWeight, "Weight", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU64), // RefTime
				primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU64), // ProofSize
			},
		)),
		primitives.NewMetadataType(metadata.TypesOptionWeight, "Option<Weight>", primitives.NewMetadataTypeDefinitionVariant(
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
			})),
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeight, "PerDispatchClass[Weight]", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Normal
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Operational
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Mandatory
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight),
		),
		primitives.NewMetadataType(metadata.TypesWeightPerClass, "WeightPerClass", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),       // BaseExtrinsic
				primitives.NewMetadataTypeDefinitionField(metadata.TypesOptionWeight), // MaxExtrinsic
				primitives.NewMetadataTypeDefinitionField(metadata.TypesOptionWeight), // MaxTotal
				primitives.NewMetadataTypeDefinitionField(metadata.TypesOptionWeight), // Reserved
			})),
		primitives.NewMetadataType(metadata.TypesPerDispatchClassWeightsPerClass, "PerDispatchClass[WeightPerClass]", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeightPerClass), // Normal
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeightPerClass), // Operational
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeightPerClass), // Mandatory
			})),
		primitives.NewMetadataType(metadata.TypesBlockWeights, "BlockWeights", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),                          // BaseBlock
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight),                          // MaxBlock
				primitives.NewMetadataTypeDefinitionField(metadata.TypesPerDispatchClassWeightsPerClass), // PerClass
			})),
		primitives.NewMetadataType(metadata.TypesSequenceU8, "Sequence[U8]", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesDigestItem, "DigestItem", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence4U8), // Engine
				primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),       // Payload
			},
		)),
		primitives.NewMetadataType(metadata.TypesSliceDigestItem, "[]DigestItem", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesDigestItem))),
		primitives.NewMetadataTypeWithParams(metadata.TypesMultiAddress, "MultiAddress", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Id",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8),
					},
					primitives.MultiAddressId,
					"MultiAddress.Id"),
				primitives.NewMetadataDefinitionVariant( // TODO: Check if this has sto be an empty tuple
					"Index",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesCompactU32),
					},
					primitives.MultiAddressIndex,
					"MultiAddress.Index"),
				primitives.NewMetadataDefinitionVariant(
					"Raw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
					},
					primitives.MultiAddressRaw,
					"MultiAddress.Raw"),
				primitives.NewMetadataDefinitionVariant(
					"Address32",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesAddress32),
					},
					primitives.MultiAddress32,
					"MultiAddress.Address32"),
				primitives.NewMetadataDefinitionVariant(
					"Address20",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence20U8),
					},
					primitives.MultiAddress20,
					"MultiAddress.Address20"),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesAddress32),
				primitives.NewMetadataTypeParameter(metadata.TypesEmptyTuple),
			}),
		primitives.NewMetadataType(metadata.CheckNonZeroSender, "CheckNonZeroSender", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckSpecVersion, "CheckSpecVersion", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckTxVersion, "CheckTxVersion", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckGenesis, "CheckGenesis", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.CheckMortality, "CheckMortality", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				// TODO: Add Era
			})),
		primitives.NewMetadataType(metadata.CheckNonce, "CheckNonce", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.CheckWeight, "CheckWeight", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
		primitives.NewMetadataType(metadata.SignedExtra, "SignedExtra", primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{
				sc.ToCompact(metadata.CheckNonZeroSender),
				sc.ToCompact(metadata.CheckSpecVersion),
				sc.ToCompact(metadata.CheckTxVersion),
				sc.ToCompact(metadata.CheckGenesis),
				sc.ToCompact(metadata.CheckMortality),
				sc.ToCompact(metadata.CheckNonce),
				sc.ToCompact(metadata.CheckWeight),
			})),
		primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress),
				primitives.NewMetadataTypeParameter(metadata.Call),
				primitives.NewMetadataTypeParameter(metadata.MultiSignature),
				primitives.NewMetadataTypeParameter(metadata.SignedExtra),
			},
		),
		primitives.NewMetadataType(metadata.SignatureEd25519, "SignatureEd25519", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence64U8)},
		)),
		primitives.NewMetadataType(metadata.SignatureSr25519, "SignatureSr25519", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence64U8)},
		)),
		primitives.NewMetadataType(metadata.SignatureEcdsa, "SignatureEcdsa", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence65U8)},
		)),
		primitives.NewMetadataType(metadata.MultiSignature, "MultiSignature", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Ed25519",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.SignatureEd25519),
					},
					primitives.MultiSignatureEd25519,
					"MultiSignature.Ed25519"),
				primitives.NewMetadataDefinitionVariant(
					"Sr25519",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.SignatureSr25519),
					},
					primitives.MultiSignatureSr25519,
					"MultiSignature.Sr25519"),
				primitives.NewMetadataDefinitionVariant(
					"Ecdsa",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.SignatureEcdsa),
					},
					primitives.MultiSignatureEcdsa,
					"MultiSignature.Ecdsa"),
			})),
		primitives.NewMetadataType(metadata.TypesEmptyTuple, "EmptyTuple", primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{},
		)),
		primitives.NewMetadataType(metadata.Runtime, "Runtime", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
	}
}
