package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// The source of the transaction.
//
// Depending on the source we might apply different validation schemes.
// For instance, we can disallow specific kinds of transactions if they were not produced
// by our local node (for instance off-chain workers).
const (
	// Transaction is already included in block.
	//
	// This means that we can't really tell where the transaction is coming from,
	// since it's already in the received block. Note that the custom validation logic
	// using either `Local` or `External` should most likely just allow `InBlock`
	// transactions as well.
	TransactionSourceInBlock sc.U8 = iota

	// Transaction is coming from a local source.
	//
	// This means that the transaction was produced internally by the node
	// (for instance an Off-Chain Worker, or an Off-Chain Call), as opposed
	// to being received over the network.
	TransactionSourceLocal

	// Transaction has been received externally.
	//
	// This means the transaction has been received from (usually) "untrusted" source,
	// for instance received over the network or RPC.
	TransactionSourceExternal
)

type TransactionSource sc.VaryingData

func NewTransactionSourceInBlock() TransactionSource {
	return TransactionSource(sc.NewVaryingData(TransactionSourceInBlock))
}

func NewTransactionSourceLocal() TransactionSource {
	return TransactionSource(sc.NewVaryingData(TransactionSourceLocal))
}

func NewTransactionSourceExternal() TransactionSource {
	return TransactionSource(sc.NewVaryingData(TransactionSourceExternal))
}

func (ts TransactionSource) Encode(buffer *bytes.Buffer) error {
	if len(ts) == 0 {
		return newTypeError("TransactionSource")
	}

	return ts[0].Encode(buffer)
}

func DecodeTransactionSource(buffer *bytes.Buffer) (TransactionSource, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return TransactionSource{}, err
	}

	switch b {
	case TransactionSourceInBlock:
		return NewTransactionSourceInBlock(), nil
	case TransactionSourceLocal:
		return NewTransactionSourceLocal(), nil
	case TransactionSourceExternal:
		return NewTransactionSourceExternal(), nil
	default:
		return TransactionSource{}, newTypeError("TransactionSource")
	}
}

func (ts TransactionSource) Bytes() []byte {
	return sc.EncodedBytes(ts)
}

func (ts TransactionSource) MetadataDefinition() *MetadataTypeDefinition {
	def := NewMetadataTypeDefinitionVariant(
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
		})
	return &def
}
