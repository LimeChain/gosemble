package mocks

import (
	"github.com/stretchr/testify/mock"
)

type IoHashing struct {
	mock.Mock
}

func (h *IoHashing) Twox128(value []byte) []byte {
	args := h.Called(value)
	return args.Get(0).([]byte)
}

func (h *IoHashing) Twox64(value []byte) []byte {
	args := h.Called(value)
	return args.Get(0).([]byte)
}

func (h *IoHashing) Blake128(value []byte) []byte {
	args := h.Called(value)
	return args.Get(0).([]byte)
}

func (h *IoHashing) Blake256(value []byte) []byte {
	args := h.Called(value)
	return args.Get(0).([]byte)
}
