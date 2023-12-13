package aura

import (
	"encoding/json"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	gcJson      = []byte("{\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]}}")
	pubKey, _   = types.NewSr25519PublicKey(sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})...)
	authorities = sc.Sequence[types.Sr25519PublicKey]{pubKey}
)

func Test_GenesisConfig_UnmarshalJSON(t *testing.T) {

	auraGc := GenesisConfig{}
	err := json.Unmarshal(gcJson, &auraGc)
	assert.NoError(t, err)
	assert.Equal(t, authorities, auraGc.Authorities)
}

func Test_CreateDefaultConfig(t *testing.T) {
	setup(timestampMinimumPeriod)

	wantGc := []byte("{\"aura\":{\"authorities\":[]}}")

	gc, err := module.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}

func Test_BuildConfig(t *testing.T) {
	setup(timestampMinimumPeriod)
	mockStorageAuthorities.On("DecodeLen").Return(sc.NewOption[sc.U64](nil), nil)
	mockStorageAuthorities.On("Put", authorities).Return()

	err := module.BuildConfig(gcJson)
	assert.NoError(t, err)
	mockStorageAuthorities.AssertCalled(t, "Put", authorities)
}
