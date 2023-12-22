package extensions

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	systemModulePath = "frame_system"
)

type CheckGenesis struct {
	module                        system.Module
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckGenesis(module system.Module) primitives.SignedExtension {
	return &CheckGenesis{
		module:                        module,
		typesInfoAdditionalSignedData: sc.NewVaryingData(primitives.H256{})}
}

func (cg CheckGenesis) Encode(*bytes.Buffer) error {
	return nil
}

func (cg CheckGenesis) Decode(*bytes.Buffer) error { return nil }

func (cg CheckGenesis) Bytes() []byte {
	return sc.EncodedBytes(cg)
}

func (cg CheckGenesis) AdditionalSigned() (primitives.AdditionalSigned, error) {
	hash, err := cg.module.StorageBlockHash(0)
	if err != nil {
		return nil, err
	}

	return sc.NewVaryingData(primitives.H256(hash)), nil
}

func (_ CheckGenesis) Validate(_who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (cg CheckGenesis) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (cg CheckGenesis) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.Pre, error) {
	_, err := cg.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (cg CheckGenesis) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) error {
	_, err := cg.ValidateUnsigned(call, info, length)
	return err
}

func (cg CheckGenesis) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact[sc.Numeric], _result *primitives.DispatchResult) error {
	return nil
}

func (cg CheckGenesis) ModulePath() string {
	return systemModulePath
}
