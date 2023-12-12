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

func (re *RuntimeExtrinsic) CreateInherents(inherentData primitives.InherentData) ([]byte, error) {
	args := re.Called(inherentData)
	if args.Get(1) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (re *RuntimeExtrinsic) CheckInherents(data primitives.InherentData, block primitives.Block) (primitives.CheckInherentsResult, error) {
	args := re.Called(data, block)

	if args.Get(1) == nil {
		return args.Get(0).(primitives.CheckInherentsResult), nil
	}

	return args.Get(0).(primitives.CheckInherentsResult), args.Get(1).(error)

}

func (re *RuntimeExtrinsic) EnsureInherentsAreFirst(block types.Block) int {
	args := re.Called(block)
	return args.Int(0)
}

func (re *RuntimeExtrinsic) OnInitialize(n sc.U64) (primitives.Weight, error) {
	args := re.Called(n)
	if args.Get(1) == nil {
		return args.Get(0).(primitives.Weight), nil
	}
	return args.Get(0).(primitives.Weight), args.Get(1).(error)
}

func (re *RuntimeExtrinsic) OnRuntimeUpgrade() primitives.Weight {
	args := re.Called()
	return args.Get(0).(primitives.Weight)
}

func (re *RuntimeExtrinsic) OnFinalize(n sc.U64) error {
	args := re.Called(n)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (re *RuntimeExtrinsic) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	args := re.Called(n, remainingWeight)
	return args.Get(0).(primitives.Weight)
}

func (re *RuntimeExtrinsic) OffchainWorker(n sc.U64) {
	re.Called(n)
}

func (re *RuntimeExtrinsic) Metadata(mdGenerator *primitives.MetadataGenerator) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV14], primitives.MetadataExtrinsicV14) {
	args := re.Called(mdGenerator)
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(sc.Sequence[primitives.MetadataModuleV14]), args.Get(2).(primitives.MetadataExtrinsicV14)
}

func (re *RuntimeExtrinsic) MetadataLatest(mdGenerator *primitives.MetadataGenerator) (sc.Sequence[primitives.MetadataType], sc.Sequence[primitives.MetadataModuleV15], primitives.MetadataExtrinsicV15, primitives.OuterEnums, primitives.CustomMetadata) {
	args := re.Called(mdGenerator)
	return args.Get(0).(sc.Sequence[primitives.MetadataType]), args.Get(1).(sc.Sequence[primitives.MetadataModuleV15]), args.Get(2).(primitives.MetadataExtrinsicV15), args.Get(3).(primitives.OuterEnums), args.Get(4).(primitives.CustomMetadata)
}
