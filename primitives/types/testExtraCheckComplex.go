package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testExtraCheckComplex struct {
	module               Module
	era                  Era
	additionalSignedData sc.VaryingData
}

func newtTestExtraCheckComplex() SignedExtension {
	return &testExtraCheckComplex{
		era:                  Era{},
		additionalSignedData: sc.VaryingData{H256{}},
	}
}

func (e testExtraCheckComplex) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckComplex) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheckComplex) Decode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckComplex) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	return sc.NewVaryingData(H256{}), nil
}

func (e testExtraCheckComplex) Validate(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	return validTransaction, nil
}

func (e testExtraCheckComplex) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	return e.Validate(AccountId[PublicKey]{}, call, info, length)
}

func (e testExtraCheckComplex) PreDispatch(who AccountId[PublicKey], call Call, info *DispatchInfo, length sc.Compact) (Pre, TransactionValidityError) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheckComplex) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) TransactionValidityError {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheckComplex) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) TransactionValidityError {
	return nil
}
