package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
)

type CheckSpecVersion struct{}

func (_ CheckSpecVersion) AdditionalSigned() (ok sc.U32, err types.TransactionValidityError) {
	return constants.RuntimeVersion.SpecVersion, err
}

func (_ CheckSpecVersion) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	ok = types.DefaultValidTransaction()
	return ok, err
}

func (v CheckSpecVersion) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = v.Validate(who, call, info, length)
	return ok, err
}
