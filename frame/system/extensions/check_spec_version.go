package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckSpecVersion struct {
	systemModule system.Module
}

func NewCheckSpecVersion(systemModule system.Module) CheckSpecVersion {
	return CheckSpecVersion{
		systemModule: systemModule,
	}
}

func (csv CheckSpecVersion) Encode(*bytes.Buffer) {}

func (csv CheckSpecVersion) Decode(*bytes.Buffer) {}

func (csv CheckSpecVersion) Bytes() []byte {
	return sc.EncodedBytes(csv)
}

func (csv CheckSpecVersion) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return sc.NewVaryingData(csv.systemModule.Constants.Version.SpecVersion), nil
}

func (_ CheckSpecVersion) Validate(_who *primitives.Address32, _call *primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (csv CheckSpecVersion) ValidateUnsigned(_call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (csv CheckSpecVersion) PreDispatch(who *primitives.Address32, call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := csv.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (csv CheckSpecVersion) PreDispatchUnsigned(call *primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := csv.ValidateUnsigned(call, info, length)
	return err
}

func (csv CheckSpecVersion) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}
