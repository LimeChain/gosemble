package inherent

import (
	types2 "github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/timestamp/module"
)

func EnsureInherentsAreFirst(block types2.Block) int {
	signedExtrinsicFound := false

	for i, extrinsic := range block.Extrinsics {
		isInherent := false

		if extrinsic.IsSigned() {
			// Signed extrinsics are not inherents
			isInherent = false
		} else {
			call := extrinsic.Function
			// Iterate through all calls and check if the given call is inherent
			switch call.CallIndex.ModuleIndex {
			case module.Module.Index():
				for _, moduleFn := range module.Module.Functions() {
					if call.CallIndex.FunctionIndex == moduleFn.Index() {
						isInherent = true
					}
				}

			}
		}

		if !isInherent {
			signedExtrinsicFound = true
		}

		if signedExtrinsicFound && isInherent {
			return i
		}
	}

	return -1
}
