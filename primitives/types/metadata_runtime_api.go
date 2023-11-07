package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	mdconstants "github.com/LimeChain/gosemble/constants/metadata"
)

const (
	coreModuleApi                   = "Core"
	metadataModuleApi               = "Metadata"
	blockBuilderModuleApi           = "BlockBuilder"
	taggedTransactionQueueModuleApi = "TaggedTransactionQueue"
	offChainWorkerModuleApi         = "OffchainWorkerApi"
	grandpaModuleApi                = "GrandpaApi"
	accountNonceModuleApi           = "AccountNonceApi"
	transactionPaymentModuleApi     = "TransactionPaymentApi"
	transactionPaymentCallModuleApi = "TransactionPaymentCallApi"
	sessionKeysModuleApi            = "SessionKeys"
	applyExtrinsicMethod            = "apply_extrinsic"
	finalizeBlockMethod             = "finalize_block"
	inherentExtrinsicsMethod        = "inherent_extrinsics"
	checkInherentsMethod            = "check_inherents"
	metadataMethod                  = "metadata"
	metadataAtVersionMethod         = "metadata_at_version"
	metadataVersionsMethod          = "metadata_versions"
	coreVersionMethod               = "version"
	coreExecuteBlockMethod          = "execute_block"
	coreInitializeBlockMethod       = "initialize_block"
	validateTransactionMethod       = "validate_transaction"
	offChainWorkerMethod            = "offchain_worker"
	grandpaAuthoritiesMethod        = "grandpa_authorities"
	accountNonceMethod              = "account_nonce"
	queryInfoMethod                 = "query_info"
	queryFeeDetailsMethod           = "query_fee_details"
	queryCallInfoMethod             = "query_call_info"
	queryCallFeeDetailsMethod       = "query_call_fee_details"
	generateSessionKeysMethod       = "generate_session_keys"
	decodeSessionKeysMethod         = "decode_session_keys"
)

type RuntimeApiMetadata struct {
	Name    sc.Str
	Methods sc.Sequence[RuntimeApiMethodMetadata]
	Docs    sc.Sequence[sc.Str]
}

func (ram RuntimeApiMetadata) Encode(buffer *bytes.Buffer) error {
	err := ram.Name.Encode(buffer)
	if err != nil {
		return err
	}
	err = ram.Methods.Encode(buffer)
	if err != nil {
		return err
	}
	return ram.Docs.Encode(buffer)
}

func DecodeRuntimeApiMetadata(buffer *bytes.Buffer) (RuntimeApiMetadata, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	methods, err := sc.DecodeSequenceWith(buffer, DecodeRuntimeApiMethodMetadata)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return RuntimeApiMetadata{}, err
	}
	return RuntimeApiMetadata{
		Name:    name,
		Methods: methods,
		Docs:    docs,
	}, nil
}

func (ram RuntimeApiMetadata) Bytes() []byte {
	return sc.EncodedBytes(ram)
}

func ApiMetadata() sc.Sequence[RuntimeApiMetadata] {
	coreMethodsMd := coreMethodsMd()
	metadataMethodsMd := metadataMethodsMd()
	blockBuilderMethodsMd := blockBuilderMethodsMd()
	taggedTransactionQueueMethodsMd := taggedTransactionQueueMethodsMd()
	offChainWorkerMethodsMd := offChainWorkerMethodsMd()
	grandpaMethodsMd := grandpaMethodsMd()
	accountNonceMethodsMd := accountNonceMethodsMd()
	transactionPaymentMethodsMd := transactionPaymentMethodsMd()
	transactionPaymentCallMethodsMd := transactionPaymentCallMethodsMd()
	sessionKeysMethodsMd := sessionKeysMethodsMd()

	return sc.Sequence[RuntimeApiMetadata]{
		RuntimeApiMetadata{
			Name:    coreModuleApi,
			Methods: coreMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `Core` runtime api that every Substrate runtime needs to implement."},
		},
		RuntimeApiMetadata{
			Name:    metadataModuleApi,
			Methods: metadataMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `Metadata` api trait that returns metadata for the runtime."},
		},
		RuntimeApiMetadata{
			Name:    blockBuilderModuleApi,
			Methods: blockBuilderMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `BlockBuilder` api trait that provides the required functionality for building a block."},
		},
		RuntimeApiMetadata{
			Name:    taggedTransactionQueueModuleApi,
			Methods: taggedTransactionQueueMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The `TaggedTransactionQueue` api trait for interfering with the transaction queue."},
		},
		RuntimeApiMetadata{
			Name:    offChainWorkerModuleApi,
			Methods: offChainWorkerMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The offchain worker api."},
		},
		RuntimeApiMetadata{
			Name:    grandpaModuleApi,
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
			Name:    accountNonceModuleApi,
			Methods: accountNonceMethodsMd,
			Docs:    sc.Sequence[sc.Str]{" The API to query account nonce."},
		},
		RuntimeApiMetadata{
			Name:    transactionPaymentModuleApi,
			Methods: transactionPaymentMethodsMd,
			Docs:    sc.Sequence[sc.Str]{},
		},
		RuntimeApiMetadata{
			Name:    transactionPaymentCallModuleApi,
			Methods: transactionPaymentCallMethodsMd,
			Docs:    sc.Sequence[sc.Str]{},
		},
		RuntimeApiMetadata{
			Name:    sessionKeysModuleApi,
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
					Type: sc.ToCompact(mdconstants.Header),
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
			Output: sc.ToCompact(mdconstants.TypeOptionOpaqueMetadata),
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
	return sc.Sequence[RuntimeApiMethodMetadata]{
		RuntimeApiMethodMetadata{
			Name:   applyExtrinsicMethod,
			Inputs: applyExtrinsicInputsMd(),
			Output: sc.ToCompact(mdconstants.TypesResult),
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
			Output: sc.ToCompact(mdconstants.TypesSequenceUncheckedExtrinsics),
			Docs:   sc.Sequence[sc.Str]{" Generate inherent extrinsics. The inherent data will vary from chain to chain."},
		},
		RuntimeApiMethodMetadata{
			Name:   checkInherentsMethod,
			Inputs: checkInherentsInputsMd(),
			Output: sc.ToCompact(mdconstants.CheckInherentsResult),
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
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "inherent",
			Type: sc.ToCompact(mdconstants.TypesInherentData),
		},
	}
}

func checkInherentsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "block",
			Type: sc.ToCompact(mdconstants.TypesBlock),
		},
		RuntimeApiMethodParamMetadata{
			Name: "data",
			Type: sc.ToCompact(mdconstants.TypesInherentData),
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
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "call",
			Type: sc.ToCompact(mdconstants.RuntimeCall),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(mdconstants.PrimitiveTypesU32),
		},
	}
}

func queryCallFeeDetailsInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "call",
			Type: sc.ToCompact(mdconstants.RuntimeCall),
		},
		RuntimeApiMethodParamMetadata{
			Name: "len",
			Type: sc.ToCompact(mdconstants.PrimitiveTypesU32),
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
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "seed",
			Type: sc.ToCompact(mdconstants.TypesOptionSequenceU8),
		},
	}
}

func decodeSessionKeysInputsMd() sc.Sequence[RuntimeApiMethodParamMetadata] {
	return sc.Sequence[RuntimeApiMethodParamMetadata]{
		RuntimeApiMethodParamMetadata{
			Name: "encoded",
			Type: sc.ToCompact(mdconstants.TypesSequenceU8),
		},
	}
}
