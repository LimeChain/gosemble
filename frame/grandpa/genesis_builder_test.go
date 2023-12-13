package grandpa

import (
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	validGcJson            = "{\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}"
	accId, _               = types.NewAccountId(sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})...)
	authorities            = sc.Sequence[types.Authority]{{Id: accId, Weight: sc.U64(1)}}
	versionedAuthorityList = types.VersionedAuthorityList{AuthorityList: authorities, Version: AuthorityVersion}
)

func Test_GenesisConfig_BuildConfig(t *testing.T) {
	for _, tt := range []struct {
		name                     string
		gcJson                   string
		wantErr                  error
		shouldAssertCalled       bool
		storageAuthorities       types.VersionedAuthorityList
		storageAuthoritiesGetErr error
	}{
		{
			name:               "valid",
			gcJson:             validGcJson,
			shouldAssertCalled: true,
		},
		{
			name:               "duplicate genesis address",
			gcJson:             "{\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1],[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]}}",
			shouldAssertCalled: true,
		},
		{
			name:    "invalid genesis address",
			gcJson:  "{\"grandpa\":{\"authorities\":[[1,1]]}}",
			wantErr: errInvalidAddrValue,
		},
		{
			name:    "invalid ss58 address",
			gcJson:  "{\"grandpa\":{\"authorities\":[[\"invalid\",1]]}}",
			wantErr: errors.New("expected at least 2 bytes in base58 decoded address"),
		},
		{
			name:    "invalid genesis weight",
			gcJson:  "{\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",\"invalid\"]]}}",
			wantErr: errInvalidWeightValue,
		},
		{
			name:   "zero authorities",
			gcJson: "{\"grandpa\":{\"authorities\":[]}}",
		},
		{
			name:                     "storage authorities error on get",
			gcJson:                   validGcJson,
			storageAuthoritiesGetErr: errors.New("err"),
			wantErr:                  errors.New("err"),
		},
		{
			name:               "storage authorities already initialized",
			gcJson:             validGcJson,
			storageAuthorities: versionedAuthorityList,
			wantErr:            errAuthoritiesAlreadyInitialized,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			mockStorageAuthorities.On("Get").Return(tt.storageAuthorities, tt.storageAuthoritiesGetErr)
			mockStorageAuthorities.On("Put", versionedAuthorityList).Return()

			err := target.BuildConfig([]byte(tt.gcJson))
			assert.Equal(t, tt.wantErr, err)

			if tt.shouldAssertCalled {
				mockStorageAuthorities.AssertCalled(t, "Get")
				mockStorageAuthorities.AssertCalled(t, "Put", versionedAuthorityList)
			}
		})
	}
}

func Test_CreateDefaultConfig(t *testing.T) {
	setup()

	wantGc := []byte("{\"grandpa\":{\"authorities\":[]}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}
