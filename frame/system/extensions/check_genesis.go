package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckGenesis struct{}

func (_ CheckGenesis) AdditionalSigned() (ok types.H256, err types.TransactionValidityError) {
	ok = types.H256(system.StorageGetBlockHash(sc.U32(0)))
	return ok, err
}

func (_ CheckGenesis) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	ok = types.DefaultValidTransaction()
	return ok, err
}

func (g CheckGenesis) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = g.Validate(who, call, info, length)
	return ok, err
}
