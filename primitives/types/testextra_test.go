package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

var (
	unknownTransactionCustomUnknownTransaction = NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
)

type testExtraCheck struct {
	hasError             sc.Bool
	value                sc.U32
	additionalSignedData sc.VaryingData
}

func newTestExtraCheck(hasError sc.Bool, value sc.U32) SignedExtension {
	return &testExtraCheck{
		hasError:             hasError,
		value:                value,
		additionalSignedData: sc.VaryingData{sc.U32(0)},
	}
}

func (e testExtraCheck) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		e.hasError,
		e.value,
	)
}

func (e testExtraCheck) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheck) Decode(buffer *bytes.Buffer) error {
	hasError, err := sc.DecodeBool(buffer)
	if err != nil {
		return err
	}
	e.hasError = hasError
	value, err := sc.DecodeU32(buffer)
	if err != nil {
		return err
	}
	e.value = value
	return nil
}

func (e testExtraCheck) AdditionalSigned() (AdditionalSigned, error) {
	if e.hasError {
		return nil, unknownTransactionCustomUnknownTransaction
	}

	return sc.NewVaryingData(e.value), nil
}

func (e testExtraCheck) Validate(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	if e.hasError {
		return ValidTransaction{}, unknownTransactionCustomUnknownTransaction
	}

	return validTransaction, nil
}

func (e testExtraCheck) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	return e.Validate(AccountId{}, call, info, length)
}

func (e testExtraCheck) PreDispatch(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (Pre, error) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheck) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheck) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) error {
	if e.hasError {
		return unknownTransactionCustomUnknownTransaction
	}

	return nil
}

//func (e testExtraCheck) Metadata() (MetadataType, MetadataSignedExtension) {
//	id := 123456
//	typ := 789
//	docs := "TestExtraCheck"
//
//	return NewMetadataTypeWithPath(
//			id,
//			docs,
//			sc.Sequence[sc.Str]{"frame_system", "extensions", "test_extra_check", "TestExtraCheck"},
//			NewMetadataTypeDefinitionCompact(sc.ToCompact(id)),
//		),
//		NewMetadataSignedExtension(sc.Str(docs), id, typ)
//}
