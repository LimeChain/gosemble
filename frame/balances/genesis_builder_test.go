package balances

import (
	// "encoding/json"
	// "errors"
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	validGcJson = "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}"
	accId, _    = types.NewAccountId(sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})...)
	balance     = sc.NewU128(uint64(1))
)

func Test_GenesisConfig_BuildConfig(t *testing.T) {
	for _, tt := range []struct {
		name               string
		gcJson             string
		wantErr            error
		shouldAssertCalled bool
		incProvidersErr    error
	}{
		{
			name:               "valid",
			gcJson:             validGcJson,
			shouldAssertCalled: true,
		},
		{
			name:    "invalid genesis address",
			gcJson:  "{\"balances\":{\"balances\":[[1,1]]}}",
			wantErr: errInvalidAddrValue,
		},
		{
			name:    "duplicate genesis balance",
			gcJson:  "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1],[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}",
			wantErr: errDuplicateBalancesInGenesis,
		},
		{
			name:    "invalid ss58 address",
			gcJson:  "{\"balances\":{\"balances\":[[\"invalid\",1]]}}",
			wantErr: errors.New("expected at least 2 bytes in base58 decoded address"),
		},
		{
			name:    "invalid genesis balance",
			gcJson:  "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",\"invalid\"]]}}",
			wantErr: errInvalidBalanceValue,
		},
		{
			name:   "zero balances",
			gcJson: "{\"aura\":{\"authorities\":[]}}",
		},
		{
			name:    "balance below existential deposit",
			gcJson:  "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",0]]}}",
			wantErr: errBalanceBelowExistentialDeposit,
		},
		{
			name:            "inc providers error",
			gcJson:          validGcJson,
			incProvidersErr: errors.New("err"),
			wantErr:         errors.New("err"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			target := setupModule()
			mockTotalIssuance := new(mocks.StorageValue[sc.U128])
			target.storage.TotalIssuance = mockTotalIssuance

			mockTotalIssuance.On("Put", balance).Return()
			mockStoredMap.On("IncProviders", accId).Return(types.IncRefStatus(0), tt.incProvidersErr)
			mockStoredMap.On("Put", accId, types.AccountInfo{Data: types.AccountData{Free: balance}}).Return()

			err := target.BuildConfig([]byte(tt.gcJson))
			assert.Equal(t, tt.wantErr, err)

			if tt.shouldAssertCalled {
				mockTotalIssuance.AssertCalled(t, "Put", balance)
				mockStoredMap.AssertCalled(t, "Put", accId, types.AccountInfo{Data: types.AccountData{Free: balance}})
				mockStoredMap.AssertCalled(t, "IncProviders", accId)
			}
		})
	}
}

func Test_GenesisConfig_CreateDefaultConfig(t *testing.T) {
	target := setupModule()

	wantGc := []byte("{\"balances\":{\"balances\":[]}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}
