package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckGenesis[N sc.Numeric] struct {
	module system.Module[N]
}

func NewCheckGenesis[N sc.Numeric](module system.Module[N]) CheckGenesis[N] {
	return CheckGenesis[N]{
		module,
	}
}

func (cg CheckGenesis[N]) Encode(*bytes.Buffer) {}

func (cg CheckGenesis[N]) Decode(*bytes.Buffer) {}

func (cg CheckGenesis[N]) Bytes() []byte {
	return sc.EncodedBytes(cg)
}

func (cg CheckGenesis[N]) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	hash := cg.module.Storage.BlockHash.Get(sc.NewNumeric[N](uint8(0)))

	return sc.NewVaryingData(primitives.H256(hash)), nil
}

func (_ CheckGenesis[N]) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cg CheckGenesis[N]) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cg CheckGenesis[N]) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := cg.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cg CheckGenesis[N]) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cg.ValidateUnsigned(call, info, length)
	return err
}

func (cg CheckGenesis[N]) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
