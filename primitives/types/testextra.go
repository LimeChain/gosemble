package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testExtraCheck struct {
	Values   []sc.Encodable
	HasError bool
}

func NewTestExtraCheck(hasError bool, values ...sc.Encodable) SignedExtension {
	return testExtraCheck{
		Values:   values,
		HasError: hasError,
	}
}

func (e testExtraCheck) Encode(buffer *bytes.Buffer) {
	for _, value := range e.Values {
		value.Encode(buffer)
	}
}

func (e testExtraCheck) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e testExtraCheck) Decode(buffer *bytes.Buffer) {}

func (e testExtraCheck) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	if e.HasError {
		return nil, NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
	}

	return sc.NewVaryingData(e.Values...), nil
}

func (e testExtraCheck) Validate(who Address32, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	if e.HasError {
		return ValidTransaction{}, NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
	}

	return validTransaction, nil
}

func (e testExtraCheck) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	return e.Validate(Address32{}, call, info, length)
}

func (e testExtraCheck) PreDispatch(who Address32, call Call, info *DispatchInfo, length sc.Compact) (Pre, TransactionValidityError) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheck) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) TransactionValidityError {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheck) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) TransactionValidityError {
	if e.HasError {
		return NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
	}

	return nil
}

func (e testExtraCheck) Metadata() (MetadataType, MetadataSignedExtension) {
	return MetadataType{}, MetadataSignedExtension{}
}
