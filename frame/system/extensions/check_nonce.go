package extensions

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce struct {
	nonce        sc.U32
	systemModule system.Module
}

func NewCheckNonce(systemModule system.Module) CheckNonce {
	return CheckNonce{systemModule: systemModule}
}

func (cn CheckNonce) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(cn.nonce).Encode(buffer)
}

func (cn *CheckNonce) Decode(buffer *bytes.Buffer) error {
	compactNonce, err := sc.DecodeCompact(buffer)
	if err != nil {
		return err
	}
	cn.nonce = sc.U32(compactNonce.ToBigInt().Uint64())
	return nil
}

func (cn CheckNonce) Bytes() []byte {
	return sc.EncodedBytes(cn)
}

func (cn CheckNonce) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(), nil
}

func (cn CheckNonce) Validate(who primitives.Address32, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	account, err := cn.systemModule.StorageAccount(who.FixedSequence)
	if err != nil {
		// TODO https://github.com/LimeChain/gosemble/issues/271
		transactionValidityError, _ := primitives.NewTransactionValidityError(sc.Str(err.Error()))
		return primitives.ValidTransaction{}, transactionValidityError
	}

	if cn.nonce < account.Nonce {
		// TODO https://github.com/LimeChain/gosemble/issues/271
		invalidTransactionStale, _ := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
		return primitives.ValidTransaction{}, invalidTransactionStale
	}

	encoded := who.Bytes()
	encoded = append(encoded, sc.ToCompact(cn.nonce).Bytes()...)
	provides := sc.Sequence[primitives.TransactionTag]{sc.BytesToSequenceU8(encoded)}

	var requires sc.Sequence[primitives.TransactionTag]
	if account.Nonce < cn.nonce {
		encoded := who.Bytes()
		encoded = append(encoded, sc.ToCompact(cn.nonce-1).Bytes()...)
		requires = sc.Sequence[primitives.TransactionTag]{sc.BytesToSequenceU8(encoded)}
	} else {
		requires = sc.Sequence[primitives.TransactionTag]{}
	}

	return primitives.ValidTransaction{
		Priority:  0,
		Requires:  requires,
		Provides:  provides,
		Longevity: primitives.TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}, nil
}

func (cn CheckNonce) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cn CheckNonce) PreDispatch(who primitives.Address32, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	account, err := cn.systemModule.StorageAccount(who.FixedSequence)
	if err != nil {
		// TODO https://github.com/LimeChain/gosemble/issues/271
		transactionValidityError, _ := primitives.NewTransactionValidityError(sc.Str(err.Error()))
		return primitives.Pre{}, transactionValidityError
	}

	if cn.nonce != account.Nonce {
		var transactionValidityError primitives.TransactionValidityError
		if cn.nonce < account.Nonce {
			// TODO https://github.com/LimeChain/gosemble/issues/271
			transactionValidityError, _ = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
		} else {
			transactionValidityError, _ = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionFuture())
		}
		return primitives.Pre{}, transactionValidityError
	}

	account.Nonce = account.Nonce + 1
	cn.systemModule.StorageAccountSet(who.FixedSequence, account)

	return primitives.Pre{}, nil
}

func (cn CheckNonce) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cn.ValidateUnsigned(call, info, length)
	return err
}

func (cn CheckNonce) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}

func (cn CheckNonce) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckNonce,
			"CheckNonce",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_nonce", "CheckNonce"},
			primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32)),
		),
		primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple)
}
