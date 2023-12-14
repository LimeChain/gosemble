package balances

import (
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

var (
	validGcJson             = "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}"
	accId, _                = types.NewAccountId(sc.BytesToSequenceU8(signature.TestKeyringPairAlice.PublicKey)...)
	balanceOne              = sc.NewU128(uint64(1))
	balanceOverMaxUint64, _ = sc.NewU128FromString("184467440737095516150")
)

func Test_GenesisConfig_BuildConfig(t *testing.T) {
	for _, tt := range []struct {
		name                     string
		gcJson                   string
		wantErr                  error
		shouldAssertCalled       bool
		tryMutateExistsErr       error
		tryMutateExistsResultErr error
		balance                  sc.U128
	}{
		{
			name:               "valid",
			gcJson:             validGcJson,
			balance:            balanceOne,
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
			name:               "balance greater than MaxUint64",
			gcJson:             "{\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",184467440737095516150]]}}",
			balance:            balanceOverMaxUint64,
			shouldAssertCalled: true,
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
			name:               "TryMutateExists error",
			gcJson:             validGcJson,
			tryMutateExistsErr: errors.New("err"),
			wantErr:            errors.New("err"),
		},
		{
			name:                     "TryMutateExists result errror",
			gcJson:                   validGcJson,
			tryMutateExistsResultErr: errors.New("err"),
			wantErr:                  types.NewDispatchErrorOther(sc.Str(errors.New("err").Error())),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			target := setupModule()
			mockTotalIssuance := new(mocks.StorageValue[sc.U128])
			target.storage.TotalIssuance = mockTotalIssuance

			mockResult := sc.Result[sc.Encodable]{}
			if tt.tryMutateExistsResultErr != nil {
				mockResult = sc.Result[sc.Encodable]{
					HasError: true,
					Value:    types.NewDispatchErrorOther(sc.Str(tt.tryMutateExistsResultErr.Error())),
				}
			}

			mockStoredMap.On("TryMutateExists", accId, mockTypeMutateAccountData).Return(mockResult, tt.tryMutateExistsErr)
			mockTotalIssuance.On("Put", tt.balance).Return()

			err := target.BuildConfig([]byte(tt.gcJson))
			assert.Equal(t, tt.wantErr, err)

			if tt.shouldAssertCalled {
				mockTotalIssuance.AssertCalled(t, "Put", tt.balance)
				mockStoredMap.AssertCalled(t, "TryMutateExists", accId, mockTypeMutateAccountData)
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
