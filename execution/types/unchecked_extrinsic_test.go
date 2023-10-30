package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	version = sc.U8(5)
)

var (
	unknownTransactionCannotLookupError = types.NewTransactionValidityError(
		types.NewUnknownTransactionCannotLookup(),
	)
	invalidTransactionAncientBirthBlockError = types.NewTransactionValidityError(
		types.NewInvalidTransactionAncientBirthBlock(),
	)
	invalidTransactionBadProofError = types.NewTransactionValidityError(
		types.NewInvalidTransactionBadProof(),
	)

	signerAddressBytes = []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	signerAddress = types.NewAddress32(sc.BytesToSequenceU8(signerAddressBytes)...)
	signer        = types.NewMultiAddressId(types.AccountId{Address32: signerAddress})

	signatureBytes = []byte{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1,
	}
	signatureEd25519 = types.NewMultiSignatureEd25519(
		types.NewSignatureEd25519(
			sc.BytesToFixedSequenceU8(signatureBytes)...,
		),
	)
	signatureSr25519 = types.NewMultiSignatureSr25519(
		types.NewSignatureSr25519(
			sc.BytesToFixedSequenceU8(signatureBytes)...,
		),
	)

	unknownMultisignature = types.MultiSignature{
		VaryingData: sc.NewVaryingData(sc.U8(3), signatureEd25519),
	}

	additionalSigned = sc.NewVaryingData(
		types.H256{
			FixedSequence: sc.FixedSequence[sc.U8]{
				0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
				0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
				0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37, 0x37,
				0x37, 0x37,
			},
		},
	)

	encodedPayloadBytes = []byte{0x38, 0x38, 0x38}
)

var (
	targetSigned   uncheckedExtrinsic
	targetUnsigned uncheckedExtrinsic

	extrinsicSignature sc.Option[types.ExtrinsicSignature]

	mockCall            *mocks.Call
	mockSignedExtra     *mocks.SignedExtra
	mockAccountIdLookup *mocks.AccountIdLookup
	mocksSignedPayload  *mocks.SignedPayload
	mockCrypto          *mocks.IoCrypto
	mockHashing         *mocks.IoHashing
)

func setup(signature types.MultiSignature) {
	mockCall = new(mocks.Call)
	mockSignedExtra = new(mocks.SignedExtra)
	mockAccountIdLookup = new(mocks.AccountIdLookup)
	mocksSignedPayload = new(mocks.SignedPayload)
	mockCrypto = new(mocks.IoCrypto)
	mockHashing = new(mocks.IoHashing)

	extrinsicSignature = sc.NewOption[types.ExtrinsicSignature](
		types.ExtrinsicSignature{
			Signer:    signer,
			Signature: signature,
			Extra:     mockSignedExtra,
		},
	)

	targetUnsigned = newTestUnsignedExtrinsic(mockCall)

	targetSigned = newTestSignedExtrinsic(
		extrinsicSignature,
		mockCall,
		mockSignedExtra,
		mocksSignedPayload,
		mockCrypto,
		mockHashing,
	)
}

func newTestUnsignedExtrinsic(call types.Call) uncheckedExtrinsic {
	return NewUnsignedUncheckedExtrinsic(call).(uncheckedExtrinsic)
}

func newTestSignedExtrinsic(
	signature sc.Option[types.ExtrinsicSignature],
	call types.Call,
	extra types.SignedExtra,
	signedPayload types.SignedPayload,
	crypto io.Crypto,
	hashing io.Hashing) uncheckedExtrinsic {

	initializer := func(call types.Call, extra types.SignedExtra) (types.SignedPayload, types.TransactionValidityError) {
		return signedPayload, nil
	}

	uxt := NewUncheckedExtrinsic(version, signature, call, extra).(uncheckedExtrinsic)
	uxt.initializePayload = initializer
	uxt.crypto = crypto
	uxt.hashing = hashing

	return uxt
}

func Test_Encode_UncheckedExtrinsic_Unsigned(t *testing.T) {
	setup(signatureEd25519)

	buffer := &bytes.Buffer{}
	mockCall.On("Encode", mock.Anything)

	targetUnsigned.Encode(buffer)

	mockCall.AssertCalled(t, "Encode", mock.Anything)
	mockSignedExtra.AssertNotCalled(t, "Encode")
	assert.Equal(t, []byte{0x4, 0x4}, buffer.Bytes())
}

