package aura

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
)

type storage struct {
	Authorities *support.StorageValue[sc.Sequence[sc.U8]]
	CurrentSlot *support.StorageValue[sc.U64]
}

func newStorage() *storage {
	return &storage{
		Authorities: support.NewStorageValue(constants.KeyAura, constants.KeyAuthorities, sc.DecodeSequence[sc.U8]),
		CurrentSlot: support.NewStorageValue(constants.KeyAura, constants.KeyCurrentSlot, sc.DecodeU64),
	}
}
