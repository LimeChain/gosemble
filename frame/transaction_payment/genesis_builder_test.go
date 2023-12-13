package transaction_payment

import (
	"encoding/json"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	gcJson = []byte("{\"transactionPayment\":{\"multiplier\":\"1\"}}")
)

func Test_GenesisConfig_UnmarshalJSON(t *testing.T) {
	transactionPaymentGc := GenesisConfig{}
	err := json.Unmarshal(gcJson, &transactionPaymentGc)
	assert.NoError(t, err)
	assert.Equal(t, sc.NewU128(1), transactionPaymentGc.Multiplier)
}

func Test_CreateDefaultConfig(t *testing.T) {
	setup()

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, gcJson, gc)
}

func Test_BuildConfig(t *testing.T) {
	setup()
	mockNextFeeMultiplier.On("Put", sc.NewU128(1)).Return()

	err := target.BuildConfig(gcJson)
	assert.NoError(t, err)
	mockNextFeeMultiplier.AssertCalled(t, "Put", sc.NewU128(1))
}
