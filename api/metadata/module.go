package metadata

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/extrinsic"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "Metadata"
	apiVersion    = 1
)

const (
	resultOkIdx sc.U8 = iota
	resultErrIdx
)

const (
	optionNoneIdx sc.U8 = iota
	optionSomeIdx
)

type Module struct {
	runtimeApiModules []primitives.RuntimeApiModule
	runtimeExtrinsic  extrinsic.RuntimeExtrinsic
	memUtils          utils.WasmMemoryTranslator
}

func New(runtimeExtrinsic extrinsic.RuntimeExtrinsic, runtimeApiModules []primitives.RuntimeApiModule) Module {
	return Module{
		runtimeApiModules: runtimeApiModules,
		runtimeExtrinsic:  runtimeExtrinsic,
		memUtils:          utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// Metadata returns the metadata of the runtime.
// Returns a pointer-size of the SCALE-encoded metadata of the runtime.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-metadata-metadata)
func (m Module) Metadata() int64 {
	metadata := m.buildMetadata()

	bMetadata := sc.Sequence[sc.U8]{}

	switch metadata.Version {
	case primitives.MetadataVersion14:
		bMetadata = sc.BytesToSequenceU8(primitives.NewMetadataV14(metadata.DataV14).Bytes())
	case primitives.MetadataVersion15:
		bMetadata = sc.BytesToSequenceU8(primitives.NewMetadataV15(metadata.DataV15).Bytes())
	default:
		log.Critical("Unknown md version")
	}

	return m.memUtils.BytesToOffsetAndSize(bMetadata.Bytes())
}

func (m Module) buildMetadata() primitives.Metadata {
	constantIdsMap := make(map[string]int)

	buildConstantsMap(constantIdsMap)

	metadataTypes := append(primitiveTypes(), basicTypes()...)

	metadataTypes = append(metadataTypes, m.runtimeTypes()...)

	types, modules, extrinsic := m.runtimeExtrinsic.Metadata(constantIdsMap)

	// append types to all
	metadataTypes = append(metadataTypes, types...)

	runtimeV14Metadata := primitives.RuntimeMetadataV14{
		Types:     metadataTypes,
		Modules:   modules,
		Extrinsic: extrinsic,
		Type:      sc.ToCompact(metadata.Runtime),
	}

	return primitives.Metadata{
		Version: primitives.MetadataVersion14,
		DataV14: runtimeV14Metadata,
	}
}

// MetadataAtVersion returns the metadata of a specific version of the runtime passed as argument.
// Returns a pointer-size of the SCALE-encoded metadata of the runtime.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-metadata-metadata)
func (m Module) MetadataAtVersion(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	version, err := sc.DecodeU32(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	constantIdsMap := make(map[string]int)

	buildConstantsMap(constantIdsMap)

	metadataTypes := append(primitiveTypes(), basicTypes()...)

	metadataTypes = append(metadataTypes, m.runtimeTypes()...)

	switch version {
	case sc.U32(primitives.MetadataVersion14):
		types, modules, extrinsicV14 := m.runtimeExtrinsic.Metadata(constantIdsMap)
		metadataTypes = append(metadataTypes, types...)
		metadataV14 := primitives.RuntimeMetadataV14{
			Types:     metadataTypes,
			Modules:   modules,
			Extrinsic: extrinsicV14,
			Type:      sc.ToCompact(metadata.Runtime),
		}
		bMetadataV14 := sc.BytesToSequenceU8(primitives.NewMetadataV14(metadataV14).Bytes())
		optionMd := sc.Option[sc.Sequence[sc.U8]]{
			HasValue: sc.Bool(true),
			Value:    bMetadataV14,
		}
		return m.memUtils.BytesToOffsetAndSize(optionMd.Bytes())
	case sc.U32(primitives.MetadataVersion15):
		typesV15, modulesV15, extrinsicV15, outerEnums, custom := m.runtimeExtrinsic.MetadataLatest(constantIdsMap)
		metadataTypes = append(metadataTypes, typesV15...)
		metadataV15 := primitives.RuntimeMetadataV15{
			Types:      metadataTypes,
			Modules:    modulesV15,
			Extrinsic:  extrinsicV15,
			Type:       sc.ToCompact(metadata.Runtime),
			Apis:       m.runtimeApiMetadata(),
			OuterEnums: outerEnums,
			Custom:     custom,
		}
		bMetadataV15 := sc.BytesToSequenceU8(primitives.NewMetadataV15(metadataV15).Bytes())
		optionMd := sc.Option[sc.Sequence[sc.U8]]{
			HasValue: sc.Bool(true),
			Value:    bMetadataV15,
		}
		return m.memUtils.BytesToOffsetAndSize(optionMd.Bytes())
	default:
		optionUnsupported := sc.Option[sc.Sequence[sc.U8]]{
			HasValue: sc.Bool(false),
		}
		return m.memUtils.BytesToOffsetAndSize(optionUnsupported.Bytes())
	}
}

func (m Module) MetadataVersions() int64 {
	bVersions := sc.Sequence[sc.U32]{
		sc.U32(primitives.MetadataVersion14), sc.U32(primitives.MetadataVersion15),
	}

	return m.memUtils.BytesToOffsetAndSize(bVersions.Bytes())
}

func (m Module) runtimeApiMetadata() sc.Sequence[primitives.RuntimeApiMetadata] {
	runtimeApiMetadata := sc.Sequence[primitives.RuntimeApiMetadata]{}

	for _, module := range m.runtimeApiModules {
		runtimeApiMetadata = append(runtimeApiMetadata, module.Metadata())
	}

	return append(runtimeApiMetadata, m.apiMetadata())
}

func (m Module) apiMetadata() primitives.RuntimeApiMetadata {
	modules := sc.Sequence[primitives.RuntimeApiMethodMetadata]{
		primitives.RuntimeApiMethodMetadata{
			Name:   "metadata",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(metadata.TypesOpaqueMetadata),
			Docs:   sc.Sequence[sc.Str]{" Returns the metadata of a runtime."},
		},
		primitives.RuntimeApiMethodMetadata{
			Name: "metadata_at_version",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
				primitives.RuntimeApiMethodParamMetadata{
					Name: "version",
					Type: sc.ToCompact(metadata.PrimitiveTypesU32),
				},
			},
			Output: sc.ToCompact(metadata.TypeOptionOpaqueMetadata),
			Docs: sc.Sequence[sc.Str]{" Returns the metadata at a given version.",
				"",
				" If the given `version` isn't supported, this will return `None`.",
				" Use [`Self::metadata_versions`] to find out about supported metadata version of the runtime."},
		},
		primitives.RuntimeApiMethodMetadata{
			Name:   "metadata_versions",
			Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(metadata.TypesSequenceU32),
			Docs: sc.Sequence[sc.Str]{" Returns the supported metadata versions.",
				"",
				" This can be used to call `metadata_at_version`."},
		},
	}

	return primitives.RuntimeApiMetadata{
		Name:    ApiModuleName,
		Methods: modules,
		Docs:    sc.Sequence[sc.Str]{" The `Metadata` api trait that returns metadata for the runtime."},
	}
}

func buildConstantsMap(constantIdsMap map[string]int) {
	constantIdsMap["Bool"] = metadata.PrimitiveTypesBool
	constantIdsMap["String"] = metadata.PrimitiveTypesString
	constantIdsMap["U8"] = metadata.PrimitiveTypesU8
	constantIdsMap["U16"] = metadata.PrimitiveTypesU16
	constantIdsMap["U32"] = metadata.PrimitiveTypesU32
	constantIdsMap["U64"] = metadata.PrimitiveTypesU64
	constantIdsMap["U128"] = metadata.PrimitiveTypesU128
	constantIdsMap["U256"] = metadata.PrimitiveTypesU256
	constantIdsMap["I8"] = metadata.PrimitiveTypesI8
	constantIdsMap["I16"] = metadata.PrimitiveTypesI16
	constantIdsMap["I32"] = metadata.PrimitiveTypesI32
	constantIdsMap["I64"] = metadata.PrimitiveTypesI64
	constantIdsMap["I128"] = metadata.PrimitiveTypesI128
}

// primitiveTypes returns all primitive types
func primitiveTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.PrimitiveTypesBool, "bool", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveBoolean)),
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
	}
}

func basicTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.TypesFixedSequence4U8, "[4]byte", primitives.NewMetadataTypeDefinitionFixedSequence(4, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence8U8, "[8]byte", primitives.NewMetadataTypeDefinitionFixedSequence(8, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence20U8, "[20]byte", primitives.NewMetadataTypeDefinitionFixedSequence(20, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence32U8, "[32]byte", primitives.NewMetadataTypeDefinitionFixedSequence(32, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence64U8, "[64]byte", primitives.NewMetadataTypeDefinitionFixedSequence(64, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesFixedSequence65U8, "[65]byte", primitives.NewMetadataTypeDefinitionFixedSequence(65, sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesSequenceU8, "[]byte", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU8))),
		primitives.NewMetadataType(metadata.TypesSequenceU32, "[]uint32", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.TypesCompactU32, "CompactU32", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32))),
		primitives.NewMetadataType(metadata.TypesCompactU64, "CompactU64", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU64))),
		primitives.NewMetadataType(metadata.TypesCompactU128, "CompactU128", primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU128))),

		primitives.NewMetadataType(metadata.TypesFixedU128, "FixedU128", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU128, "u128"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		primitives.NewMetadataTypeWithPath(metadata.TypesAddress32, "Address32", sc.Sequence[sc.Str]{"sp_core", "crypto", "AccountId32"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]")},
		)),

		primitives.NewMetadataTypeWithPath(metadata.TypesKeyTypeId, "KeyTypeId", sc.Sequence[sc.Str]{"sp_core", "crypto", "KeyTypeId"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "[u8; 4]")},
		)),

		primitives.NewMetadataTypeWithPath(metadata.TypesAccountData, "AccountData", sc.Sequence[sc.Str]{"pallet_balances", "AccountData"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "misc_frozen", "Balance"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "fee_frozen", "Balance"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesAccountInfo, "AccountInfo", sc.Sequence[sc.Str]{"frame_system", "AccountInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "nonce", "Index"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "consumers", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "providers", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "sufficients", "RefCount"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAccountData, "data", "AccountData"),
			},
		)),

		primitives.NewMetadataTypeWithPath(metadata.TypesWeight, "Weight", sc.Sequence[sc.Str]{"sp_weights", "weight_v2", "Weight"}, primitives.NewMetadataTypeDefinitionComposite(
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
		primitives.NewMetadataTypeWithParam(metadata.TypesPerDispatchClassU32,
			"PerDispatchClass[U32]",
			sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "normal", "T"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "operational", "T"),
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU32, "mandatory", "T"),
				},
			),
			primitives.NewMetadataTypeParameter(metadata.PrimitiveTypesU32, "T"),
		),

		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureEd25519, "SignatureEd25519", sc.Sequence[sc.Str]{"sp_core", "ed25519", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesEd25519PubKey, "SignatureEd25519 Public", sc.Sequence[sc.Str]{"sp_core", "ed25519", "Public"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence32U8, "[u8; 32]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureSr25519, "SignatureSr25519", sc.Sequence[sc.Str]{"sp_core", "sr25519", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence64U8, "[u8; 64]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesSignatureEcdsa, "SignatureEcdsa", sc.Sequence[sc.Str]{"sp_core", "ecdsa", "Signature"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence65U8, "[u8; 65]")},
			)),
		primitives.NewMetadataTypeWithPath(metadata.TypesMultiSignature, "MultiSignature", sc.Sequence[sc.Str]{"sp_runtime", "MultiSignature"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Ed25519",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureEd25519, "ed25519::Signature"),
						},
						primitives.MultiSignatureEd25519,
						"MultiSignature.Ed25519"),
					primitives.NewMetadataDefinitionVariant(
						"Sr25519",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureSr25519, "sr25519::Signature"),
						},
						primitives.MultiSignatureSr25519,
						"MultiSignature.Sr25519"),
					primitives.NewMetadataDefinitionVariant(
						"Ecdsa",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSignatureEcdsa, "ecdsa::Signature"),
						},
						primitives.MultiSignatureEcdsa,
						"MultiSignature.Ecdsa"),
				})),

		primitives.NewMetadataType(metadata.TypesEmptyTuple, "EmptyTuple", primitives.NewMetadataTypeDefinitionTuple(
			sc.Sequence[sc.Compact]{},
		)),

		primitives.NewMetadataType(metadata.TypesTupleU32U32, "(U32, U32)",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.PrimitiveTypesU32), sc.ToCompact(metadata.PrimitiveTypesU32)})),

		primitives.NewMetadataTypeWithParams(metadata.TypesMultiAddress, "MultiAddress", sc.Sequence[sc.Str]{"sp_runtime", "multiaddress", "MultiAddress"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Id",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesAddress32, "AccountId"),
					},
					primitives.MultiAddressId,
					"MultiAddress.Id"),
				primitives.NewMetadataDefinitionVariant(
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

		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchInfo, "DispatchInfo", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchInfo"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesWeight, "weight", "Weight"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesDispatchClass, "class", "DispatchClass"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesPays, "pays_fee", "Pays"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchClass, "DispatchClass", sc.Sequence[sc.Str]{"frame_support", "dispatch", "DispatchClass"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Normal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					sc.U8(primitives.DispatchClassNormal),
					"DispatchClass.Normal"),
				primitives.NewMetadataDefinitionVariant(
					"Operational",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					sc.U8(primitives.DispatchClassOperational),
					"DispatchClass.Operational"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					sc.U8(primitives.DispatchClassMandatory),
					"DispatchClass.Mandatory"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesPays, "Pays", sc.Sequence[sc.Str]{"frame_support", "dispatch", "Pays"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Yes",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					sc.U8(primitives.PaysYes),
					"Pays.Yes"),
				primitives.NewMetadataDefinitionVariant(
					"No",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					sc.U8(primitives.PaysNo),
					"Pays.No"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesDispatchError, "DispatchError", sc.Sequence[sc.Str]{"sp_runtime", "DispatchError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorOther,
					"DispatchError.Other"),
				primitives.NewMetadataDefinitionVariant(
					"Cannotlookup",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorCannotLookup,
					"DispatchError.Cannotlookup"),
				primitives.NewMetadataDefinitionVariant(
					"BadOrigin",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorBadOrigin,
					"DispatchError.BadOrigin"),
				primitives.NewMetadataDefinitionVariant(
					"Module",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesModuleError, "ModuleError"),
					},
					primitives.DispatchErrorModule,
					"DispatchError.Module"),
				primitives.NewMetadataDefinitionVariant(
					"ConsumerRemaining",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorConsumerRemaining,
					"DispatchError.ConsumerRemaining"),
				primitives.NewMetadataDefinitionVariant(
					"NoProviders",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorNoProviders,
					"DispatchError.NoProviders"),
				primitives.NewMetadataDefinitionVariant(
					"TooManyConsumers",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorTooManyConsumers,
					"DispatchError.TooManyConsumers"),
				primitives.NewMetadataDefinitionVariant(
					"Token",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTokenError, "TokenError"),
					},
					primitives.DispatchErrorToken,
					"DispatchError.Token"),
				primitives.NewMetadataDefinitionVariant(
					"Arithmetic",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesArithmeticError, "ArithmeticError"),
					},
					primitives.DispatchErrorArithmetic,
					"DispatchError.Arithmetic"),
				primitives.NewMetadataDefinitionVariant(
					"TransactionalError",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesTransactionalError, "TransactionalError"),
					},
					primitives.DispatchErrorTransactional,
					"DispatchError.TransactionalError"),
				primitives.NewMetadataDefinitionVariant(
					"Exhausted",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorExhausted,
					"DispatchError.Exhausted"),
				primitives.NewMetadataDefinitionVariant(
					"Corruption",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorCorruption,
					"DispatchError.Corruption"),
				primitives.NewMetadataDefinitionVariant(
					"Unavailable",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DispatchErrorUnavailable,
					"DispatchError.Unavailable"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesModuleError, "ModuleError", sc.Sequence[sc.Str]{"sp_runtime", "ModuleError"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU8, "index", "u8"),
				primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesFixedSequence4U8, "error", "[u8; MAX_MODULE_ERROR_ENCODED_SIZE]"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.TypesTokenError, "TokenError", sc.Sequence[sc.Str]{"sp_runtime", "TokenError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"NoFunds",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorNoFunds,
					"TokenError.NoFunds"),
				primitives.NewMetadataDefinitionVariant(
					"WouldDie",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorWouldDie,
					"TokenError.WouldDie"),
				primitives.NewMetadataDefinitionVariant(
					"Mandatory",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorBelowMinimum,
					"TokenError.BelowMinimum"),
				primitives.NewMetadataDefinitionVariant(
					"CannotCreate",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorCannotCreate,
					"TokenError.CannotCreate"),
				primitives.NewMetadataDefinitionVariant(
					"UnknownAsset",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorUnknownAsset,
					"TokenError.UnknownAsset"),
				primitives.NewMetadataDefinitionVariant(
					"Frozen",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorFrozen,
					"TokenError.Frozen"),
				primitives.NewMetadataDefinitionVariant(
					"Unsupported",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TokenErrorUnsupported,
					"TokenError.Unsupported"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesArithmeticError, "ArithmeticError", sc.Sequence[sc.Str]{"sp_arithmetic", "ArithmeticError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Underflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorUnderflow,
					"ArithmeticError.Underflow"),
				primitives.NewMetadataDefinitionVariant(
					"Overflow",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorOverflow,
					"ArithmeticError.Overflow"),
				primitives.NewMetadataDefinitionVariant(
					"DivisionByZero",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.ArithmeticErrorDivisionByZero,
					"ArithmeticError.DivisionByZero"),
			})),
		primitives.NewMetadataTypeWithPath(metadata.TypesTransactionalError, "TransactionalError", sc.Sequence[sc.Str]{"sp_runtime", "TransactionalError"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"LimitReached",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TransactionalErrorLimitReached,
					"TransactionalError.LimitReached"),
				primitives.NewMetadataDefinitionVariant(
					"NoLayer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.TransactionalErrorNoLayer,
					"TransactionalError.NoLayer"),
			})),

		primitives.NewMetadataType(metadata.TypesVecTopics, "Vec<Topics>", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesH256))),

		primitives.NewMetadataTypeWithPath(metadata.TypesDigestItem, "DigestItem", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "DigestItem"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"PreRuntime",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestItemPreRuntime,
					"DigestItem.PreRuntime"),
				primitives.NewMetadataDefinitionVariant(
					"Consensus",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestItemConsensusMessage,
					"DigestItem.Consensus"),
				primitives.NewMetadataDefinitionVariant(
					"Seal",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesFixedSequence4U8, "ConsensusEngineId"),
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestItemSeal,
					"DigestItem.Seal"),
				primitives.NewMetadataDefinitionVariant(
					"Other",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
					},
					primitives.DigestItemOther,
					"DigestItem.Other"),
				primitives.NewMetadataDefinitionVariant(
					"RuntimeEnvironmentUpdated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					primitives.DigestItemRuntimeEnvironmentUpgraded,
					"DigestItem.RuntimeEnvironmentUpdated"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesDigest, "sp_runtime generic digest Digest", sc.Sequence[sc.Str]{"sp_runtime", "generic", "digest", "Digest"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesSliceDigestItem, "logs", "Vec<DigestItem>"),
				})),
		primitives.NewMetadataType(metadata.TypesSliceDigestItem, "Vec<DigestItem>", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesDigestItem))),

		primitives.NewMetadataTypeWithParams(metadata.Header, "Header",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "header", "Header"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "Hash::Output"), // parent_hash
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesCompactU32, "Number"), // number
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "Hash::Output"), // state_root
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "Hash::Output"), // extrinsics_root
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesDigest, "Digest"),     // digest
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.PrimitiveTypesU32, "Number"),
				primitives.NewMetadataEmptyTypeParameter("Hash"),
			},
		),

		primitives.NewMetadataTypeWithPath(metadata.TypesOpaqueMetadata,
			"sp_core OpaqueMetadata",
			sc.Sequence[sc.Str]{"sp_core", "OpaqueMetadata"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesSequenceU8, "Vec<u8>"),
				})),

		primitives.NewMetadataTypeWithParam(metadata.TypeOptionOpaqueMetadata, "Option<OpaqueMetadata>", sc.Sequence[sc.Str]{"Option"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					optionNoneIdx,
					""),
				primitives.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesOpaqueMetadata),
					},
					optionSomeIdx,
					""),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesOpaqueMetadata, "T")),

		// type 31
		primitives.NewMetadataTypeWithParams(metadata.TypesResultEmptyTuple, "Result", sc.Sequence[sc.Str]{"Result"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Ok",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesEmptyTuple),
					},
					resultOkIdx,
					""),
				primitives.NewMetadataDefinitionVariant(
					"Err",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesDispatchError),
					},
					resultErrIdx, ""),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesEmptyTuple, "T"),
				primitives.NewMetadataTypeParameter(metadata.TypesDispatchError, "E")}),

		// type 869
		primitives.NewMetadataTypeWithParams(metadata.TypesResult, "Result", sc.Sequence[sc.Str]{"Result"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Ok",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesResultEmptyTuple),
					},
					resultOkIdx,
					""),
				primitives.NewMetadataDefinitionVariant(
					"Err",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesTransactionValidityError),
					},
					resultErrIdx, ""),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesResultEmptyTuple, "T"),
				primitives.NewMetadataTypeParameter(metadata.TypesTransactionValidityError, "E"),
			}),
		primitives.NewMetadataType(metadata.TypesSequenceUncheckedExtrinsics, "[]byte", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.UncheckedExtrinsic))),
		//type 876
		primitives.NewMetadataType(metadata.TypesTuple8U8SequenceU8, "([8]bytes, []byte])",
			primitives.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesFixedSequence8U8), sc.ToCompact(metadata.TypesSequenceU8)})),
		// type 875
		primitives.NewMetadataType(metadata.TypesSequenceTuple8U8SequenceU8, "[]byte", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTuple8U8SequenceU8))),
		// type 874
		primitives.NewMetadataTypeWithParams(metadata.TypesBTreeMap, "BTreeMap",
			sc.Sequence[sc.Str]{"BTreeMap"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceTuple8U8SequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesFixedSequence8U8, "K"),
				primitives.NewMetadataTypeParameter(metadata.TypesSequenceU8, "V"),
			},
		),
		primitives.NewMetadataTypeWithPath(metadata.TypesInherentData, "sp_inherents InherentData", sc.Sequence[sc.Str]{"sp_inherents", "InherentData"}, primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesBTreeMap, "BTreeMap<InherentIdentifier, Vec<u8>>"),
			})),

		primitives.NewMetadataTypeWithPath(metadata.CheckInherentsResult, "sp_inherents CheckInherentsResult", sc.Sequence[sc.Str]{"sp_inherents", "CheckInherentsResult"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "bool"),
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "bool"),
				primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesInherentData, "InherentData"),
			},
			)),
		primitives.NewMetadataTypeWithParam(metadata.TypesOptionSequenceU8, "Option<Seq<U8>>", sc.Sequence[sc.Str]{"Option"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{},
					optionNoneIdx,
					""),
				primitives.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
					},
					optionSomeIdx,
					""),
			}),
			primitives.NewMetadataTypeParameter(metadata.TypesSequenceU8, "T")),
		primitives.NewMetadataType(metadata.TypesSequenceSequenceU8, "[][]byte", primitives.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesSequenceU8))),
	}
}

func (m Module) runtimeTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesRuntimeVersion, "sp_version RuntimeVersion", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, primitives.NewMetadataTypeDefinitionComposite(
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
		primitives.NewMetadataType(metadata.Runtime, "Runtime", primitives.NewMetadataTypeDefinitionComposite(
			sc.Sequence[primitives.MetadataTypeDefinitionField]{})),
	}
}
