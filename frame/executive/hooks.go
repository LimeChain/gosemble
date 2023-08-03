package executive

import (
	"strconv"

	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

func onRuntimeUpgrade() types.Weight {
	return types.WeightFromParts(200, 0)
}

func onIdle(n types.BlockNumber, remainingWeight types.Weight) types.Weight {
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	// log.Trace(fmt.Sprintf("on_idle %v, %v)", n, remainingWeight))
	log.Trace("on_idle " + strconv.Itoa(int(n)) + " " + strconv.Itoa(int(remainingWeight.RefTime)))
	return types.WeightFromParts(175, 0)
}
