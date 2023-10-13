package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ExtrinsicInitializer interface {
	NewChecked(signed sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) primitives.CheckedExtrinsic
}

type extrinsicInitializer struct{}

func NewExtrinsicInitializer() ExtrinsicInitializer {
	return extrinsicInitializer{}
}

func (ex extrinsicInitializer) NewChecked(signer sc.Option[primitives.Address32], function primitives.Call, extra primitives.SignedExtra) primitives.CheckedExtrinsic {
	return checkedExtrinsic{
		signer:        signer,
		function:      function,
		extra:         extra,
		transactional: support.NewTransactional[primitives.PostDispatchInfo, primitives.DispatchError](),
	}
}
