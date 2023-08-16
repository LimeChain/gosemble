package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckMortality struct {
	era          primitives.Era
	systemModule system.Module
}

func NewCheckMortality(systemModule system.Module) CheckMortality {
	return CheckMortality{
		systemModule: systemModule,
	}
}

func (cm CheckMortality) Encode(buffer *bytes.Buffer) {
	cm.era.Encode(buffer)
}

func (cm *CheckMortality) Decode(buffer *bytes.Buffer) {
	cm.era = primitives.DecodeEra(buffer)
}

func (cm CheckMortality) Bytes() []byte {
	return sc.EncodedBytes(cm)
}

func (cm CheckMortality) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	current := sc.U64(cm.systemModule.Storage.BlockNumber.Get()) // TODO: impl saturated_into::<u64>()
	n := sc.U32(cm.era.Birth(current))                           // TODO: impl saturated_into::<T::BlockNumber>()

	if !cm.systemModule.Storage.BlockHash.Exists(n) {
		return nil, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
	}

	return sc.NewVaryingData(primitives.H256(cm.systemModule.Storage.BlockHash.Get(n))), nil
}

func (cm CheckMortality) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	currentU64 := sc.U64(cm.systemModule.Storage.BlockNumber.Get()) // TODO: per module implementation

	validTill := cm.era.Death(currentU64)

	ok := primitives.DefaultValidTransaction()
	ok.Longevity = validTill.SaturatingSub(currentU64)

	return ok, nil
}

func (cm CheckMortality) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cm CheckMortality) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := cm.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cm CheckMortality) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cm.ValidateUnsigned(call, info, length)
	return err
}

func (cm CheckMortality) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
