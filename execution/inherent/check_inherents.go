package inherent

import (
	cts "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/types"
)

func CheckInherents(data types.InherentData, block types.Block) types.CheckInherentsResult {
	result := types.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false

		call := extrinsic.Function

		switch call.CallIndex.ModuleIndex {
		case timestamp.Module.Index:
			for funcKey := range timestamp.Module.Functions {
				if call.CallIndex.FunctionIndex == timestamp.Module.Functions[funcKey].Index {
					isInherent = true
					err := timestamp.CheckInherent(call, data)
					if err != nil {
						err := result.PutError(cts.InherentIdentifier, err.(types.IsFatalError))
						if err != nil {
							panic(err)
						}

						if result.FatalError {
							return result
						}
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

	// TODO: go through all required pallets with required inherents

	return result
}
