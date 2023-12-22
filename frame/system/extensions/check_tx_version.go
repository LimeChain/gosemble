package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckTxVersion struct {
	systemModule                  system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckTxVersion(module system.Module) primitives.SignedExtension {
	return &CheckTxVersion{
		systemModule:                  module,
		typesInfoAdditionalSignedData: sc.NewVaryingData(sc.U32(0)),
	}
}

func (ctv CheckTxVersion) Encode(*bytes.Buffer) error {
	return nil
}

func (ctv CheckTxVersion) Decode(*bytes.Buffer) error { return nil }

func (ctv CheckTxVersion) Bytes() []byte {
	return sc.EncodedBytes(ctv)
}

func (ctv CheckTxVersion) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return sc.NewVaryingData(ctv.systemModule.Version().TransactionVersion), nil
}

func (_ CheckTxVersion) Validate(_who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (ctv CheckTxVersion) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.Pre, error) {
	_, err := ctv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (ctv CheckTxVersion) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) error {
	_, err := ctv.ValidateUnsigned(call, info, length)
	return err
}

func (ctv CheckTxVersion) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact[sc.Numeric], _result *primitives.DispatchResult) error {
	return nil
}

func (ctv CheckTxVersion) ModulePath() string {
	return systemModulePath
}
