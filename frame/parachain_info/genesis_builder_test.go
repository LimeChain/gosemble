package parachain_info

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	validGcJson        = "{\"parachainInfo\":{\"parachainId\":2000}}"
	parachainId sc.U32 = 2000
)

func Test_GenesisConfig_BuildConfig(t *testing.T) {
	target := setup()

	mockStorageParachainId.On("Put", parachainId).Return()

	err := target.BuildConfig([]byte(validGcJson))
	assert.NoError(t, err)

	mockStorageParachainId.AssertCalled(t, "Put", parachainId)
}

func Test_GenesisConfig_BuildConfig_EmptyBytes(t *testing.T) {
	target := setup()

	err := target.BuildConfig([]byte{})
	assert.EqualError(t, err, "unexpected end of JSON input")
}

func Test_GenesisConfig_BuildConfig_InvalidParachainId(t *testing.T) {
	target := setup()

	err := target.BuildConfig([]byte("{\"parachainInfo\":{\"parachainId\":\"error\"}}"))
	assert.EqualError(t, err, "json: cannot unmarshal string into Go struct field .parachainInfo.parachainId of type uint32")
}

func Test_GenesisConfig_CreateDefaultConfig(t *testing.T) {
	target := setup()

	expectedGc := []byte("{\"parachainInfo\":{\"parachainId\":100}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, expectedGc, gc)
}
