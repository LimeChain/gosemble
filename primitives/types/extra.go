package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// implements SignedExtension
type SignedExtra struct {
	// Extra data, E, is a tuple containing additional meta data about the extrinsic and the system it is meant to be executed in.
	// E := (Tmor, N, Pt)

	// NonZeroSender MultiAddress

	Era Era

	// N: a compact integer containing the nonce of the sender.
	// The nonce must be incremented by one for each extrinsic created,
	// otherwise the Polkadot network will reject the extrinsic.
	Nonce sc.U64 // encode as Compact

	// Pt: a compact integer containing the transactor pay including tip.
	Fee sc.U64 // encode as Compact
	// 	TransactionPayment sc.Compact

	Weight Weight

	// SpecVersion sc.U32
	// TxVersion   sc.U32
	// Genesis     H256
}

func (e SignedExtra) Encode(buffer *bytes.Buffer) {
	e.Era.Encode(buffer)
	sc.ToCompact(e.Nonce).Encode(buffer)
	sc.ToCompact(e.Fee).Encode(buffer)
}

func DecodeExtra(buffer *bytes.Buffer) SignedExtra {
	e := SignedExtra{}
	e.Era = DecodeEra(buffer)
	e.Nonce = sc.U64(sc.U128(sc.DecodeCompact(buffer)).ToBigInt().Uint64())
	e.Fee = sc.U64(sc.U128(sc.DecodeCompact(buffer)).ToBigInt().Uint64())
	return e
}

func (e SignedExtra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e SignedExtra) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	return AdditionalSigned{
		SpecVersion:   sc.U32(RuntimeVersion{}.SpecVersion),
		FormatVersion: ExtrinsicFormatVersion,
		// GenesisHash:   H256(),
		// BlockHash: H256(),
		// TransactionVersion sc.U32
		// BlockNumber
	}, nil
}
