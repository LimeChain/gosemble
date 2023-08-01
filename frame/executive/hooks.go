package executive

import (
	"fmt"

	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

func onRuntimeUpgrade() types.Weight {
	return types.WeightFromParts(200, 0)
}

func onIdle(n types.BlockNumber, remainingWeight types.Weight) types.Weight {
	log.Trace(fmt.Sprintf("on_idle %v, %v)", n, remainingWeight))
	return types.WeightFromParts(175, 0)
}
