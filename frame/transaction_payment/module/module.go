package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/constants/transaction_payment"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TransactionPaymentModule struct {
}

func NewTransactionPaymentModule() TransactionPaymentModule {
	return TransactionPaymentModule{}
}

func (tpm TransactionPaymentModule) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (tpm TransactionPaymentModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tpm TransactionPaymentModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (tpm TransactionPaymentModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return sc.Sequence[primitives.MetadataType]{}, primitives.MetadataModule{
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
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesTransactionPaymentEvents)),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"OperationalFeeMultiplier",
				sc.ToCompact(metadata.PrimitiveTypesU8),
				sc.BytesToSequenceU8(sc.U8(5).Bytes()),
				"A fee multiplier for `Operational` extrinsics to compute \"virtual tip\" to boost their  `priority` ",
			),
		},
		Error: sc.NewOption[sc.Compact](nil),
		Index: transaction_payment.ModuleIndex,
	}
}
