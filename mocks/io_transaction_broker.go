package mocks

import "github.com/stretchr/testify/mock"

type IoTransactionBroker struct {
	mock.Mock
}

func (m *IoTransactionBroker) Start() {
	m.Called()
}

func (m *IoTransactionBroker) Commit() {
	m.Called()
}

func (m *IoTransactionBroker) Rollback() {
	m.Called()
}
