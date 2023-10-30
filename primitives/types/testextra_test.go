package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type testExtraCheck struct {
	hasError sc.Bool
	value    sc.U32
}

func newTestExtraCheck(hasError sc.Bool, value sc.U32) SignedExtension {
	return &testExtraCheck{
		hasError: hasError,
		value:    value,
	}
}

func (e testExtraCheck) Encode(buffer *bytes.Buffer) {
	e.hasError.Encode(buffer)
	e.value.Encode(buffer)
}

func (e testExtraCheck) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheck) Decode(buffer *bytes.Buffer) {
	e.hasError = sc.DecodeBool(buffer)
	e.value = sc.DecodeU32(buffer)
}

func (e testExtraCheck) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
	if e.hasError {
		return nil, NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
	}

	return sc.NewVaryingData(e.value), nil
}

func (e testExtraCheck) Validate(who Address32, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	if e.hasError {
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
	if e.hasError {
		return NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
	}

	return nil
}

func (e testExtraCheck) Metadata() (MetadataType, MetadataSignedExtension) {
	id := 123456
	typ := 789
	docs := "TestExtraCheck"

	return NewMetadataTypeWithPath(
			id,
			docs,
			sc.Sequence[sc.Str]{"frame_system", "extensions", "test_extra_check", "TestExtraCheck"},
			NewMetadataTypeDefinitionCompact(sc.ToCompact(id)),
		),
		NewMetadataSignedExtension(sc.Str(docs), id, typ)
}