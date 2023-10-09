package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ExtrinsicInitializer interface {
	NewChecked(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) types.CheckedExtrinsic
}

type extrinsicInitializer struct{}

func NewExtrinsicInitializer() ExtrinsicInitializer {
	return extrinsicInitializer{}
}

func (ex extrinsicInitializer) NewChecked(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) types.CheckedExtrinsic {
	return checkedExtrinsic{
		signed:        signed,
		function:      function,
		extra:         extra,
		transactional: support.NewTransactional[primitives.PostDispatchInfo, primitives.DispatchError](),
	}
}
