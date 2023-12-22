package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckSpecVersion struct {
	systemModule                  system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckSpecVersion(systemModule system.Module) primitives.SignedExtension {
	return &CheckSpecVersion{
		systemModule:                  systemModule,
		typesInfoAdditionalSignedData: sc.NewVaryingData(sc.U32(0)),
	}
}

func (csv CheckSpecVersion) Encode(*bytes.Buffer) error {
	return nil
}

func (csv CheckSpecVersion) Decode(*bytes.Buffer) error { return nil }

func (csv CheckSpecVersion) Bytes() []byte {
	return sc.EncodedBytes(csv)
}

func (csv CheckSpecVersion) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return sc.NewVaryingData(csv.systemModule.Version().SpecVersion), nil
}

func (_ CheckSpecVersion) Validate(_who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (csv CheckSpecVersion) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (csv CheckSpecVersion) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, error) {
	_, err := csv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (csv CheckSpecVersion) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) error {
	_, err := csv.ValidateUnsigned(call, info, length)
	return err
}

func (csv CheckSpecVersion) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _dispatchErr error) error {
	return nil
}

func (csv CheckSpecVersion) ModulePath() string {
	return systemModulePath
}
