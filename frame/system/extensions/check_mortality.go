package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckMortality struct {
	era                           primitives.Era
	systemModule                  system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckMortality(systemModule system.Module) primitives.SignedExtension {
	return &CheckMortality{systemModule: systemModule,
		typesInfoAdditionalSignedData: sc.NewVaryingData(primitives.H256{})}
}

func (cm CheckMortality) Encode(buffer *bytes.Buffer) error {
	return cm.era.Encode(buffer)
}

func (cm *CheckMortality) Decode(buffer *bytes.Buffer) error {
	era, err := primitives.DecodeEra(buffer)
	if err != nil {
		return err
	}
	cm.era = era
	return nil
}

func (cm CheckMortality) Bytes() []byte {
	return sc.EncodedBytes(cm)
}

func (cm CheckMortality) AdditionalSigned() (primitives.AdditionalSigned, error) {
	current, err := cm.systemModule.StorageBlockNumber()
	if err != nil {
		return nil, err
	}
	n := cm.era.Birth(current)

	if !cm.systemModule.StorageBlockHashExists(n) {
		return nil, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
	}

	blockHash, err := cm.systemModule.StorageBlockHash(n)
	if err != nil {
		return nil, err
	}
	hash, err := primitives.NewH256(blockHash.FixedSequence...)
	if err != nil {
		return nil, err
	}
	return sc.NewVaryingData(hash), nil
}

func (cm CheckMortality) Validate(_who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	currentBlockNum, err := cm.systemModule.StorageBlockNumber()
	if err != nil {
		return primitives.ValidTransaction{}, err
	}

	validTill := cm.era.Death(currentBlockNum)

	ok := primitives.DefaultValidTransaction()
	ok.Longevity = sc.SaturatingSubU64(validTill, currentBlockNum)

	return ok, nil
}

func (cm CheckMortality) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (cm CheckMortality) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.Pre, error) {
	_, err := cm.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cm CheckMortality) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) error {
	_, err := cm.ValidateUnsigned(call, info, length)
	return err
}

func (cm CheckMortality) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact[sc.Numeric], _result *primitives.DispatchResult) error {
	return nil
}

func (cm CheckMortality) ModulePath() string {
	return systemModulePath
}
