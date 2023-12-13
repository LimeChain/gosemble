package types

import (
	"bytes"
	"reflect"
	"strings"

	sc "github.com/LimeChain/goscale"
)

const (
	basePath = "github.com/LimeChain/gosemble/"
)

// a check that has multiple varying signed data
type testExtraCheckComplex struct {
	module                        Module
	era                           Era
	hash                          H256
	value                         sc.U64
	typesInfoAdditionalSignedData sc.VaryingData
}

func newtTestExtraCheckComplex() SignedExtension {
	return &testExtraCheckComplex{
		era:                           Era{},
		typesInfoAdditionalSignedData: sc.NewVaryingData(H256{}, sc.U32(0), sc.U64(0), H512{}, Ed25519PublicKey{}, Weight{}),
	}
}

func (e testExtraCheckComplex) Encode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckComplex) ModulePath() string {
	pkgPath := reflect.TypeOf(e).PkgPath()
	_, pkgPath, _ = strings.Cut(pkgPath, basePath)
	pkgPath, _, _ = strings.Cut(pkgPath, "/extensions")
	return strings.Replace(pkgPath, "/", "_", 1)
}

func (e testExtraCheckComplex) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e *testExtraCheckComplex) Decode(buffer *bytes.Buffer) error {
	return nil
}

func (e testExtraCheckComplex) AdditionalSigned() (AdditionalSigned, error) {
	return sc.NewVaryingData(H256{}, sc.U32(0), sc.U64(0), H512{}, Ed25519PublicKey{}), nil
}

func (e testExtraCheckComplex) Validate(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	validTransaction := DefaultValidTransaction()
	validTransaction.Priority = 1

	return validTransaction, nil
}

func (e testExtraCheckComplex) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	return e.Validate(AccountId{}, call, info, length)
}

func (e testExtraCheckComplex) PreDispatch(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (Pre, error) {
	_, err := e.Validate(who, call, info, length)
	return Pre{}, err
}

func (e testExtraCheckComplex) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error {
	_, err := e.ValidateUnsigned(call, info, length)
	return err
}

func (e testExtraCheckComplex) PostDispatch(pre sc.Option[Pre], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) error {
	return nil
}
