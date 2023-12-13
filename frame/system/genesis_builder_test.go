package system

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	bytes69   = []byte{69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69}
	hash69, _ = types.NewBlake2bHash(sc.BytesToFixedSequenceU8(bytes69)...)
	lrui      = types.LastRuntimeUpgradeInfo{SpecVersion: 2, SpecName: "test-spec"}
)

func Test_CreateDefaultConfig(t *testing.T) {
	target := setupModule()

	wantGc := []byte("{\"system\":{}}")

	gc, err := target.CreateDefaultConfig()
	assert.NoError(t, err)
	assert.Equal(t, wantGc, gc)
}

func Test_BuildConfig(t *testing.T) {
	target := setupModule()
	mockStorageBlockHash.On("Put", sc.U64(0), hash69).Return()
	mockStorageParentHash.On("Put", hash69).Return()
	mockStorageLastRuntimeUpgrade.On("Put", lrui).Return()
	mockStorageExtrinsicIndex.On("Put", sc.U32(0)).Return()

	err := target.BuildConfig([]byte{})
	assert.NoError(t, err)
	mockStorageBlockHash.AssertCalled(t, "Put", sc.U64(0), hash69)
	mockStorageParentHash.AssertCalled(t, "Put", hash69)
	mockStorageLastRuntimeUpgrade.AssertCalled(t, "Put", lrui)
	mockStorageExtrinsicIndex.AssertCalled(t, "Put", sc.U32(0))
}
