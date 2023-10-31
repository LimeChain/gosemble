package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	// The call of the transaction is not expected. Reject
	InvalidTransactionCall sc.U8 = iota

	// General error to do with the inability to pay some fees (e.g. account balance too low). Reject
	InvalidTransactionPayment

	// General error to do with the transaction not yet being valid (e.g. nonce too high). Don't Reject
	InvalidTransactionFuture

	// General error to do with the transaction being outdated (e.g. nonce too low). Reject
	InvalidTransactionStale

	// General error to do with the transaction's proofs (e.g. signature). Reject
	//
	// # Possible causes
	//
	// When using a signed extension that provides additional data for signing, it is required
	// that the signing and the verifying side use the same additional data. Additional
	// data will only be used to generate the signature, but will not be part of the transaction
	// itself. As the verifying side does not know which additional data was used while signing
	// it will only be able to assume a bad signature and cannot express a more meaningful error.
	InvalidTransactionBadProof

	// The transaction birth block is ancient. Reject
	//
	// # Possible causes
	//
	// For `FRAME`-based runtimes this would be caused by `current block number
	// - Era::birth block number > BlockHashCount`. (e.g. in Polkadot `BlockHashCount` = 2400, so
	//   a
	// transaction with birth block number 1337 would be valid up until block number 1337 + 2400,
	// after which point the transaction would be considered to have an ancient birth block.)
	InvalidTransactionAncientBirthBlock

	// The transaction would exhaust the resources of the current block. Don't Reject
	//
	// The transaction might be valid, but there are not enough resources
	// left in the current block.
	InvalidTransactionExhaustsResources

	// Any other custom invalid validity that is not covered by this enum. Reject
	InvalidTransactionCustom // + sc.U8

	// An extrinsic with mandatory dispatch resulted in an error. Reject
	// This is indicative of either a malicious validator or a buggy `provide_inherent`.
	// In any case, it can result in dangerously overweight blocks and therefore if
	// found, invalidates the block.
	InvalidTransactionBadMandatory

	// An extrinsic with a mandatory dispatch tried to be validated.
	// This is invalid; only inherent extrinsics are allowed to have mandatory dispatches.
	InvalidTransactionMandatoryValidation

	// The sending address is disabled or known to be invalid.
	InvalidTransactionBadSigner
)

type InvalidTransaction struct {
	sc.VaryingData
}

func NewInvalidTransactionCall() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionCall)}
}

func NewInvalidTransactionPayment() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionPayment)}
}

func NewInvalidTransactionFuture() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionFuture)}
}

func NewInvalidTransactionStale() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionStale)}
}

func NewInvalidTransactionBadProof() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadProof)}
}

func NewInvalidTransactionAncientBirthBlock() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionAncientBirthBlock)}
}

func NewInvalidTransactionExhaustsResources() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionExhaustsResources)}
}

func NewInvalidTransactionCustom(customError sc.U8) InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionCustom, customError)}
}

func NewInvalidTransactionBadMandatory() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadMandatory)}
}

func NewInvalidTransactionMandatoryValidation() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionMandatoryValidation)}
}

func NewInvalidTransactionBadSigner() InvalidTransaction {
	return InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadSigner)}
}

func DecodeInvalidTransaction(buffer *bytes.Buffer) (InvalidTransaction, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return InvalidTransaction{}, err
	}

	switch b {
	case InvalidTransactionCall:
		return NewInvalidTransactionCall(), nil
	case InvalidTransactionPayment:
		return NewInvalidTransactionPayment(), nil
	case InvalidTransactionFuture:
		return NewInvalidTransactionFuture(), nil
	case InvalidTransactionStale:
		return NewInvalidTransactionStale(), nil
	case InvalidTransactionBadProof:
		return NewInvalidTransactionBadProof(), nil
	case InvalidTransactionAncientBirthBlock:
		return NewInvalidTransactionAncientBirthBlock(), nil
	case InvalidTransactionExhaustsResources:
		return NewInvalidTransactionExhaustsResources(), nil
	case InvalidTransactionCustom:
		v, err := sc.DecodeU8(buffer)
		if err != nil {
			return InvalidTransaction{}, err
		}
		return NewInvalidTransactionCustom(v), nil
	case InvalidTransactionBadMandatory:
		return NewInvalidTransactionBadMandatory(), nil
	case InvalidTransactionMandatoryValidation:
		return NewInvalidTransactionMandatoryValidation(), nil
	case InvalidTransactionBadSigner:
		return NewInvalidTransactionBadSigner(), nil
	default:
		log.Critical("invalid InvalidTransaction type")
	}

	panic("unreachable")
}
