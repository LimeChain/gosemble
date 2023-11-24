package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testExtraCheckEmpty struct {
	additionalSignedData sc.VaryingData
}

func newTestExtraCheckEmpty() SignedExtension {
	return &testExtraCheckEmpty{
		additionalSignedData: sc.VaryingData{},
	}
}

func (e testExtraCheckEmpty) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckEmpty) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheckEmpty) Decode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckEmpty) AdditionalSigned() (AdditionalSigned, error) {
	return sc.NewVaryingData(), nil
}

func (e testExtraCheckEmpty) Validate(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	return validTransaction, nil
}

func (e testExtraCheckEmpty) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	return e.Validate(AccountId[PublicKey]{}, call, info, length)
}

func (e testExtraCheckEmpty) PreDispatch(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (Pre, error) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheckEmpty) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheckEmpty) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) error {
	return nil
}
