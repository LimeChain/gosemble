package system

import (
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce sc.U32

func (n CheckNonce) AdditionalSigned() (ok sc.Empty, err types.TransactionValidityError) {
	ok = sc.Empty{}
	return ok, err
}

func (n CheckNonce) Validate(who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _lenght sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	// TODO: check if we can use just who
	account := system.StorageGetAccount((*who).FixedSequence)

	if sc.U32(n) < account.Nonce {
		err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.StaleError))
		return ok, err
	}

	encoded := (*who).Bytes()
	encoded = append(encoded, sc.ToCompact(sc.U32(n)).Bytes()...)
	provides := sc.Sequence[types.TransactionTag]{sc.BytesToSequenceU8(encoded)}

	var requires sc.Sequence[types.TransactionTag]
	if account.Nonce < sc.U32(n) {
		encoded := (*who).Bytes()
		encoded = append(encoded, sc.ToCompact(sc.U32(n)-1).Bytes()...)
		requires = sc.Sequence[types.TransactionTag]{sc.BytesToSequenceU8(encoded)}
	} else {
		requires = sc.Sequence[types.TransactionTag]{}
	}

	ok = types.ValidTransaction{
		Priority:  0,
		Requires:  requires,
		Provides:  provides,
		Longevity: types.TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}

	return ok, err
}

func (n CheckNonce) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	account := system.StorageGetAccount(who.FixedSequence)

	if sc.U32(n) != account.Nonce {
		if sc.U32(n) < account.Nonce {
			err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.StaleError))
		} else {
			err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.FutureError))
		}
		return ok, err
	}

	account.Nonce += 1
	system.StorageSetAccount(who.FixedSequence, account)

	ok = types.Pre{}
	return ok, err
}
