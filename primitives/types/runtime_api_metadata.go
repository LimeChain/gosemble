package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/api/account_nonce"
	blockbuilder "github.com/LimeChain/gosemble/api/block_builder"
	"github.com/LimeChain/gosemble/api/core"
	"github.com/LimeChain/gosemble/api/grandpa"
	md "github.com/LimeChain/gosemble/api/metadata"
	ow "github.com/LimeChain/gosemble/api/offchain_worker"
	"github.com/LimeChain/gosemble/api/session_keys"
	"github.com/LimeChain/gosemble/api/tagged_transaction_queue"
	"github.com/LimeChain/gosemble/api/transaction_payment"
	"github.com/LimeChain/gosemble/api/transaction_payment_call"
	mdconstants "github.com/LimeChain/gosemble/constants/metadata"
)

const (
	applyExtrinsicMethod      = "apply_extrinsic"
	finalizeBlockMethod       = "finalize_block"
	inherentExtrinsicsMethod  = "inherent_extrinsics"
	checkInherentsMethod      = "check_inherents"
	metadataMethod            = "metadata"
	metadataAtVersionMethod   = "metadata_at_version"
	metadataVersionsMethod    = "metadata_versions"
	coreVersionMethod         = "version"
	validateTransactionMethod = "validate_transaction"
	offChainWorkerMethod      = "offchain_worker"
	grandpaAuthoritiesMethod  = "grandpa_authorities"
	accountNonceMethod        = "account_nonce"
	queryInfoMethod           = "query_info"
	queryFeeDetailsMethod     = "query_fee_details"
	queryCallInfoMethod       = "query_call_info"
	queryCallFeeDetailsMethod = "query_call_fee_details"
	generateSessionKeysMethod = "generate_session_keys"
	decodeSessionKeysMethod   = "decode_session_keys"
)

type RuntimeApiMetadata struct {
	Name    sc.Str
	Methods sc.Sequence[RuntimeApiMethodMetadata]
	Docs    sc.Sequence[sc.Str]
}

func (ram RuntimeApiMetadata) Encode(buffer *bytes.Buffer) {
	ram.Name.Encode(buffer)
	ram.Methods.Encode(buffer)
	ram.Docs.Encode(buffer)
}

func DecodeRuntimeApiMetadata(buffer *bytes.Buffer) RuntimeApiMetadata {
	return RuntimeApiMetadata{
		Name:    sc.DecodeStr(buffer),
		Methods: sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodMetadata),
		Docs:    sc.DecodeSequence[sc.Str](buffer),
	}
}

func (ram RuntimeApiMetadata) Bytes() []byte {
	return sc.EncodedBytes(ram)
}

func ApiMetadata() sc.Sequence[RuntimeApiMetadata] {
	coreMethodsMd := coreMethodsMd()
	metadataMethodsMd := metadataMethodsMd()
	blockbuilderMethodsMd := blockBuilderMethodsMd()
	taggedTransactionQueueMethodsMd := taggedTransactionQueueMethodsMd()
	offChainWorkerMethodsMd := offChainWorkerMethodsMd()
	grandpaMethodsMd := grandpaMethodsMd()
	accountNonceMethodsMd := accountNonceMethodsMd()
	transactionPaymentMethodsMd := transactionPaymentMethodsMd()
	transactionPaymentCallMethodsMd := transactionPaymentCallMethodsMd()
	sessionKeysMethodsMd := sessionKeysMethodsMd()
	return sc.Sequence[RuntimeApiMetadata]{
		RuntimeApiMetadata{
			Name:    core.ApiModuleName,
			Methods: coreMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `Core` runtime api that every Substrate runtime needs to implement."},
		},
		RuntimeApiMetadata{
			Name:    md.ApiModuleName,
			Methods: metadataMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `Metadata` api trait that returns metadata for the runtime."},
		},
		RuntimeApiMetadata{
			Name:    blockbuilder.ApiModuleName,
			Methods: blockbuilderMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `BlockBuilder` api trait that provides the required functionality for building a block."},
		},
		RuntimeApiMetadata{
			Name:    tagged_transaction_queue.ApiModuleName,
			Methods: taggedTransactionQueueMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `TaggedTransactionQueue` api trait for interfering with the transaction queue."},
		},
		RuntimeApiMetadata{
			Name:    ow.ApiModuleName,
			Methods: offChainWorkerMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The offchain worker api."},
		},
		RuntimeApiMetadata{
			Name:    grandpa.ApiModuleName,
			Methods: grandpaMethodsMd,
			Docs: sc.Sequence[sc.Str]{
				" APIs for integrating the GRANDPA finality gadget into runtimes.",
				" This should be implemented on the runtime side.",
				"",
				" This is primarily used for negotiating authority-set changes for the",
				" gadget. GRANDPA uses a signaling model of changing authority sets:",
				" changes should be signaled with a delay of N blocks, and then automatically",
				" applied in the runtime after those N blocks have passed.",
				"",
				" The consensus protocol will coordinate the handoff externally.",
			},
		},
		RuntimeApiMetadata{
			Name:    account_nonce.ApiModuleName,
			Methods: accountNonceMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The API to query account nonce."},
		},
		RuntimeApiMetadata{
			Name:    transaction_payment.ApiModuleName,
			Methods: transactionPaymentMethodsMd,
			Docs:    sc.Sequence[sc.Str]{},
		},
		RuntimeApiMetadata{
			Name:    transaction_payment_call.ApiModuleName,
			Methods: transactionPaymentCallMethodsMd,
			Docs:    sc.Sequence[sc.Str]{},
		},
		RuntimeApiMetadata{
			Name:    session_keys.ApiModuleName,
			Methods: sessionKeysMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" Session keys runtime api."},
		},
	}
}

func coreMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	versionOutput := NewMetadataTypeWithPath(mdconstants.TypesRuntimeVersion, "sp_version RuntimeVersion", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesString), // spec_name
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesString), // impl_name
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesU32),    // authoring_version
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesU32),    // spec_version
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesU32),    // impl_version
			NewMetadataTypeDefinitionField(mdconstants.TypesRuntimeApis),     // apis
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesU32),    // transaction_version
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesU8),     // state_version
		}))
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   coreVersionMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(versionOutput),
			Docs:   sc.Sequence[sc.Str]{" Returns the version of the runtime."},
		},
	}
}

func metadataMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	metadataOutput := NewMetadataTypeWithPath(mdconstants.TypesOpaqueMetadata, // TODO: Verify it's correct
		"sp_core OpaqueMetadata",
		sc.Sequence[sc.Str]{"sp_core", "OpaqueMetadata"},
		NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(mdconstants.TypesSequenceU8, "Vec<u8>"),
			}))

	metadataDocs := sc.Str(" Returns the metadata of a runtime.")

	metadataAtVersionOutput := NewMetadataTypeWithParam(mdconstants.TypesOptionWeight, "Option", sc.Sequence[sc.Str]{"Option"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"None",
				sc.Sequence[MetadataTypeDefinitionField]{},
				0,
				""),
			NewMetadataDefinitionVariant(
				"Some",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(mdconstants.TypesOpaqueMetadata),
				},
				1,
				""),
		}),
		NewMetadataTypeParameter(mdconstants.TypesOpaqueMetadata, "T"),
	)

	metadataVersionsOutput := NewMetadataType(mdconstants.TypesSequenceU32, "[]byte", NewMetadataTypeDefinitionSequence(sc.ToCompact(mdconstants.PrimitiveTypesU32)))

	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   metadataMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(metadataOutput),
			Docs:   sc.Sequence[sc.Str]{metadataDocs},
		},
		RuntimeApiMethodMetadata{
			Name:   metadataAtVersionMethod,
			Inputs: metadataAtVersionsInputsMd(),
			Output: sc.ToCompact(metadataAtVersionOutput),
			Docs: sc.Sequence[sc.Str]{" Returns the metadata at a given version.",
				"",
				" If the given `version` isn't supported, this will return `None`.",
				" Use [`Self::metadata_versions`] to find out about supported metadata version of the runtime."},
		},
		RuntimeApiMethodMetadata{
			Name:   metadataVersionsMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(metadataVersionsOutput),
			Docs: sc.Sequence[sc.Str]{" Returns the supported metadata versions.",
				"",
				" This can be used to call `metadata_at_version`."},
		},
	}
}

func metadataAtVersionsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "version",
			Type: sc.ToCompact(NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))),
		},
	}
}

func blockBuilderMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	applyExtrinsicOutput := NewMetadataTypeWithParams(mdconstants.TypesResult, "Result", sc.Sequence[sc.Str]{"Result"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"Ok",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(mdconstants.TypesEmptyResult),
				},
				mdconstants.TypesResultOk,
				""),
			NewMetadataDefinitionVariant(
				"Err",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(mdconstants.TypesTransactionalError), // TODO: Is this the correct constant ?
				},
				mdconstants.TypesResultErr, ""),
		}),
		sc.Sequence[MetadataTypeParameter]{
			NewMetadataTypeParameter(mdconstants.TypesEmptyResult, "T"),
			NewMetadataTypeParameter(mdconstants.TypesTransactionalError, "E"), // TODO: Is this the correct constant ?
		})

	finalizeBlockOutput := getHeaderType()

	inherentExtrinsicsOutput := NewMetadataType(mdconstants.TypesSequenceUncheckedExtrinsics, "[]byte", NewMetadataTypeDefinitionSequence(sc.ToCompact(mdconstants.UncheckedExtrinsic)))

	checkInherentsOutput := NewMetadataTypeWithPath(mdconstants.CheckInherentsResult, "sp_inherents CheckInherentsResult", sc.Sequence[sc.Str]{"sp_inherents", "CheckInherentsResult"},
		NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(mdconstants.TypesResultOk, "bool"),
			NewMetadataTypeDefinitionFieldWithName(mdconstants.TypesResultErr, "bool"),
			NewMetadataTypeDefinitionFieldWithName(mdconstants.TypesInherentData, "inherentData"),
		},
		))

	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   applyExtrinsicMethod,
			Inputs: applyExtrinsicInputsMd(),
			Output: sc.ToCompact(applyExtrinsicOutput),
			Docs: sc.Sequence[sc.Str]{" Apply the given extrinsic.",
				"",
				" Returns an inclusion outcome which specifies if this extrinsic is included in",
				" this block or not."},
		},
		RuntimeApiMethodMetadata{
			Name:   finalizeBlockMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(finalizeBlockOutput),
			Docs:   sc.Sequence[sc.Str]{" Finish the current block."},
		},
		RuntimeApiMethodMetadata{
			Name:   inherentExtrinsicsMethod,
			Inputs: inherentExtrinsicsInputsMd(),
			Output: sc.ToCompact(inherentExtrinsicsOutput),
			Docs:   sc.Sequence[sc.Str]{" Generate inherent extrinsics. The inherent data will vary from chain to chain."},
		},
		RuntimeApiMethodMetadata{
			Name:   checkInherentsMethod,
			Inputs: checkInherentsInputsMd(),
			Output: sc.ToCompact(checkInherentsOutput),
			Docs:   sc.Sequence[sc.Str]{" Check that the inherents are valid. The inherent data will vary from chain to chain."},
		},
	}
}

func applyExtrinsicInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	uncheckedExtrinsicType := getUncheckedExtrinsicType()
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "Extrinsic",
			Type: sc.ToCompact(uncheckedExtrinsicType),
		},
	}
}

func inherentExtrinsicsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	inherentType := NewMetadataTypeWithPath(mdconstants.TypesRuntimeVersion, "sp_inherents InherentData", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesString), // data TODO: Encode the BTreeMap
		}))
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "inherent",
			Type: sc.ToCompact(inherentType),
		},
	}
}

func checkInherentsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {

	//t := NewMetadataTypeWithPath(mdconstants.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
	//	NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
	//		NewMetadataTypeDefinitionField(mdconstants.TypesFixedSequence32U8)}))

	blockType := NewMetadataTypeWithParams(mdconstants.UncheckedExtrinsic, "UncheckedExtrinsic",
		sc.Sequence[sc.Str]{"sp_runtime", "generic", "block", "Block"},
		NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionField(mdconstants.TypesH256),       // parent_hash
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // number
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // state_root
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // extrinsics_root
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // digest
			}),
		sc.Sequence[MetadataTypeParameter]{
			NewMetadataTypeParameter(mdconstants.Header, "Header"), // TODO: Is this correct ?
			NewMetadataTypeParameter(mdconstants.RuntimeCall, "Extrinsic"),
		},
	)

	inherentType := NewMetadataTypeWithPath(mdconstants.TypesRuntimeVersion, "sp_inherents InherentData", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesString), // data TODO: Encode the BTreeMap
		}))
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "block",
			Type: sc.ToCompact(blockType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "data",
			Type: sc.ToCompact(inherentType),
		},
	}
}

func taggedTransactionQueueMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   validateTransactionMethod,
			Inputs: validateTransactionInputsMd(),
			Output: sc.ToCompact(5), // TODO
			Docs: sc.Sequence[sc.Str]{" Validate the transaction.",
				"",
				" This method is invoked by the transaction pool to learn details about given transaction.",
				" The implementation should make sure to verify the correctness of the transaction",
				" against current state. The given `block_hash` corresponds to the hash of the block",
				" that is used as current state.",
				"",
				" Note that this call may be performed by the pool multiple times and transactions",
				" might be verified in any possible order."},
		},
	}
}

func validateTransactionInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	sourceType := NewMetadataTypeWithPath(mdconstants.TypesTransactionSource, "TransactionSource", sc.Sequence[sc.Str]{"sp_runtime", "transaction_validity", "TransactionSource"},
		NewMetadataTypeDefinitionVariant(
			sc.Sequence[MetadataDefinitionVariant]{
				NewMetadataDefinitionVariant(
					"InBlock",
					sc.Sequence[MetadataTypeDefinitionField]{},
					TransactionSourceInBlock,
					"TransactionSourceInBlock"),
				NewMetadataDefinitionVariant(
					"Local",
					sc.Sequence[MetadataTypeDefinitionField]{},
					TransactionSourceLocal,
					"TransactionSourceLocal"),
				NewMetadataDefinitionVariant(
					"External",
					sc.Sequence[MetadataTypeDefinitionField]{},
					TransactionSourceExternal,
					"TransactionSourceExternal"),
			}))
	uncheckedExtrinsicType := getUncheckedExtrinsicType()

	blockHashType := NewMetadataTypeWithPath(mdconstants.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
		NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionField(mdconstants.TypesFixedSequence32U8)}))
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "source",
			Type: sc.ToCompact(sourceType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "tx",
			Type: sc.ToCompact(uncheckedExtrinsicType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "block_hash",
			Type: sc.ToCompact(blockHashType),
		},
	}
}

func offChainWorkerMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	emptyTupleType := NewMetadataType(mdconstants.TypesEmptyTuple, "EmptyTuple", NewMetadataTypeDefinitionTuple(
		sc.Sequence[sc.Compact]{},
	))
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   offChainWorkerMethod,
			Inputs: offChainWorkerInputsMd(),
			Output: sc.ToCompact(emptyTupleType),
			Docs:   sc.Sequence[sc.Str]{" Starts the off-chain task for given block header."},
		},
	}
}

func offChainWorkerInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	//t := NewMetadataTypeWithPath(mdconstants.TypesH256, "primitives H256", sc.Sequence[sc.Str]{"primitive_types", "H256"},
	//	NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
	//		NewMetadataTypeDefinitionField(mdconstants.TypesFixedSequence32U8)}))

	headerType := getHeaderType()

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "header",
			Type: sc.ToCompact(headerType),
		},
	}
}

func grandpaMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   grandpaAuthoritiesMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(5), // TODO: Encode
			Docs: sc.Sequence[sc.Str]{
				" Get the current GRANDPA authorities and weights. This should not change except",
				" for when changes are scheduled and the corresponding delay has passed.",
				"",
				" When called at block B, it will return the set of authorities that should be",
				" used to finalize descendants of this block (B+1, B+2, ...). The block B itself",
				" is finalized by the authorities from block B-1.",
			},
		},
	}
}

func accountNonceMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   accountNonceMethod,
			Inputs: accountNonceInputsMd(),
			Output: sc.ToCompact(u32Type),
			Docs:   sc.Sequence[sc.Str]{" Get current account nonce of given `AccountId`."},
		},
	}
}

func accountNonceInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	accountType := NewMetadataTypeWithPath(mdconstants.TypesAddress32, "Address32", sc.Sequence[sc.Str]{"sp_core", "crypto", "AccountId32"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{NewMetadataTypeDefinitionFieldWithName(mdconstants.TypesFixedSequence32U8, "[u8; 32]")},
	))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "account",
			Type: sc.ToCompact(accountType),
		},
	}
}

func transactionPaymentMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   queryInfoMethod,
			Inputs: queryInfoInputsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
		RuntimeApiMethodMetadata{
			Name:   queryFeeDetailsMethod,
			Inputs: queryFeeDetailsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
	}
}

func queryInfoInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	uncheckedExtrinsicType := getUncheckedExtrinsicType()
	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "uxt",
			Type: sc.ToCompact(uncheckedExtrinsicType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type),
		},
	}
}

func queryFeeDetailsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	uncheckedExtrinsicType := getUncheckedExtrinsicType()

	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "uxt",
			Type: sc.ToCompact(uncheckedExtrinsicType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type),
		},
	}
}

func transactionPaymentCallMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   queryCallInfoMethod,
			Inputs: queryCallInfoInputsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
		RuntimeApiMethodMetadata{
			Name:   queryCallFeeDetailsMethod,
			Inputs: queryCallFeeDetailsInputsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
	}
}

func queryCallInfoInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	// TODO: A MetadataType with 61 variants ?
	callType := NewMetadataTypeWithPath(mdconstants.RuntimeCall, "RuntimeCall", sc.Sequence[sc.Str]{"kitchensink_runtime", "RuntimeCall"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"System",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorOther,
				"DispatchError.Other"),
			NewMetadataDefinitionVariant(
				"Utility",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorCannotLookup,
				"DispatchError.Cannotlookup"),
			NewMetadataDefinitionVariant(
				"Babe",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorBadOrigin,
				"DispatchError.BadOrigin"),
		}))
	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "call",
			Type: sc.ToCompact(callType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type),
		},
	}
}

func queryCallFeeDetailsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))

	// TODO: A MetadataType with 61 variants ?
	callType := NewMetadataTypeWithPath(mdconstants.RuntimeCall, "RuntimeCall", sc.Sequence[sc.Str]{"kitchensink_runtime", "RuntimeCall"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"System",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorOther,
				"DispatchError.Other"),
			NewMetadataDefinitionVariant(
				"Utility",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorCannotLookup,
				"DispatchError.Cannotlookup"),
			NewMetadataDefinitionVariant(
				"Babe",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorBadOrigin,
				"DispatchError.BadOrigin"),
		}))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "call",
			Type: sc.ToCompact(callType),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type),
		},
	}
}

func sessionKeysMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   generateSessionKeysMethod,
			Inputs: generateSessionKeysInputsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
		RuntimeApiMethodMetadata{
			Name:   decodeSessionKeysMethod,
			Inputs: decodeSessionKeysInputsMd(),
			Output: sc.ToCompact(5),
			Docs:   sc.Sequence[sc.Str]{}, // TODO: Add docs
		},
	}
}

func generateSessionKeysInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	seedType := NewMetadataTypeWithParam(mdconstants.TypesOptionWeight, "Option", sc.Sequence[sc.Str]{"Option"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"None",
				sc.Sequence[MetadataTypeDefinitionField]{},
				0,
				""),
			NewMetadataDefinitionVariant(
				"Some",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8),
				},
				1,
				""),
		}),
		NewMetadataTypeParameter(mdconstants.TypesSequenceU8, "T"),
	)

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "seed",
			Type: sc.ToCompact(seedType),
		},
	}
}

func decodeSessionKeysInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	encodedType := NewMetadataType(mdconstants.TypesSequenceU8, "[]byte", NewMetadataTypeDefinitionSequence(sc.ToCompact(mdconstants.PrimitiveTypesU8)))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "encoded",
			Type: sc.ToCompact(encodedType),
		},
	}
}

// A helper function since this type is often used
func getUncheckedExtrinsicType() MetadataType {
	return NewMetadataTypeWithParams(mdconstants.UncheckedExtrinsic, "UncheckedExtrinsic",
		sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
		NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8),
			}),
		sc.Sequence[MetadataTypeParameter]{
			NewMetadataTypeParameter(mdconstants.TypesMultiAddress, "Address"),
			NewMetadataTypeParameter(mdconstants.RuntimeCall, "Call"), // TODO: Is this correct ?
			NewMetadataTypeParameter(mdconstants.TypesMultiSignature, "Signature"),
			NewMetadataTypeParameter(mdconstants.SignedExtra, "Extra"),
		},
	)
}

func getHeaderType() MetadataType {
	return NewMetadataTypeWithParams(mdconstants.TypesOffChainWorker, "OffChainWorker",
		sc.Sequence[sc.Str]{"sp_runtime", "generic", "header", "Header"},
		NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionField(mdconstants.TypesH256),       // parent_hash // TODO: Is this correct ?
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // number
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // state_root
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // extrinsics_root
				NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8), // digest
			}),
		sc.Sequence[MetadataTypeParameter]{
			NewMetadataTypeParameter(mdconstants.PrimitiveTypesU64, "Number"), // TODO: Is this correct ?
			NewMetadataTypeParameter(mdconstants.TypesH256, "Hash"),
		},
	)
}
