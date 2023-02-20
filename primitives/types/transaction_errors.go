package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// Errors that can occur while checking the validity of a transaction.
type TransactionValidityError sc.VaryingData

func NewTransactionValidityError(value sc.Encodable) TransactionValidityError {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction:
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	return TransactionValidityError(sc.NewVaryingData(value))
}

func (e TransactionValidityError) Encode(buffer *bytes.Buffer) {
	value := e[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(InvalidTransaction)):
		buffer.Write([]byte{0x00})
	case reflect.TypeOf(*new(UnknownTransaction)):
		buffer.Write([]byte{0x01})
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	value.Encode(buffer)
}

func DecodeTransactionValidityError(buffer *bytes.Buffer) TransactionValidityError {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		value := DecodeInvalidTransaction(buffer)
		return NewTransactionValidityError(value)
	case 1:
		value := DecodeUnknownTransaction(buffer)
		return NewTransactionValidityError(value)
	default:
		log.Critical("invalid TransactionValidityError type")
	}

	panic("unreachable")
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

const (
	// The call of the transaction is not expected. Reject
	CallError sc.U8 = iota

	// General error to do with the inability to pay some fees (e.g. account balance too low). Reject
	PaymentError

	// General error to do with the transaction not yet being valid (e.g. nonce too high). Don't Reject
	FutureError

	// General error to do with the transaction being outdated (e.g. nonce too low). Reject
	StaleError

	// General error to do with the transaction's proofs (e.g. signature). Reject
	//
	// # Possible causes
	//
	// When using a signed extension that provides additional data for signing, it is required
	// that the signing and the verifying side use the same additional data. Additional
	// data will only be used to generate the signature, but will not be part of the transaction
	// itself. As the verifying side does not know which additional data was used while signing
	// it will only be able to assume a bad signature and cannot express a more meaningful error.
	BadProofError

	// The transaction birth block is ancient. Reject
	//
	// # Possible causes
	//
	// For `FRAME`-based runtimes this would be caused by `current block number
	// - Era::birth block number > BlockHashCount`. (e.g. in Polkadot `BlockHashCount` = 2400, so
	//   a
	// transaction with birth block number 1337 would be valid up until block number 1337 + 2400,
	// after which point the transaction would be considered to have an ancient birth block.)
	AncientBirthBlockError

	// The transaction would exhaust the resources of the current block. Don't Reject
	//
	// The transaction might be valid, but there are not enough resources
	// left in the current block.
	ExhaustsResourcesError

	// Any other custom invalid validity that is not covered by this enum. Reject
	CustomInvalidTransactionError // + sc.U8

	// An extrinsic with mandatory dispatch resulted in an error. Reject
	// This is indicative of either a malicious validator or a buggy `provide_inherent`.
	// In any case, it can result in dangerously overweight blocks and therefore if
	// found, invalidates the block.
	BadMandatoryError

	// An extrinsic with a mandatory dispatch tried to be validated.
	// This is invalid; only inherent extrinsics are allowed to have mandatory dispatches.
	MandatoryValidationError

	// The sending address is disabled or known to be invalid.
	BadSignerError
)

type InvalidTransaction sc.VaryingData

func NewInvalidTransaction(values ...sc.Encodable) InvalidTransaction {
	switch values[0] {
	case CallError, PaymentError, FutureError, StaleError, BadProofError, AncientBirthBlockError, ExhaustsResourcesError, BadMandatoryError, MandatoryValidationError, BadSignerError:
		return InvalidTransaction(sc.NewVaryingData(values[0]))
	case CustomInvalidTransactionError:
		return InvalidTransaction(sc.NewVaryingData(values[0:2]...))
	default:
		log.Critical("invalid InvalidTransaction type")
	}

	panic("unreachable")
}

