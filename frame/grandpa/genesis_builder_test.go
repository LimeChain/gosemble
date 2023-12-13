package grandpa

import (
	"encoding/json"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	gcJson                 = []byte("{\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}")
	pubKey, _              = types.NewEd25519PublicKey(sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})...)
	accId                  = types.NewAccountId[types.PublicKey](pubKey)
	authorities            = sc.Sequence[types.Authority]{{Id: accId, Weight: sc.U64(1)}}
	versionedAuthorityList = types.VersionedAuthorityList{AuthorityList: authorities, Version: AuthorityVersion}
)

func Test_GenesisConfig_UnmarshalJSON(t *testing.T) {
	grandpaGc := GenesisConfig{}
	err := json.Unmarshal(gcJson, &grandpaGc)
	assert.NoError(t, err)
	assert.Equal(t, authorities, grandpaGc.Authorities)
}

func Test_CreateDefaultConfig(t *testing.T) {
	setup()

	wantGc := []byte("{\"grandpa\":{\"authorities\":[]}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}

func Test_BuildConfig(t *testing.T) {
	setup()
	mockStorageAuthorities.On("Get").Return(types.VersionedAuthorityList{}, nil)
	mockStorageAuthorities.On("Put", versionedAuthorityList).Return()

	err := target.BuildConfig(gcJson)
	assert.NoError(t, err)
	mockStorageAuthorities.AssertCalled(t, "Get")
	mockStorageAuthorities.AssertCalled(t, "Put", versionedAuthorityList)
}
