package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	mockAuraModule      *mocks.AuraModule
	mockSystemModule    *mocks.SystemModule
	mockGrandpaModule   *mocks.GrandpaModule
	mockTxPaymentModule *mocks.TransactionPaymentModule
	modules             = []Module{mockAuraModule, mockSystemModule, mockGrandpaModule, mockTxPaymentModule}
)

func setup() {
	mockAuraModule = new(mocks.AuraModule)
	mockSystemModule = new(mocks.SystemModule)
	mockGrandpaModule = new(mocks.GrandpaModule)
	mockTxPaymentModule = new(mocks.TransactionPaymentModule)
}

func Test_GetModule(t *testing.T) {
	setup()
	mockAuraModule.On("GetIndex").Return(sc.U8(0))
	mockAuraModule.On("GetIndex").Return(sc.U8(1))
	mockAuraModule.On("GetIndex").Return(sc.U8(2))
	mockAuraModule.On("GetIndex").Return(sc.U8(3))
	for i := 0; i < 4; i++ {
		m, ok := GetModule(0, modules)
		assert.True(t, ok)
		assert.Equal(t, m.GetIndex(), i)
	}
}
