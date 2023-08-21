package extensions

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce[N sc.Numeric] struct {
	nonce        sc.U32
	systemModule system.Module[N]
}

func NewCheckNonce[N sc.Numeric](systemModule system.Module[N]) CheckNonce[N] {
	return CheckNonce[N]{systemModule: systemModule}
}

func (cn CheckNonce[N]) Encode(buffer *bytes.Buffer) {
	sc.ToCompact(cn.nonce).Encode(buffer)
}

func (cn *CheckNonce[N]) Decode(buffer *bytes.Buffer) {
	cn.nonce = sc.U32(sc.U128(sc.DecodeCompact(buffer)).ToBigInt().Uint64())
}

func (cn CheckNonce[N]) Bytes() []byte {
	return sc.EncodedBytes(cn)
}

func (cn CheckNonce[N]) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(), nil
}

func (cn CheckNonce[N]) Validate(who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	account := cn.systemModule.Storage.Account.Get((*who).FixedSequence)

	if cn.nonce < account.Nonce {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
	}

	encoded := (*who).Bytes()
	encoded = append(encoded, sc.ToCompact(cn.nonce).Bytes()...)
	provides := sc.Sequence[primitives.TransactionTag]{sc.BytesToSequenceU8(encoded)}

	var requires sc.Sequence[primitives.TransactionTag]
	if account.Nonce < cn.nonce {
		encoded := (*who).Bytes()
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

func (cn CheckNonce[N]) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cn CheckNonce[N]) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	account := cn.systemModule.Storage.Account.Get(who.FixedSequence)

	if cn.nonce != account.Nonce {
		var err primitives.TransactionValidityError
		if cn.nonce < account.Nonce {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
		} else {
			err = primitives.NewTransactionValidityError(primitives.NewInvalidTransactionFuture())
		}
		return primitives.Pre{}, err
	}

	account.Nonce += 1
	cn.systemModule.Storage.Account.Put(who.FixedSequence, account)

	return primitives.Pre{}, nil
}

func (cn CheckNonce[N]) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cn.ValidateUnsigned(call, info, length)
	return err
}

func (cn CheckNonce[N]) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}

func (cn CheckNonce[N]) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckNonce,
			"CheckNonce",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_nonce", "CheckNonce"},
			primitives.NewMetadataTypeDefinitionCompact(sc.ToCompact(metadata.PrimitiveTypesU32)),
		),
		primitives.NewMetadataSignedExtension("CheckNonce", metadata.CheckNonce, metadata.TypesEmptyTuple)
}
