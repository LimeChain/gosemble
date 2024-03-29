package aura

import (
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/stretchr/testify/assert"
)

var (
	validGcJson = "{\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]}}"
	pubKey, _   = types.NewSr25519PublicKey(sc.BytesToSequenceU8(signature.TestKeyringPairAlice.PublicKey)...)
	authorities = sc.Sequence[types.Sr25519PublicKey]{pubKey}
)

func Test_GenesisConfig_BuildConfig(t *testing.T) {
	for _, tt := range []struct {
		name               string
		gcJson             string
		expectedErr        error
		decodeLen          sc.Option[sc.U64]
		decodeLenErr       error
		maxAuthorities     sc.Option[sc.U32]
		shouldAssertCalled bool
	}{
		{
			name:               "valid",
			gcJson:             validGcJson,
			shouldAssertCalled: true,
			decodeLen:          sc.NewOption[sc.U64](nil),
		},
		{
			name:               "duplicate genesis address",
			gcJson:             "{\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\", \"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]}}",
			shouldAssertCalled: true,
			decodeLen:          sc.NewOption[sc.U64](nil),
		},
		{
			name:        "invalid ss58 address",
			gcJson:      "{\"aura\":{\"authorities\":[\"invalid\"]}}",
			expectedErr: errors.New("expected at least 2 bytes in base58 decoded address"),
		},
		{
			name:   "zero authorities",
			gcJson: "{\"aura\":{\"authorities\":[]}}",
		},
		{
			name:         "storage authorities DecodeLen error",
			gcJson:       validGcJson,
			decodeLenErr: errors.New("err"),
			expectedErr:  errors.New("err"),
		},
		{
			name:        "storage authorities DecodeLen has value",
			gcJson:      validGcJson,
			decodeLen:   sc.NewOption[sc.U64](sc.U64(1)),
			expectedErr: errAuthoritiesAlreadyInitialized,
		},
		{
			name:           "authorities exceed max authorities",
			gcJson:         validGcJson,
			maxAuthorities: sc.NewOption[sc.U32](sc.U32(0)),
			expectedErr:    errAuthoritiesExceedMaxAuthorities,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			setup(timestampMinimumPeriod)
			mockStorageAuthorities.On("DecodeLen").Return(tt.decodeLen, tt.decodeLenErr)
			mockStorageAuthorities.On("Put", authorities).Return()
			if tt.maxAuthorities.HasValue {
				module.config.MaxAuthorities = tt.maxAuthorities.Value
			}

			err := module.BuildConfig([]byte(tt.gcJson))
			assert.Equal(t, tt.expectedErr, err)

			if tt.shouldAssertCalled {
				mockStorageAuthorities.AssertCalled(t, "Put", authorities)
			}
		})
	}
}

func Test_CreateDefaultConfig(t *testing.T) {
	setup(timestampMinimumPeriod)

	expectedGc := []byte("{\"aura\":{\"authorities\":[]}}")

	gc, err := module.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, expectedGc, gc)
}