func Test_Encode_UncheckedExtrinsic_Signed(t *testing.T) {
	setup(signatureEd25519)

	buffer := &bytes.Buffer{}
	mockCall.On("Encode", mock.Anything)
	mockSignedExtra.On("Encode", mock.Anything)

	targetSigned.Encode(buffer)

	mockCall.AssertCalled(t, "Encode", mock.Anything)
	mockSignedExtra.AssertCalled(t, "Encode", mock.Anything)
	assert.Equal(t, []byte{
		0x8d, 0x1, // length
		5,                                                                                                 // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		// call
	}, buffer.Bytes())
}

func Test_Bytes_UncheckedExtrinsic_Unsigned(t *testing.T) {
	setup(signatureEd25519)

	mockCall.On("Encode", mock.Anything)

	encoded := targetUnsigned.Bytes()

	mockCall.AssertCalled(t, "Encode", mock.Anything)
	assert.Equal(t, []byte{0x4, 0x4}, encoded)
}

func Test_Bytes_UncheckedExtrinsic_Signed(t *testing.T) {
	setup(signatureEd25519)

	mockCall.On("Encode", mock.Anything)
	mockSignedExtra.On("Encode", mock.Anything)

	encoded := targetSigned.Bytes()

	mockCall.AssertCalled(t, "Encode", mock.Anything)
	mockSignedExtra.AssertCalled(t, "Encode", mock.Anything)
	assert.Equal(t, []byte{
		0x8d, 0x1, // length
		5,                                                                                                 // version
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // signer
		0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, // signature,
		// extra
		// call
	}, encoded)
}

func Test_Signature(t *testing.T) {
	setup(signatureEd25519)

	assert.Equal(t, extrinsicSignature, targetSigned.Signature())
}

func Test_Function(t *testing.T) {
	setup(signatureEd25519)

	assert.Equal(t, mockCall, targetSigned.Function())
}

func Test_Extra(t *testing.T) {
	setup(signatureEd25519)

	assert.Equal(t, mockSignedExtra, targetSigned.Extra())
}

func Test_IsSigned(t *testing.T) {
	setup(signatureEd25519)

	assert.Equal(t, false, targetUnsigned.IsSigned())
	assert.Equal(t, true, targetSigned.IsSigned())
}

func Test_Check_UnsignedUncheckedExtrinsic(t *testing.T) {
	setup(signatureEd25519)
	expect := NewCheckedExtrinsic(sc.NewOption[types.Address32](nil), mockCall, types.SignedExtra(nil)).(checkedExtrinsic)

	lookup := types.DefaultAccountIdLookup()

	result, err := targetUnsigned.Check(lookup)

	assert.Nil(t, err)
	checked := result.(checkedExtrinsic)
	assert.Equal(t, expect.extra, checked.extra)
	assert.Equal(t, expect.signer, checked.signer)
	assert.Equal(t, expect.function, checked.function)
}

func Test_Check_SignedUncheckedExtrinsic_LookupError(t *testing.T) {
	setup(signatureEd25519)

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).
		Return(types.Address32{}, unknownTransactionCannotLookupError)

	res, err := targetSigned.Check(mockAccountIdLookup)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mockSignedExtra.AssertNotCalled(t, "AdditionalSigned")
	mocksSignedPayload.AssertNotCalled(t, "UsingEncoded")
	mockCrypto.AssertNotCalled(t, "Ed25519Verify", mock.Anything, mock.Anything, mock.Anything)
	assert.Equal(t, unknownTransactionCannotLookupError, err)
	assert.Equal(t, nil, res)
}

func Test_Check_SignedUncheckedExtrinsic_AncientBirthBlockError(t *testing.T) {
	setup(signatureEd25519)

	targetSigned.initializePayload = types.NewSignedPayload
	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mockSignedExtra.On("AdditionalSigned").Return(types.AdditionalSigned{}, invalidTransactionAncientBirthBlockError)

	res, err := targetSigned.Check(mockAccountIdLookup)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mockSignedExtra.AssertCalled(t, "AdditionalSigned")
	mocksSignedPayload.AssertNotCalled(t, "UsingEncoded")
	mockCrypto.AssertNotCalled(t, "Ed25519Verify", mock.Anything, mock.Anything, mock.Anything)
	assert.Equal(t, invalidTransactionAncientBirthBlockError, err)
	assert.Equal(t, nil, res)
}

