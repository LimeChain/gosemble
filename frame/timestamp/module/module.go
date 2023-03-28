package module

import (
	sc "github.com/LimeChain/goscale"
	ts "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/timestamp/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TimestampModule struct {
	functions map[sc.U8]primitives.Call
}

func NewTimestampModule() TimestampModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[ts.FunctionSetIndex] = dispatchables.NewSetCall(nil)

	return TimestampModule{
		functions: functions,
	}
}

func (tm TimestampModule) Functions() map[sc.U8]primitives.Call {
	return tm.functions
}

func (tm TimestampModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tm TimestampModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}
