package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TimestampModule struct {
	functions map[sc.U8]support.FunctionMetadata
}

func NewTimestampModule() TimestampModule {
	functions := make(map[sc.U8]support.FunctionMetadata)
	functions[timestamp.FunctionSetIndex] = dispatchables.FnSet{}

	return TimestampModule{
		functions: functions,
	}
}

func (tm TimestampModule) Functions() map[sc.U8]support.FunctionMetadata {
	return tm.functions
}

func (tm TimestampModule) PreDispatch(_ support.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tm TimestampModule) ValidateUnsigned(_ primitives.TransactionSource, _ support.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}
