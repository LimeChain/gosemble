package mocks

import (
	"github.com/stretchr/testify/mock"
)

type IoMisc struct {
	mock.Mock
}

func (m *IoMisc) PrintHex(data []byte) {
	m.Called(data)
}

func (m *IoMisc) PrintUtf8(data []byte) {
	m.Called(data)
}

func (m *IoMisc) RuntimeVersion(codeBlob []byte) []byte {
	args := m.Called(codeBlob)

	return args.Get(0).([]byte)
}
