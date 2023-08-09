package main

import (
	"bytes"
	"testing"
	"time"

	gossamertypes "github.com/ChainSafe/gossamer/dot/types"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/execution/types"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/stretchr/testify/assert"
)

func Test_BlockBuilder_Inherent_Extrinsics(t *testing.T) {
	idata := gossamertypes.NewInherentData()
	time := time.Now().UnixMilli()
	err := idata.SetInherent(gossamertypes.Timstap0, uint64(time))
	decoder := types.NewModuleDecoder(modules, newSignedExtra())

	assert.NoError(t, err)

	call := tsm.NewSetCall(timestamp.ModuleIndex, timestamp.FunctionSetIndex, sc.NewVaryingData(sc.ToCompact(time)), nil, nil, nil)

	expectedExtrinsic := types.NewUnsignedUncheckedExtrinsic(call)

	ienc, err := idata.Encode()
	assert.NoError(t, err)

	rt, _ := newTestRuntime(t)

	inherentExt, err := rt.Exec("BlockBuilder_inherent_extrinsics", ienc)
	assert.NoError(t, err)

	assert.NotNil(t, inherentExt)

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{inherentExt[0]})

	totalInherents := sc.DecodeCompact(buffer)
	assert.Equal(t, int64(1), totalInherents.ToBigInt().Int64())
	buffer.Reset()

	buffer.Write(inherentExt[1:])
	extrinsic := decoder.DecodeUncheckedExtrinsic(buffer)

	assert.Equal(t, expectedExtrinsic.Bytes(), extrinsic.Bytes())
}
