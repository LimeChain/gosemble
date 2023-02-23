package system

import (
	"math"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonce sc.U64

func (n CheckNonce) Validate(who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _lenght sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	// TODO: check if we can use just who
	accountNonce := system.StorageAccountNonce((*who).FixedSequence) // account = Account::<T>::get(who)

	if sc.U64(n) < sc.U64(accountNonce) {
		err = types.NewTransactionValidityError(types.NewInvalidTransaction(types.StaleError))
		return ok, err
	}

	encoded := (*who).Bytes()
	encoded = append(encoded, sc.ToCompact(sc.U64(n)).Bytes()...) // TODO: confirm if it is compact encoded
	provides := sc.Sequence[types.TransactionTag]{sc.BytesToSequenceU8(encoded)}

	var requires sc.Sequence[types.TransactionTag]
	if sc.U64(accountNonce) < sc.U64(n) {
		encoded := (*who).Bytes()
		encoded = append(encoded, sc.ToCompact(sc.U64(n)-1).Bytes()...)
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
