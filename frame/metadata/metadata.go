package metadata

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/constants/metadata"
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
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassWeight, "PerDispatchClass[Weight]", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Normal
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Operational
				primitives.NewMetadataTypeDefinitionField(metadata.TypesWeight), // Mandatory
			},
		),
			primitives.NewMetadataTypeParameter(metadata.TypesWeight),
		),
		primitives.NewMetadataType(metadata.TypesSequenceU8, "Sequence[U8]", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesDigestItem, "DigestItem", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence4U8), // Engine
				primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),       // Payload
			},
		)),
		primitives.NewMetadataType(metadata.TypesSliceDigestItem, "[]DigestItem", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesDigestItem))),
		primitives.NewMetadataType(metadata.TypesMultiAddress, "MultiAddress", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Id",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8),
					},
					primitives.MultiAddressId,
					"MultiAddress.Id"),
				primitives.NewMetadataDefinitionVariant(
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
			})),
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
