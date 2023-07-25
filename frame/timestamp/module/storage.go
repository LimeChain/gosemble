package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
)

type storage struct {
	Now       *support.StorageValue[sc.U64]
	DidUpdate *support.StorageValue[sc.Bool]
}

func newStorage() *storage {
	return &storage{
		Now:       support.NewStorageValue(constants.KeyTimestamp, constants.KeyNow, sc.DecodeU64),
		DidUpdate: support.NewStorageValue(constants.KeyTimestamp, constants.KeyDidUpdate, sc.DecodeBool),
	}
}
