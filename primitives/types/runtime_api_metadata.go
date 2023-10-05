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
	coreExecuteBlockMethod    = "execute_block"
	coreInitializeBlockMethod = "initialize_block"
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
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   coreVersionMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(mdconstants.TypesRuntimeVersion),
			Docs:   sc.Sequence[sc.Str]{" Returns the version of the runtime."},
		},
		RuntimeApiMethodMetadata{
			Name: coreExecuteBlockMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{
				RuntimeApiMethodParamMetadata{
					Name: "block",
					Type: sc.ToCompact(mdconstants.TypesBlock),
				},
			},
			Output: sc.ToCompact(mdconstants.TypesEmptyTuple),
			Docs:   sc.Sequence[sc.Str]{" Execute the given block."},
		},
		RuntimeApiMethodMetadata{
			Name: coreInitializeBlockMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{
				RuntimeApiMethodParamMetadata{
					Name: "header",
					Type: sc.ToCompact(mdconstants.Header), // TODO: Is this correct?
				},
			},
			Output: sc.ToCompact(mdconstants.TypesEmptyTuple),
			Docs:   sc.Sequence[sc.Str]{" Initialize a block with the given header."},
		},
	}
}

func metadataMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   metadataMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(mdconstants.TypesOpaqueMetadata),
			Docs:   sc.Sequence[sc.Str]{" Returns the metadata of a runtime."},
		},
		RuntimeApiMethodMetadata{
			Name:   metadataAtVersionMethod,
			Inputs: metadataAtVersionsInputsMd(),
			Output: sc.ToCompact(mdconstants.TypeOption),
			Docs: sc.Sequence[sc.Str]{" Returns the metadata at a given version.",
				"",
				" If the given `version` isn't supported, this will return `None`.",
				" Use [`Self::metadata_versions`] to find out about supported metadata version of the runtime."},
		},
		RuntimeApiMethodMetadata{
			Name:   metadataVersionsMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(mdconstants.TypesSequenceU32),
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
			Type: sc.ToCompact(mdconstants.PrimitiveTypesU32),
		},
	}
}

func blockBuilderMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	applyExtrinsicOutput_869 := NewMetadataTypeWithParams(mdconstants.TypesResult, "Result", sc.Sequence[sc.Str]{"Result"}, NewMetadataTypeDefinitionVariant(
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
			Output: sc.ToCompact(applyExtrinsicOutput.Id),
			Docs: sc.Sequence[sc.Str]{" Apply the given extrinsic.",
				"",
				" Returns an inclusion outcome which specifies if this extrinsic is included in",
				" this block or not."},
		},
		RuntimeApiMethodMetadata{
			Name:   finalizeBlockMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(mdconstants.Header),
			Docs:   sc.Sequence[sc.Str]{" Finish the current block."},
		},
		RuntimeApiMethodMetadata{
			Name:   inherentExtrinsicsMethod,
			Inputs: inherentExtrinsicsInputsMd(),
			Output: sc.ToCompact(inherentExtrinsicsOutput.Id),
			Docs:   sc.Sequence[sc.Str]{" Generate inherent extrinsics. The inherent data will vary from chain to chain."},
		},
		RuntimeApiMethodMetadata{
			Name:   checkInherentsMethod,
			Inputs: checkInherentsInputsMd(),
			Output: sc.ToCompact(checkInherentsOutput.Id),
			Docs:   sc.Sequence[sc.Str]{" Check that the inherents are valid. The inherent data will vary from chain to chain."},
		},
	}
}

func applyExtrinsicInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "Extrinsic",
			Type: sc.ToCompact(mdconstants.UncheckedExtrinsic),
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
			Type: sc.ToCompact(inherentType.Id),
		},
	}
}

func checkInherentsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	inherentType := NewMetadataTypeWithPath(mdconstants.TypesRuntimeVersion, "sp_inherents InherentData", sc.Sequence[sc.Str]{"sp_version", "RuntimeVersion"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionField(mdconstants.PrimitiveTypesString), // TODO: Encode the BTreeMap
		}))
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "block",
			Type: sc.ToCompact(mdconstants.TypesBlock),
		},
		RuntimeApiMethodParamMetadata{
			Name: "data",
			Type: sc.ToCompact(inherentType.Id),
		},
	}
}

func taggedTransactionQueueMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   validateTransactionMethod,
			Inputs: validateTransactionInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesResultValidityTransaction),
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
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "source",
			Type: sc.ToCompact(mdconstants.TypesTransactionSource),
		},
		RuntimeApiMethodParamMetadata{
			Name: "tx",
			Type: sc.ToCompact(mdconstants.UncheckedExtrinsic),
		},
		RuntimeApiMethodParamMetadata{
			Name: "block_hash",
			Type: sc.ToCompact(mdconstants.TypesH256),
		},
	}
}

func offChainWorkerMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   offChainWorkerMethod,
			Inputs: offChainWorkerInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesEmptyTuple),
			Docs:   sc.Sequence[sc.Str]{" Starts the off-chain task for given block header."},
		},
	}
}

func offChainWorkerInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "header",
			Type: sc.ToCompact(mdconstants.Header),
		},
	}
}

func grandpaMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   grandpaAuthoritiesMethod,
			Inputs: sc.Sequence[RuntimeApiMethodParamMetadata]{},
			Output: sc.ToCompact(mdconstants.TypesSequenceTupleGrandpaAppPublic),
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
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   accountNonceMethod,
			Inputs: accountNonceInputsMd(),
			Output: sc.ToCompact(mdconstants.PrimitiveTypesU32),
			Docs:   sc.Sequence[sc.Str]{" Get current account nonce of given `AccountId`."},
		},
	}
}

func accountNonceInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "account",
			Type: sc.ToCompact(mdconstants.TypesAddress32),
		},
	}
}

func transactionPaymentMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   queryInfoMethod,
			Inputs: queryInfoInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesTransactionPaymentRuntimeDispatchInfo),
			Docs:   sc.Sequence[sc.Str]{},
		},
		RuntimeApiMethodMetadata{
			Name:   queryFeeDetailsMethod,
			Inputs: queryFeeDetailsMd(),
			Output: sc.ToCompact(mdconstants.TypesTransactionPaymentFeeDetails),
			Docs:   sc.Sequence[sc.Str]{},
		},
	}
}

func queryInfoInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "uxt",
			Type: sc.ToCompact(mdconstants.UncheckedExtrinsic),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(mdconstants.PrimitiveTypesU32),
		},
	}
}

func queryFeeDetailsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "uxt",
			Type: sc.ToCompact(mdconstants.UncheckedExtrinsic),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(mdconstants.PrimitiveTypesU32),
		},
	}
}

func transactionPaymentCallMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   queryCallInfoMethod,
			Inputs: queryCallInfoInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesTransactionPaymentRuntimeDispatchInfo),
			Docs:   sc.Sequence[sc.Str]{" Query information of a dispatch class, weight, and fee of a given encoded `Call`."},
		},
		RuntimeApiMethodMetadata{
			Name:   queryCallFeeDetailsMethod,
			Inputs: queryCallFeeDetailsInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesTransactionPaymentFeeDetails),
			Docs:   sc.Sequence[sc.Str]{" Query fee details of a given encoded `Call`."},
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
			Type: sc.ToCompact(callType.Id),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type.Id),
		},
	}
}

func queryCallFeeDetailsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	u32Type := NewMetadataType(mdconstants.PrimitiveTypesU32, "U32", NewMetadataTypeDefinitionPrimitive(MetadataDefinitionPrimitiveU32))

	callType := NewMetadataTypeWithPath(mdconstants.RuntimeCall, "RuntimeCall", sc.Sequence[sc.Str]{"kitchensink_runtime", "RuntimeCall"}, NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"System",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(mdconstants.TypesSequenceU8),
				},
				DispatchErrorOther,
				""),
			NewMetadataDefinitionVariant(
				"Timestamp",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorBadOrigin,
				""),
			NewMetadataDefinitionVariant(
				"Balances",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorBadOrigin,
				""),
			NewMetadataDefinitionVariant(
				"Grandpa",
				sc.Sequence[MetadataTypeDefinitionField]{},
				DispatchErrorBadOrigin,
				""),
		}))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "call",
			Type: sc.ToCompact(callType.Id),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(u32Type.Id),
		},
	}
}

func sessionKeysMethodsMd() sc.Sequence[RuntimeApiMethodMetadata] {
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   generateSessionKeysMethod,
			Inputs: generateSessionKeysInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesSequenceU8),
			Docs: sc.Sequence[sc.Str]{
				" Generate a set of session keys with optionally using the given seed.",
				" The keys should be stored within the keystore exposed via runtime",
				" externalities.",
				"",
				" The seed needs to be a valid `utf8` string.",
				"",
				" Returns the concatenated SCALE encoded public keys.",
			},
		},
		RuntimeApiMethodMetadata{
			Name:   decodeSessionKeysMethod,
			Inputs: decodeSessionKeysInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesOptionTupleSequenceU8KeyTypeId),
			Docs: sc.Sequence[sc.Str]{
				" Decode the given public session keys.",
				"",
				" Returns the list of public raw public keys + key type.",
			},
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
			Type: sc.ToCompact(seedType.Id),
		},
	}
}

func decodeSessionKeysInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	encodedType := NewMetadataType(mdconstants.TypesSequenceU8, "[]byte", NewMetadataTypeDefinitionSequence(sc.ToCompact(mdconstants.PrimitiveTypesU8)))

	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "encoded",
			Type: sc.ToCompact(encodedType.Id),
		},
	}
}
