package extensions

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonZeroAddress struct {
	typesInfoAdditionalSignedData sc.VaryingData
}

func NewCheckNonZeroAddress() primitives.SignedExtension {
	return &CheckNonZeroAddress{
		typesInfoAdditionalSignedData: sc.NewVaryingData(),
	}
}

func (c CheckNonZeroAddress) AdditionalSigned() (primitives.AdditionalSigned, error) {
	return primitives.AdditionalSigned{}, nil
}

func (c CheckNonZeroAddress) Encode(*bytes.Buffer) error {
	return nil
}

func (c CheckNonZeroAddress) Decode(*bytes.Buffer) error { return nil }

func (c CheckNonZeroAddress) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func (c CheckNonZeroAddress) Validate(who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	if reflect.DeepEqual(who, constants.ZeroAccountId) {
		return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadSigner())
	}

	return primitives.DefaultValidTransaction(), nil
}

func (c CheckNonZeroAddress) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.ValidTransaction, error) {
	return primitives.DefaultValidTransaction(), nil
}

func (c CheckNonZeroAddress) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) (primitives.Pre, error) {
	_, err := c.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (c CheckNonZeroAddress) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact[sc.Numeric]) error {
	_, err := c.ValidateUnsigned(call, info, length)
	return err
}

func (c CheckNonZeroAddress) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact[sc.Numeric], _result *primitives.DispatchResult) error {
	return nil
}

func (c CheckNonZeroAddress) ModulePath() string {
	return systemModulePath
}
