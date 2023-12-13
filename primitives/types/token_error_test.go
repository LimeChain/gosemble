package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_TokenError(t *testing.T) {
	for _, tt := range []struct {
		name       string
		newErr     TokenError
		wantErr    error
		wantErrMsg string
	}{
		{
			name:       "TokenErrorNoFunds",
			newErr:     NewTokenErrorNoFunds(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorNoFunds)),
			wantErrMsg: "Funds are unavailable",
		},
		{
			name:       "TokenErrorWouldDie",
			newErr:     NewTokenErrorWouldDie(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorWouldDie)),
			wantErrMsg: "Account that must exist would die",
		},
		{
			name:       "TokenErrorBelowMinimum",
			newErr:     NewTokenErrorBelowMinimum(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorBelowMinimum)),
			wantErrMsg: "Account cannot exist with the funds that would be given",
		},
		{
			name:       "TokenErrorCannotCreate",
			newErr:     NewTokenErrorCannotCreate(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorCannotCreate)),
			wantErrMsg: "Account cannot be created",
		},
		{
			name:       "TokenErrorUnknownAsset",
			newErr:     NewTokenErrorUnknownAsset(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorUnknownAsset)),
			wantErrMsg: "The asset in question is unknown",
		},
		{
			name:       "TokenErrorFrozen",
			newErr:     NewTokenErrorFrozen(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorFrozen)),
			wantErrMsg: "Funds exist but are frozen",
		},
		{
			name:       "TokenErrorUnsupported",
			newErr:     NewTokenErrorUnsupported(),
			wantErr:    TokenError(sc.NewVaryingData(TokenErrorUnsupported)),
			wantErrMsg: "Operation is not supported by the asset",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			err := tt.newErr.Encode(buffer)
			assert.NoError(t, err)

			haveErr, err := DecodeTokenError(buffer)
			assert.NoError(t, err)

			assert.Equal(t, tt.wantErr, haveErr)
			assert.Equal(t, tt.wantErrMsg, haveErr.Error())
		})
	}
}

func Test_DecodeTokenError_TypeError(t *testing.T) {
	for _, tt := range []struct {
		name    string
		errType sc.Encodable
	}{
		{
			name:    "invalid type",
			errType: sc.U8(7),
		},
		{
			name:    "nil",
			errType: sc.Empty{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			err := tt.errType.Encode(buffer)
			assert.NoError(t, err)

			_, err = DecodeTokenError(buffer)
			assert.Error(t, err)
			assert.Equal(t, "not a valid 'TokenError' type", err.Error())
		})
	}
}
