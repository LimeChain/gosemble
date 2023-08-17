package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckTxVersion[N sc.Numeric] struct {
	systemModule system.Module[N]
}

func NewCheckTxVersion[N sc.Numeric](module system.Module[N]) CheckTxVersion[N] {
	return CheckTxVersion[N]{
		systemModule: module,
	}
}

func (ctv CheckTxVersion[N]) Encode(*bytes.Buffer) {}

func (ctv CheckTxVersion[N]) Decode(*bytes.Buffer) {}

func (ctv CheckTxVersion[N]) Bytes() []byte {
	return sc.EncodedBytes(ctv)
}

func (ctv CheckTxVersion[N]) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(ctv.systemModule.Constants.Version.TransactionVersion), nil
}

func (_ CheckTxVersion[N]) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion[N]) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion[N]) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := ctv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (ctv CheckTxVersion[N]) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := ctv.ValidateUnsigned(call, info, length)
	return err
}

func (ctv CheckTxVersion[N]) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
