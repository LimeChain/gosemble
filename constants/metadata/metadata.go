package metadata

const (
	PrimitiveTypesBool = iota
	PrimitiveTypesChar
	PrimitiveTypesString
	PrimitiveTypesU8
	PrimitiveTypesU16
	PrimitiveTypesU32
	PrimitiveTypesU64
	PrimitiveTypesU128
	PrimitiveTypesU256
	PrimitiveTypesI8
	PrimitiveTypesI16
	PrimitiveTypesI32
	PrimitiveTypesI64
	PrimitiveTypesI128
	PrimitiveTypesI256

	TypesFixedSequence4U8
	TypesFixedSequence20U8
	TypesFixedSequence32U8
	TypesFixedSequence64U8
	TypesFixedSequence65U8

	TypesFixedU128

	TypesSequenceU8

	TypesCompactU32
	TypesCompactU64
	TypesCompactU128

	TypesH256
	TypesVecBlockNumEventIndex
	TypesTupleU32U32
	TypesSystemEvent
	TypesSystemEventStorage
	TypesEventRecord
	TypesPhase
	TypesRuntimeEvent
	TypesDispatchInfo
	TypesDispatchClass
	TypesPays
	TypesDispatchError
	TypesModuleError
	TypesTokenError
	TypesArithmeticError
	TypesTransactionalError
	TypesBalancesEvent
	TypesBalanceStatus
	TypesVecTopics
	TypesLastRuntimeUpgradeInfo
	TypesSystemErrors
	TypesBlockLength
	TypesPerDispatchClassU32
	TypesDbWeight
	TypesRuntimeVersion
	TypesRuntimeApis
	TypesRuntimeVecApis
	TypesTupleApiIdU32
	TypesApiId
	TypesAuraStorageAuthorites
	TypesSequencePubKeys
	TypesAuthorityId
	TypesSr25519PubKey
	TypesAuraSlot
	TypesBalancesErrors
	TypesTransactionPaymentReleases
	TypesTransactionPaymentEvents

	TypesEmptyTuple

	TypesAddress32
	TypesMultiAddress

	TypesAccountData
	TypesAccountInfo

	TypesWeight
	TypesOptionWeight
	TypesPerDispatchClassWeight
	TypesPerDispatchClassWeightsPerClass
	TypesWeightPerClass
	TypesBlockWeights

	TypesDigestItem
	TypesSliceDigestItem

	SignatureEd25519
	SignatureSr25519
	SignatureEcdsa
	MultiSignature

	Call

	SystemCalls
	TimestampCalls
	BalancesCalls

	UncheckedExtrinsic
	SignedExtra

	CheckNonZeroSender
	CheckSpecVersion
	CheckTxVersion
	CheckGenesis
	CheckMortality
	CheckNonce
	CheckWeight
	CheckTransactionPayment

	Runtime
)
