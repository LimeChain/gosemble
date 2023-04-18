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

	TypesSequenceU8

	TypesCompactU32
	TypesCompactU64
	TypesCompactU128

	TypesEmptyTuple

	TypesAddress32
	TypesMultiAddress

	TypesAccountData
	TypesAccountInfo

	TypesWeight
	TypesPerDispatchClassWeight

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
