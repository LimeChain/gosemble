package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckMortality[N sc.Numeric] struct {
	era          primitives.Era
	systemModule system.Module[N]
}

func NewCheckMortality[N sc.Numeric](systemModule system.Module[N]) CheckMortality[N] {
	return CheckMortality[N]{
		systemModule: systemModule,
	}
}

func (cm CheckMortality[N]) Encode(buffer *bytes.Buffer) {
	cm.era.Encode(buffer)
}

func (cm *CheckMortality[N]) Decode(buffer *bytes.Buffer) {
	cm.era = primitives.DecodeEra(buffer)
}

func (cm CheckMortality[N]) Bytes() []byte {
	return sc.EncodedBytes(cm)
}

func (cm CheckMortality[N]) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	current := sc.U64(cm.systemModule.Storage.BlockNumber.Get()) // TODO: impl saturated_into::<u64>()
	n := sc.U32(cm.era.Birth(current))                           // TODO: impl saturated_into::<T::BlockNumber>()

	if !cm.systemModule.Storage.BlockHash.Exists(N(n)) {
		return nil, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
	}

	return sc.NewVaryingData(primitives.H256(cm.systemModule.Storage.BlockHash.Get(N(n)))), nil
}

func (cm CheckMortality[N]) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	currentU64 := sc.U64(cm.systemModule.Storage.BlockNumber.Get()) // TODO: per module implementation

	validTill := cm.era.Death(currentU64)

	ok := primitives.DefaultValidTransaction()
	ok.Longevity = validTill.SaturatingSub(currentU64)

	return ok, nil
}

func (cm CheckMortality[N]) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cm CheckMortality[N]) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := cm.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cm CheckMortality[N]) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cm.ValidateUnsigned(call, info, length)
	return err
}

func (cm CheckMortality[N]) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
