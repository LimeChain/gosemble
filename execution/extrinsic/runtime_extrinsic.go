package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeExtrinsic struct {
	modules map[sc.U8]types.Module
}

func New(modules map[sc.U8]types.Module) RuntimeExtrinsic {
	return RuntimeExtrinsic{modules: modules}
}

func (re RuntimeExtrinsic) Module(index sc.U8) (module types.Module, isFound bool) {
	m, ok := re.modules[index]
	return m, ok
}

func (re RuntimeExtrinsic) CreateInherents(inherentData primitives.InherentData) []byte {
	i := 0
	var result []byte

	for _, module := range re.modules {
		inherent := module.CreateInherent(inherentData)

		if inherent.HasValue {
			i++
			extrinsic := types.NewUnsignedUncheckedExtrinsic(inherent.Value)
			result = append(result, extrinsic.Bytes()...)
		}
	}

	if i == 0 {
		return []byte{}
	}

	return append(sc.ToCompact(i).Bytes(), result...)
}

func (re RuntimeExtrinsic) CheckInherents(data primitives.InherentData, block types.Block) primitives.CheckInherentsResult {
	result := primitives.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false
		call := extrinsic.Function

		for _, module := range re.modules {
			if module.IsInherent(call) {
				isInherent = true

				err := module.CheckInherent(extrinsic.Function, data)
				if err != nil {
					e := err.(primitives.IsFatalError)
					err := result.PutError(module.InherentIdentifier(), e)
					// TODO: log depending on error type - handle_put_error_result
					if err != nil {
						log.Critical(err.Error())
					}

					if e.IsFatal() {
						return result
					}
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	return result
}

// EnsureInherentsAreFirst checks if the inherents are before non-inherents.
func (re RuntimeExtrinsic) EnsureInherentsAreFirst(block types.Block) int {
	signedExtrinsicFound := false

	for i, extrinsic := range block.Extrinsics {
		isInherent := false

		if extrinsic.IsSigned() {
			signedExtrinsicFound = true
		} else {
			call := extrinsic.Function

			for _, module := range re.modules {
				if module.IsInherent(call) {
					isInherent = true
				}
			}
		}

		if signedExtrinsicFound && isInherent {
			return i
		}
	}

	return -1
}

func (re RuntimeExtrinsic) OnInitialize(n sc.U32) primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnInitialize(n))
	}

	return weight
}

func (re RuntimeExtrinsic) OnRuntimeUpgrade() primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnRuntimeUpgrade())
	}

	return weight
}

func (re RuntimeExtrinsic) OnFinalize(n sc.U32) {
	for _, m := range re.modules {
		m.OnFinalize(n)
	}
}

func (re RuntimeExtrinsic) OnIdle(n sc.U32, remainingWeight primitives.Weight) primitives.Weight {
	weight := primitives.WeightZero()
	for _, m := range re.modules {
		adjustedRemainingWeight := remainingWeight.SaturatingSub(weight)
		weight = weight.SaturatingAdd(m.OnIdle(n, adjustedRemainingWeight))
	}

	return weight
}

func (re RuntimeExtrinsic) OffchainWorker(n sc.U32) {
	for _, m := range re.modules {
		m.OffchainWorker(n)
	}
}
