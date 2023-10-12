package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type RuntimeExtrinsic struct {
	mock.Mock
}

func (re *RuntimeExtrinsic) Module(index sc.U8) (module primitives.Module, isFound bool) {
	args := re.Called(index)
	return args.Get(0).(primitives.Module), args.Bool(1)
}

func (re *RuntimeExtrinsic) CreateInherents(inherentData primitives.InherentData) []byte {
	args := re.Called(inherentData)
	return args.Get(0).([]byte)
}

func (re *RuntimeExtrinsic) CheckInherents(data primitives.InherentData, block primitives.Block) primitives.CheckInherentsResult {
	args := re.Called(data, block)
	return args.Get(0).(primitives.CheckInherentsResult)
}

func (re *RuntimeExtrinsic) EnsureInherentsAreFirst(block types.Block) int {
	args := re.Called(block)
	return args.Int(0)
}

func (re *RuntimeExtrinsic) OnInitialize(n sc.U64) primitives.Weight {
	args := re.Called(n)
	return args.Get(0).(primitives.Weight)
}

func (re *RuntimeExtrinsic) OnRuntimeUpgrade() primitives.Weight {
	args := re.Called()
	return args.Get(0).(primitives.Weight)
}

func (re *RuntimeExtrinsic) OnFinalize(n sc.U64) {
	re.Called(n)
}

func (re *RuntimeExtrinsic) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	args := re.Called(n, remainingWeight)
	return args.Get(0).(primitives.Weight)
}

func (re *RuntimeExtrinsic) OffchainWorker(n sc.U64) {
	re.Called(n)
}

func (re *RuntimeExtrinsic) Metadata() (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModule], primitives.MetadataExtrinsic) {
	args := re.Called()
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(sc.Sequence[primitives.MetadataModule]), args.Get(2).(primitives.MetadataExtrinsic)
}
