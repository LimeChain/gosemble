package aura

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
)

var (
	keyAura        = []byte("Aura")
	keyAuthorities = []byte("Authorities")
	keyCurrentSlot = []byte("CurrentSlot")
)

type storage struct {
	Authorities support.StorageValue[sc.Sequence[sc.U8]]
	CurrentSlot support.StorageValue[sc.U64]
}

func newStorage() *storage {
	return &storage{
		Authorities: support.NewHashStorageValue(keyAura, keyAuthorities, sc.DecodeSequence[sc.U8]),
		CurrentSlot: support.NewHashStorageValue(keyAura, keyCurrentSlot, sc.DecodeU64),
	}
}
