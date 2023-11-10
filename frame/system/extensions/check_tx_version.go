package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckTxVersion struct {
	systemModule system.Module
}

func NewCheckTxVersion(module system.Module) CheckTxVersion {
	return CheckTxVersion{
		systemModule: module,
	}
}

func (ctv CheckTxVersion) Encode(*bytes.Buffer) error {
	return nil
}

func (ctv CheckTxVersion) Decode(*bytes.Buffer) error { return nil }

func (ctv CheckTxVersion) Bytes() []byte {
	return sc.EncodedBytes(ctv)
}

func (ctv CheckTxVersion) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(ctv.systemModule.Version().TransactionVersion), nil
}

func (_ CheckTxVersion) Validate(_who primitives.AccountId[primitives.PublicKey], _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) PreDispatch(who primitives.AccountId[primitives.PublicKey], call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := ctv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (ctv CheckTxVersion) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := ctv.ValidateUnsigned(call, info, length)
	return err
}

func (ctv CheckTxVersion) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}

func (ctv CheckTxVersion) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckTxVersion,
			"CheckTxVersion",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_tx_version", "CheckTxVersion"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
		),
		primitives.NewMetadataSignedExtension("CheckTxVersion", metadata.CheckTxVersion, metadata.PrimitiveTypesU32)
}
