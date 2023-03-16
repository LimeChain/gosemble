package timestamp

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/timestamp"
	dispatchables "github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
)

var Module = TimestampModule{}

type TimestampModule struct {
	Set dispatchables.FnSet
	// TODO: add more dispatchables
}

func (m TimestampModule) Functions() []support.FunctionMetadata {
	return []support.FunctionMetadata{
		m.Set,
		// TODO: add more dispatchables
	}
}

func (m TimestampModule) Index() sc.U8 {
	return timestamp.ModuleIndex
}
