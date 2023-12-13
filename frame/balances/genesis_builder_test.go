package balances

import (
	"encoding/json"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	gcJson   = []byte("{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}")
	accId, _ = types.NewAccountId(sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})...)
	balance  = sc.NewU128(uint64(1))
	balances = []gcAccountBalance{{AccountId: accId, Balance: balance}}
)

func Test_GenesisConfig_UnmarshalJSON(t *testing.T) {
	balancesGc := GenesisConfig{}
	err := json.Unmarshal(gcJson, &balancesGc)
	assert.NoError(t, err)
	assert.Equal(t, balances, balancesGc.Balances)
}

func Test_CreateDefaultConfig(t *testing.T) {
	target := setupModule()

	wantGc := []byte("{\"balances\":{\"balances\":[]}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}

func Test_BuildConfig(t *testing.T) {
	target := setupModule()
	mockTotalIssuance := new(mocks.StorageValue[sc.U128])
	target.storage.TotalIssuance = mockTotalIssuance

	mockTotalIssuance.On("Put", balance).Return()
	mockStoredMap.On("IncProviders", accId).Return(types.IncRefStatus(0), nil)
	mockStoredMap.On("Put", accId, types.AccountInfo{Data: types.AccountData{Free: balance}}).Return()

	err := target.BuildConfig(gcJson)
	assert.NoError(t, err)
	mockTotalIssuance.AssertCalled(t, "Put", balance)
	mockStoredMap.AssertCalled(t, "Put", accId, types.AccountInfo{Data: types.AccountData{Free: balance}})
	mockStoredMap.AssertCalled(t, "IncProviders", accId)
}
