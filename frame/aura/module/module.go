package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/metadata"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type AuraModule struct {
}

func NewAuraModule() AuraModule {
	return AuraModule{}
}

func (am AuraModule) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (am AuraModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (am AuraModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (am AuraModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return sc.Sequence[primitives.MetadataType]{}, primitives.MetadataModule{
		Name: "Aura",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "Aura",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"Authorities",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraStorageAuthorites)),
					"The current authority set."),
				primitives.NewMetadataModuleStorageEntry(
					"CurrentSlot",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraSlot)),
					"The current slot of this block.   This will be set in `on_initialize`."),
			},
		}),
		Call:      sc.NewOption[sc.Compact](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     aura.ModuleIndex,
	}
}
