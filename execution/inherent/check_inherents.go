package inherent

import (
	cts "github.com/LimeChain/gosemble/constants/timestamp"
	types2 "github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/LimeChain/gosemble/primitives/types"
)

func CheckInherents(data types.InherentData, block types2.Block) types.CheckInherentsResult {
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
		case module.Module.Index():
			for _, moduleFn := range module.Module.Functions() {
				if call.CallIndex.FunctionIndex == moduleFn.Index() {
					isInherent = true
					err := timestamp.CheckInherent(call.Args, data)
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
