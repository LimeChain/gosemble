package extensions

import (
	"bytes"
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce struct {
	nonce                         sc.U32
	systemModule                  system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckNonce(systemModule system.Module) primitives.SignedExtension {
	return &CheckNonce{
		systemModule:                  systemModule,
		typesInfoAdditionalSignedData: sc.NewVaryingData(),
	}
}

func (cn CheckNonce) Encode(buffer *bytes.Buffer) error {
	return sc.ToCompact(cn.nonce).Encode(buffer)
}

func (cn *CheckNonce) Decode(buffer *bytes.Buffer) error {
	compactNonce, err := sc.DecodeCompact[sc.U128](buffer)
	if err != nil {
		return err
	}
	cn.nonce = sc.U32(compactNonce.ToBigInt().Uint64())
	return nil
}

func (cn CheckNonce) Bytes() []byte {
	return sc.EncodedBytes(cn)
}

func (cn CheckNonce) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return sc.NewVaryingData(), nil
}

func (cn CheckNonce) Validate(who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, error) {
	account, err := cn.systemModule.StorageAccount(who)
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	if cn.nonce < account.Nonce {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionStale())
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

func (cn CheckNonce) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (cn CheckNonce) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, error) {
	account, err := cn.systemModule.StorageAccount(who)
	if err != nil {
		return primitives.Pre{}, err
	}

	if cn.nonce != account.Nonce {
		var txErr sc.Encodable
		if cn.nonce < account.Nonce {
			txErr = primitives.NewInvalidTransactionStale()
		} else {
			txErr = primitives.NewInvalidTransactionFuture()
		}
		return primitives.Pre{}, primitives.NewTransactionValidityError(txErr)
	}

	account.Nonce = account.Nonce + 1
	cn.systemModule.StorageAccountSet(who, account)

	return primitives.Pre{}, nil
}

func (cn CheckNonce) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) error {
	_, err := cn.ValidateUnsigned(call, info, length)
	return err
}

func (cn CheckNonce) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) error {
	return nil
}

func (cn CheckNonce) ModulePath() string {
	return systemModulePath
}
