package main

import (
	"bytes"
	"testing"

	// "github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_JSON(t *testing.T) {
	rt, _ := newTestRuntime(t)
	gcStr := "{\"aura\":{\"authorities\":[]},\"balances\":{\"balances\":[]}}"

	res, err := rt.Exec("GenesisBuilder_create_default_config", []byte{})
	assert.NoError(t, err)

	bzDecoded, err := sc.DecodeSequence[sc.U8](bytes.NewBuffer(res))
	assert.NoError(t, err)

	assert.Equal(t, []byte(gcStr), sc.SequenceU8ToBytes(bzDecoded))
	// assert.Equal(t, [sc.Str(gcStr).Bytes()], res)
}

func Test_BuildConfig(t *testing.T) {
	rt, _ := newTestRuntime(t)

	gcStr := "{\"aura\": {\"authorities\":[\"testtest\",\"test2\"]},\"balances\":{\"balances\":[[\"5D34dL5prEUaGNQtPPZ3yN5Y6BnkfXunKXXz6fo7ZJbLwRRH\",100000000000000000],[\"5GBNeWRhZc2jXu7D55rBimKYDk8PGk8itRYFTPfC8RJLKG5o\",100000000000000000]]}}"

	_, err := rt.Exec("GenesisBuilder_build_config", sc.BytesToSequenceU8([]byte(gcStr)).Bytes())
	assert.NoError(t, err)
}

func Test_Encode_Decode(t *testing.T) {
	gcBz := []byte("{\"aura\": {\"authorities\":[\"testtest\",\"test2\"]},\"balances\":{\"balances\":[[\"5D34dL5prEUaGNQtPPZ3yN5Y6BnkfXunKXXz6fo7ZJbLwRRH\",100000000000000000],[\"5GBNeWRhZc2jXu7D55rBimKYDk8PGk8itRYFTPfC8RJLKG5o\",100000000000000000]]}}")
	gcSeqU8 := sc.BytesToSequenceU8(gcBz)
	buffEnc := bytes.Buffer{}
	err := gcSeqU8.Encode(&buffEnc)
	assert.NoError(t, err)
	assert.NotEqual(t, gcBz, buffEnc.Bytes())
	assert.Equal(t, buffEnc.Bytes(), gcSeqU8.Bytes())
	assert.Equal(t, gcBz, sc.SequenceU8ToBytes(gcSeqU8))

	buffDec := bytes.Buffer{}
	buffDec.Write(gcSeqU8.Bytes())
	gcDec, err := sc.DecodeSequence[sc.U8](&buffDec)
	assert.NoError(t, err)

	assert.Equal(t, gcSeqU8.Bytes(), gcDec.Bytes())
}
