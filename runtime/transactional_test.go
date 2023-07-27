package main

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/config"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/testable/dispatchables"
	"github.com/LimeChain/gosemble/frame/testable/module"

	"github.com/stretchr/testify/assert"
)

// TODO: in the test case "Commit_Then_Rollback" the host
// panics with "fatal error: exitsyscall: syscall frame is no longer valid"

func Test_Storage_Layer_Rollback_Then_Commit(t *testing.T) {
	rt, storage := newTestRuntime(t)

	call := dispatchables.NewTestCall(config.TestableIndex, module.FunctionTestIndex, sc.NewVaryingData(sc.Sequence[sc.U8]{}))

	extrinsic := types.NewUnsignedUncheckedExtrinsic(call)

	_, err := rt.Exec("BlockBuilder_apply_extrinsic", extrinsic.Bytes())
	assert.NoError(t, err)

	assert.Equal(t, []byte{1}, (*storage).Get([]byte("testvalue")))
}
