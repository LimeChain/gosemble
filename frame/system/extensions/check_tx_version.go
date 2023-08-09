package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system/module"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckTxVersion struct {
	systemModule module.SystemModule
}

func NewCheckTxVersion(module module.SystemModule) CheckTxVersion {
	return CheckTxVersion{
		systemModule: module,
	}
}

func (ctv CheckTxVersion) Encode(*bytes.Buffer) {}

func (ctv CheckTxVersion) Decode(*bytes.Buffer) {}

func (ctv CheckTxVersion) Bytes() []byte {
	return sc.EncodedBytes(ctv)
}

func (ctv CheckTxVersion) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(ctv.systemModule.Constants.Version.TransactionVersion), nil
}

func (_ CheckTxVersion) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := ctv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (ctv CheckTxVersion) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := ctv.ValidateUnsigned(call, info, length)
	return err
}

func (ctv CheckTxVersion) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
