package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testExtraCheckEra struct {
	module               Module
	era                  Era
	additionalSignedData sc.VaryingData
}

func newtTestExtraCheckEra() SignedExtension {
	return &testExtraCheckEra{
		era:                  Era{},
		additionalSignedData: sc.VaryingData{H256{}},
	}
}

func (e testExtraCheckEra) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckEra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheckEra) Decode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckEra) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	return sc.NewVaryingData(H256{}), nil
}

func (e testExtraCheckEra) Validate(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	return validTransaction, nil
}

func (e testExtraCheckEra) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	return e.Validate(AccountId[PublicKey]{}, call, info, length)
}

func (e testExtraCheckEra) PreDispatch(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (Pre, TransactionValidityError) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheckEra) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) TransactionValidityError {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheckEra) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) TransactionValidityError {
	return nil
}
