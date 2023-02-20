package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/trie"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_BlockBuilder_Inherent_Extrinsics(t *testing.T) {
	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))

	assert.NoError(t, err)

	expectedExtrinsic := types.NewUnsignedUncheckedExtrinsic(types.Call{
		CallIndex: types.CallIndex{
			ModuleIndex:   timestamp.Module.Index,
			FunctionIndex: timestamp.Module.Functions["set"].Index,
		},
		Args: sc.BytesToSequenceU8(sc.ToCompact(time).Bytes()),
	})

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	storage := trie.NewEmptyTrie()
	rt := wasmer.NewTestInstanceWithTrie(t, WASM_RUNTIME, storage)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	assert.NotNil(t, inherentExt)

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := types.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic, extrinsic)
}
