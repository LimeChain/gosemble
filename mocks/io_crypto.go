package mocks

import "github.com/stretchr/testify/mock"

type IoCrypto struct {
	mock.Mock
}

func (m *IoCrypto) EcdsaGenerate(keyTypeId []byte, seed []byte) []byte {
	args := m.Called(keyTypeId, seed)
	return args.Get(0).([]byte)
}

func (m *IoCrypto) EcdsaRecoverCompressed(signature []byte, msg []byte) []byte {
	args := m.Called(signature, msg)
	return args.Get(0).([]byte)
}

func (m *IoCrypto) Ed25519Generate(keyTypeId []byte, seed []byte) []byte {
	args := m.Called(keyTypeId, seed)

	return args.Get(0).([]byte)
}

func (m *IoCrypto) Ed25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	args := m.Called(signature, message, pubKey)

	return args.Get(0).(bool)
}

func (m *IoCrypto) Sr25519Generate(keyTypeId []byte, seed []byte) []byte {
	args := m.Called(keyTypeId, seed)

	return args.Get(0).([]byte)
}

func (m *IoCrypto) Sr25519Verify(signature []byte, message []byte, pubKey []byte) bool {
	args := m.Called(signature, message, pubKey)

	return args.Get(0).(bool)
}
