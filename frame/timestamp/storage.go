package timestamp

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
)

var (
	keyTimestamp = []byte("Timestamp")
	keyDidUpdate = []byte("DidUpdate")
	keyNow       = []byte("Now")
)

type storage struct {
	Now       support.StorageValue[sc.U64]
	DidUpdate support.StorageValue[sc.Bool]
}

func newStorage() *storage {
	return &storage{
		Now:       support.NewHashStorageValue(keyTimestamp, keyNow, sc.DecodeU64),
		DidUpdate: support.NewHashStorageValue(keyTimestamp, keyDidUpdate, sc.DecodeBool),
	}
}
