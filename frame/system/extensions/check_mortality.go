package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
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

func (cm CheckMortality) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	current, err := cm.systemModule.StorageBlockNumber() // TODO: impl saturated_into::<u64>()
	if err != nil {
		// TODO https://github.com/LimeChain/gosemble/issues/271
		transactionValidityError, _ := primitives.NewTransactionValidityError(sc.Str(err.Error()))
		return nil, transactionValidityError
	}
	n := cm.era.Birth(current) // TODO: impl saturated_into::<T::BlockNumber>()

	if !cm.systemModule.StorageBlockHashExists(n) {
		invalidTransactionAncientBirthBlock, _ := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionAncientBirthBlock())
		return nil, invalidTransactionAncientBirthBlock
	}

	blockHash, err := cm.systemModule.StorageBlockHash(n)
	if err != nil {
		transactionValidityError, _ := primitives.NewTransactionValidityError(sc.Str(err.Error()))
		return nil, transactionValidityError
	}
	hash, err := primitives.NewH256(blockHash.FixedSequence...)
	if err != nil {
		transactionValidityError, _ := primitives.NewTransactionValidityError(sc.Str(err.Error()))
		return nil, transactionValidityError
	}
	return sc.NewVaryingData(hash), nil
}

func (cm CheckMortality) Validate(_who primitives.Address32, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	currentBlockNum, err := cm.systemModule.StorageBlockNumber() // TODO: per module implementation
	if err != nil {
		log.Critical(err.Error())
	}

	validTill := cm.era.Death(currentBlockNum)

	ok := primitives.DefaultValidTransaction()
	ok.Longevity = sc.SaturatingSubU64(validTill, currentBlockNum)

	return ok, nil
}

func (cm CheckMortality) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (cm CheckMortality) PreDispatch(who primitives.Address32, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := cm.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cm CheckMortality) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := cm.ValidateUnsigned(call, info, length)
	return err
}

func (cm CheckMortality) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}

func (cm CheckMortality) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckMortality,
			"CheckMortality",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_mortality", "CheckMortality"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesEra, "Era"),
				},
			),
		),
		primitives.NewMetadataSignedExtension("CheckMortality", metadata.CheckMortality, metadata.TypesH256)
}
