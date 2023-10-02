package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/types"
)

type Module interface {
	types.InherentProvider
	hooks.DispatchModule
	GetIndex() sc.U8
	Functions() map[sc.U8]types.Call
	PreDispatch(call types.Call) (sc.Empty, types.TransactionValidityError)
	ValidateUnsigned(source types.TransactionSource, call types.Call) (types.ValidTransaction, types.TransactionValidityError)
	Metadata() (sc.Sequence[types.MetadataType], types.MetadataModule)
}
