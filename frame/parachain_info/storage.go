package parachain_info

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
)

var (
	keyAura        = []byte("ParachainInfo")
	keyParachainId = []byte("ParachainId")
)

type storage struct {
	ParachainId support.StorageValue[sc.U32]
}

func newStorage() *storage {
	return &storage{
		ParachainId: support.NewHashStorageValueWithDefault(keyAura, keyParachainId, sc.DecodeU32, &defaultParachainId),
	}
}