func (e InvalidTransaction) Encode(buffer *bytes.Buffer) {
	switch e[0] {
	case CallError:
		sc.U8(0).Encode(buffer)
	case PaymentError:
		sc.U8(1).Encode(buffer)
	case FutureError:
		sc.U8(2).Encode(buffer)
	case StaleError:
		sc.U8(3).Encode(buffer)
	case BadProofError:
		sc.U8(4).Encode(buffer)
	case AncientBirthBlockError:
		sc.U8(5).Encode(buffer)
	case ExhaustsResourcesError:
		sc.U8(6).Encode(buffer)
	case CustomInvalidTransactionError:
		sc.U8(7).Encode(buffer)
		e[1].Encode(buffer)
	case BadMandatoryError:
		sc.U8(8).Encode(buffer)
	case MandatoryValidationError:
		sc.U8(9).Encode(buffer)
	case BadSignerError:
		sc.U8(10).Encode(buffer)
	default:
		log.Critical("invalid InvalidTransaction type")
	}
}

func DecodeInvalidTransaction(buffer *bytes.Buffer) InvalidTransaction {
	b := sc.DecodeU8(buffer)

	switch b {
	case sc.U8(0):
		return NewInvalidTransaction(CallError)
	case sc.U8(1):
		return NewInvalidTransaction(PaymentError)
	case sc.U8(2):
		return NewInvalidTransaction(FutureError)
	case sc.U8(3):
		return NewInvalidTransaction(StaleError)
	case sc.U8(4):
		return NewInvalidTransaction(BadProofError)
	case sc.U8(5):
		return NewInvalidTransaction(AncientBirthBlockError)
	case sc.U8(6):
		return NewInvalidTransaction(ExhaustsResourcesError)
	case sc.U8(7):
		v := sc.DecodeU8(buffer)
		return NewInvalidTransaction(CustomInvalidTransactionError, v)
	case sc.U8(8):
		return NewInvalidTransaction(BadMandatoryError)
	case sc.U8(9):
		return NewInvalidTransaction(MandatoryValidationError)
	case sc.U8(10):
		return NewInvalidTransaction(BadSignerError)
	default:
		log.Critical("invalid InvalidTransaction type")
	}

	panic("unreachable")
}

func (e InvalidTransaction) Bytes() []byte {
	return sc.EncodedBytes(e)
}

const (
	// Could not lookup some information that is required to validate the transaction. Reject
	CannotLookupError sc.U8 = iota

	// No validator found for the given unsigned transaction. Reject
	NoUnsignedValidatorError

	// Any other custom unknown validity that is not covered by this type. Reject
	CustomUnknownTransactionError // + sc.U8
)

type UnknownTransaction sc.VaryingData

func NewUnknownTransaction(values ...sc.Encodable) UnknownTransaction {
	switch values[0] {
	case CannotLookupError, NoUnsignedValidatorError:
		return UnknownTransaction(sc.NewVaryingData(values[0]))
	case CustomUnknownTransactionError:
		return UnknownTransaction(sc.NewVaryingData(values[0:2]...))
	default:
		log.Critical("invalid UnknownTransaction type")
	}

	panic("unreachable")
}

func (e UnknownTransaction) Encode(buffer *bytes.Buffer) {
	switch e[0] {
	case CannotLookupError:
		sc.U8(0).Encode(buffer)
	case NoUnsignedValidatorError:
		sc.U8(1).Encode(buffer)
	case CustomUnknownTransactionError:
		sc.U8(2).Encode(buffer)
	default:
		log.Critical("invalid UnknownTransaction type")
	}
}

func DecodeUnknownTransaction(buffer *bytes.Buffer) UnknownTransaction {
	b := sc.DecodeU8(buffer)

	switch b {
	case sc.U8(0):
		return NewUnknownTransaction(CannotLookupError)
	case sc.U8(1):
		return NewUnknownTransaction(NoUnsignedValidatorError)
	case sc.U8(2):
		v := sc.DecodeU8(buffer)
		return NewUnknownTransaction(CustomUnknownTransactionError, v)
	default:
		log.Critical("invalid UnknownTransaction type")
	}

	panic("unreachable")
}

func (e UnknownTransaction) Bytes() []byte {
	return sc.EncodedBytes(e)
}
