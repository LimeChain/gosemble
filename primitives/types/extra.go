package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Extra data, E, is a tuple containing additional meta data about the extrinsic and the system it is meant to be executed in.
// E := (Tmor, N, Pt)
type Extra struct {
	// Tmor: contains the SCALE encoded mortality of the extrinsic
	// Mortality sc.Sequence[sc.U8]
	Era ExtrinsicEra

	// N: a compact integer containing the nonce of the sender.
	// The nonce must be incremented by one for each extrinsic created,
	// otherwise the Polkadot network will reject the extrinsic.
	Nonce sc.Compact // sc.U64

	// Pt: a compact integer containing the transactor pay including tip.
	Fee sc.Compact // sc.U64
}

// type SignedExtra struct {
// 	NonZeroSender MultiAddress
// 	SpecVersion   sc.U32
// 	TxVersion     sc.U32
// 	Genesis       H256
// 	Era           ExtrinsicEra
// 	Nonce         sc.Compact
// 	// Weight
// 	TransactionPayment sc.Compact
// }

func (e Extra) Encode(buffer *bytes.Buffer) {
	// e.ExtrinsicEra.Encode(buffer)
	e.Nonce.Encode(buffer)
	e.Fee.Encode(buffer)
}

func DecodeExtra(buffer *bytes.Buffer) Extra {
	e := Extra{}
	// e.ExtrinsicEra = DecodeExtrinsicEra(buffer)
	e.Nonce = sc.DecodeCompact(buffer)
	e.Fee = sc.DecodeCompact(buffer)
	return e
}

func (e Extra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e Extra) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	return AdditionalSigned{
		SpecVersion:   sc.U32(RuntimeVersion{}.SpecVersion),
		FormatVersion: ExtrinsicFormatVersion,
		// GenesisHash:   H256(),
		// BlockHash: H256(),
		// TransactionVersion sc.U32
		// BlockNumber
	}, nil
}
