package main

import (
	"testing"

	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/stretchr/testify/assert"
)

func Test_Sign_Sr25519_Message(t *testing.T) {
	rt, _ := newTestRuntime(t)

	_, err := rt.Exec("SessionKeys_generate_session_keys", sc.NewOption[sc.U8](nil).Bytes())
	assert.NoError(t, err)

	keyTypeID := aura.KeyTypeId[:]
	publicKey := rt.Keystore().Aura.PublicKeys()[0].Encode()

	msg := []byte("test")
	msgArg, err := scale.Marshal(msg)
	assert.NoError(t, err)

	args := sc.BytesToSequenceU8(append(append(keyTypeID, publicKey...), msgArg...)).Bytes()
	res, err := rt.Exec("Example_sign_sr25519_message", args)
	assert.NoError(t, err)

	sig := res[1:]

	ok, err := rt.Keystore().Aura.Keypairs()[0].Public().Verify(msg, sig)
	assert.NoError(t, err)
	assert.True(t, ok)

	ok, err = rt.Keystore().Aura.Keypairs()[0].Public().Verify([]byte("invalid"), sig)
	assert.NoError(t, err)
	assert.False(t, ok)
}
