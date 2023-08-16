package transaction_payment

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type Module struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[sc.U32]
	Index     sc.U8
	Config    *Config
	Constants *consts
	storage   *storage
}

func New(index sc.U8, config *Config) Module {
	return Module{
		Index:     index,
		Config:    config,
		Constants: newConstants(config.OperationalFeeMultiplier),
		storage:   newStorage(),
	}
}

func (m Module) GetIndex() sc.U8 {
	return m.Index
}

func (m Module) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (m Module) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name: "TransactionPayment",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "TransactionPayment",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"NextFeeMultiplier",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesFixedU128)),
					"NextFeeMultiplier"),
				primitives.NewMetadataModuleStorageEntry(
					"StorageVersion",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesTransactionPaymentReleases)),
					"StorageVersion"),
			},
		}),
		Call:  sc.NewOption[sc.Compact](nil),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesTransactionPaymentEvent)),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"OperationalFeeMultiplier",
				sc.ToCompact(metadata.PrimitiveTypesU8),
				sc.BytesToSequenceU8(m.Constants.OperationalFeeMultiplier.Bytes()),
				"A fee multiplier for `Operational` extrinsics to compute \"virtual tip\" to boost their  `priority` ",
			),
		},
		Error: sc.NewOption[sc.Compact](nil),
		Index: m.Index,
	}
}

func (m Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesTransactionPaymentReleases, "Releases", sc.Sequence[sc.Str]{"pallet_transaction_payment", "Releases"}, primitives.NewMetadataTypeDefinitionVariant(
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

		primitives.NewMetadataTypeWithParam(metadata.TypesTransactionPaymentEvent, "pallet_transaction_payment pallet Event", sc.Sequence[sc.Str]{"pallet_transaction_payment", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"TransactionFeePaid",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "actual_fee", "BalanceOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "tip", "BalanceOf<T>"),
					},
					0,
					"Event.TransactionFeePaid"),
			}), primitives.NewMetadataEmptyTypeParameter("T")),

		primitives.NewMetadataTypeWithParam(metadata.ChargeTransactionPayment, "ChargeTransactionPayment", sc.Sequence[sc.Str]{"pallet_transaction_payment", "ChargeTransactionPayment"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU128, "BalanceOf<T>"),
			}),
			primitives.NewMetadataEmptyTypeParameter("T"),
		),
	}
}

func (m Module) ComputeFee(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.Balance {
	return m.ComputeFeeDetails(len, info, tip).FinalFee()
}

func (m Module) ComputeFeeDetails(len sc.U32, info primitives.DispatchInfo, tip primitives.Balance) primitives.FeeDetails {
	return m.computeFeeRaw(len, info.Weight, tip, info.PaysFee, info.Class)
}

func (m Module) ComputeActualFee(len sc.U32, info primitives.DispatchInfo, postInfo primitives.PostDispatchInfo, tip primitives.Balance) primitives.Balance {
	return m.computeActualFeeDetails(len, info, postInfo, tip).FinalFee()
}

func (m Module) computeActualFeeDetails(len sc.U32, info primitives.DispatchInfo, postInfo primitives.PostDispatchInfo, tip primitives.Balance) primitives.FeeDetails {
	return m.computeFeeRaw(len, postInfo.CalcActualWeight(&info), tip, postInfo.Pays(&info), info.Class)
}

func (m Module) computeFeeRaw(len sc.U32, weight primitives.Weight, tip primitives.Balance, paysFee primitives.Pays, class primitives.DispatchClass) primitives.FeeDetails {
	if paysFee[0] == primitives.PaysYes { // TODO: type safety
		unadjustedWeightFee := m.weightToFee(weight)
		multiplier := m.storage.NextFeeMultiplier.Get()
		// Storage value is FixedU128, which is different from U128.
		// It implements a decimal fixed point number, which is `1 / VALUE`
		// Example: FixedU128, VALUE is 1_000_000_000_000_000_000.
		// FixedU64, VALUE is 1_000_000_000.
		fixedU128Div := big.NewInt(1_000_000_000_000_000_000)
		bnAdjustedWeightFee := new(big.Int).Mul(multiplier.ToBigInt(), unadjustedWeightFee.ToBigInt())
		adjustedWeightFee := sc.NewU128FromBigInt(new(big.Int).Div(bnAdjustedWeightFee, fixedU128Div)) // TODO: Create FixedU128 type

		lenFee := m.lengthToFee(len)
		baseFee := m.weightToFee(m.Config.BlockWeights.Get(class).BaseExtrinsic)

		inclusionFee := sc.NewOption[primitives.InclusionFee](primitives.NewInclusionFee(baseFee, lenFee, adjustedWeightFee))

		return primitives.FeeDetails{
			InclusionFee: inclusionFee,
			Tip:          tip,
		}
	}

	return primitives.FeeDetails{
		InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		Tip:          tip,
	}
}

func (m Module) lengthToFee(length sc.U32) primitives.Balance {
	return m.Config.LengthToFee.WeightToFee(primitives.WeightFromParts(sc.U64(length), 0))
}

func (m Module) weightToFee(weight primitives.Weight) primitives.Balance {
	cappedWeight := weight.Min(m.Config.BlockWeights.MaxBlock)

	return m.Config.WeightToFee.WeightToFee(cappedWeight)
}
