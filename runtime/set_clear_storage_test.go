package main

import (
	"bytes"
	"testing"

	cscale "github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	ctypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/stretchr/testify/assert"
)

func Test_Set_key_in_storage(t *testing.T) {
	rt, storage := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	_, err = rt.Exec("Set_key_in_storage", extEnc.Bytes())
	assert.NoError(t, err)

	keyRes := (*storage).Get([]byte("Test"))

	expected := []byte("Set")

	assert.Equal(t, keyRes, expected)
}

func Test_Clear_key_in_storage(t *testing.T) {
	rt, storage := newTestRuntime(t)
	metadata := runtimeMetadata(t, rt)

	call, err := ctypes.NewCall(metadata, "System.remark", []byte{})
	assert.NoError(t, err)

	extrinsic := ctypes.NewExtrinsic(call)

	extEnc := bytes.Buffer{}
	encoder := cscale.NewEncoder(&extEnc)
	err = extrinsic.Encode(*encoder)
	assert.NoError(t, err)

	_, err = rt.Exec("Clear_key_in_storage", extEnc.Bytes())
	assert.NoError(t, err)

	keyRes := (*storage).Get([]byte("Test"))

	assert.Equal(t, []uint8([]byte(nil)), keyRes)
}