func Test_Check_SignedUncheckedExtrinsic_BadProofError(t *testing.T) {
	setup(signatureEd25519)

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mocksSignedPayload.On("Bytes").Return(encodedPayloadBytes)
	mockCrypto.On("Ed25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes).Return(false)

	res, err := targetSigned.Check(mockAccountIdLookup)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mocksSignedPayload.AssertCalled(t, "Bytes")
	mockHashing.AssertNotCalled(t, "Blake256", mock.Anything)
	mockCrypto.AssertCalled(t, "Ed25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes)
	assert.Equal(t, invalidTransactionBadProofError, err)
	assert.Equal(t, nil, res)
}

func Test_Check_SignedUncheckedExtrinsic_LongEncoding_BadProofError(t *testing.T) {
	setup(signatureEd25519)

	signedPayloadBytes := make([]byte, 257)
	blakeHashBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mocksSignedPayload.On("Bytes").Return(signedPayloadBytes)
	mockHashing.On("Blake256", signedPayloadBytes).Return(blakeHashBytes)
	mockCrypto.On("Ed25519Verify", signatureBytes, blakeHashBytes, signerAddressBytes).Return(false)

	res, err := targetSigned.Check(mockAccountIdLookup)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mocksSignedPayload.AssertCalled(t, "Bytes")
	mockHashing.AssertCalled(t, "Blake256", signedPayloadBytes)
	mockCrypto.AssertCalled(t, "Ed25519Verify", signatureBytes, blakeHashBytes, signerAddressBytes)
	assert.Equal(t, invalidTransactionBadProofError, err)
	assert.Equal(t, nil, res)
}

func Test_Check_SignedUncheckedExtrinsic_Success(t *testing.T) {
	setup(signatureEd25519)
	expect := NewCheckedExtrinsic(sc.NewOption[types.Address32](signerAddress), mockCall, mockSignedExtra).(checkedExtrinsic)

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mocksSignedPayload.On("Bytes").Return(encodedPayloadBytes)
	mockCrypto.On("Ed25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes).Return(true)

	result, err := targetSigned.Check(mockAccountIdLookup)

	assert.Nil(t, err)
	checked := result.(checkedExtrinsic)
	assert.Equal(t, expect.extra, checked.extra)
	assert.Equal(t, expect.signer, checked.signer)
	assert.Equal(t, expect.function, checked.function)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mocksSignedPayload.AssertCalled(t, "Bytes")
	mockHashing.AssertNotCalled(t, "Blake256", mock.Anything)
	mockCrypto.AssertCalled(t, "Ed25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes)
}

func Test_Check_SignedUncheckedExtrinsic_Success_Sr25519(t *testing.T) {
	setup(signatureSr25519)
	expect := NewCheckedExtrinsic(sc.NewOption[types.Address32](signerAddress), mockCall, mockSignedExtra).(checkedExtrinsic)

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mocksSignedPayload.On("Bytes").Return(encodedPayloadBytes)
	mockCrypto.On("Sr25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes).Return(true)

	result, err := targetSigned.Check(mockAccountIdLookup)

	assert.Nil(t, err)
	checked := result.(checkedExtrinsic)
	assert.Equal(t, expect.extra, checked.extra)
	assert.Equal(t, expect.signer, checked.signer)
	assert.Equal(t, expect.function, checked.function)

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mocksSignedPayload.AssertCalled(t, "Bytes")
	mockHashing.AssertNotCalled(t, "Blake256", mock.Anything)
	mockCrypto.AssertCalled(t, "Sr25519Verify", signatureBytes, encodedPayloadBytes, signerAddressBytes)
}

func Test_Check_SignedUncheckedExtrinsic_UnknownSignatureType(t *testing.T) {
	setup(unknownMultisignature)

	mockAccountIdLookup.On("Lookup", extrinsicSignature.Value.Signer).Return(signerAddress, nil)
	mocksSignedPayload.On("Bytes").Return(encodedPayloadBytes)

	assert.PanicsWithValue(t, "invalid MultiSignature type in Verify", func() {
		targetSigned.Check(mockAccountIdLookup)
	})

	mockAccountIdLookup.AssertCalled(t, "Lookup", extrinsicSignature.Value.Signer)
	mocksSignedPayload.AssertCalled(t, "Bytes")
	mockHashing.AssertNotCalled(t, "Blake256", mock.Anything)
	mockCrypto.AssertNotCalled(t, "Sr25519Verify", mock.Anything, mock.Anything, mock.Anything)
}
