package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeExtrinsic[N sc.Numeric] struct {
	modules map[sc.U8]types.Module[N]
}

func New[N sc.Numeric](modules map[sc.U8]types.Module[N]) RuntimeExtrinsic[N] {
	return RuntimeExtrinsic[N]{modules: modules}
}

func (re RuntimeExtrinsic[N]) Module(index sc.U8) (module types.Module[N], isFound bool) {
	m, ok := re.modules[index]
	return m, ok
}

func (re RuntimeExtrinsic[N]) CreateInherents(inherentData primitives.InherentData) []byte {
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

func (re RuntimeExtrinsic[N]) CheckInherents(data primitives.InherentData, block types.Block[N]) primitives.CheckInherentsResult {
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
func (re RuntimeExtrinsic[N]) EnsureInherentsAreFirst(block types.Block[N]) int {
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

func (re RuntimeExtrinsic[N]) OnInitialize(n N) primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnInitialize(n))
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OnRuntimeUpgrade() primitives.Weight {
	weight := primitives.Weight{}
	for _, m := range re.modules {
		weight = weight.Add(m.OnRuntimeUpgrade())
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OnFinalize(n N) {
	for _, m := range re.modules {
		m.OnFinalize(n)
	}
}

func (re RuntimeExtrinsic[N]) OnIdle(n N, remainingWeight primitives.Weight) primitives.Weight {
	weight := primitives.WeightZero()
	for _, m := range re.modules {
		adjustedRemainingWeight := remainingWeight.SaturatingSub(weight)
		weight = weight.SaturatingAdd(m.OnIdle(n, adjustedRemainingWeight))
	}

	return weight
}

func (re RuntimeExtrinsic[N]) OffchainWorker(n N) {
	for _, m := range re.modules {
		m.OffchainWorker(n)
	}
}
