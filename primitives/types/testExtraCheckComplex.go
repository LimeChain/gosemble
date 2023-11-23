package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// a check that has multiple varying signed data
type testExtraCheckComplex struct {
	module               Module
	era                  Era
	hash                 H256
	value                sc.U64
	additionalSignedData sc.VaryingData
}

func newtTestExtraCheckComplex() SignedExtension {
	return &testExtraCheckComplex{
		era:                  Era{},
		additionalSignedData: sc.NewVaryingData(H256{}, sc.U32(0), sc.U64(0), H512{}, Ed25519PublicKey{}),
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
	return sc.NewVaryingData(H256{}, sc.U32(0), sc.U64(0), H512{}, Ed25519PublicKey{}), nil
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