package mocks

import "github.com/stretchr/testify/mock"

type IoTrie struct {
	mock.Mock
}

func (m *IoTrie) Blake2256OrderedRoot(key []byte, version int32) []byte {
	args := m.Called(key, version)

	return args[0].([]byte)
}
