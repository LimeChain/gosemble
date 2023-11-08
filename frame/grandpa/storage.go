package grandpa

import (
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

var (
	keyGrandpaAuthorities = []byte(":grandpa_authorities")
)

type storage struct {
	Authorities support.StorageValue[primitives.VersionedAuthorityList]
}

func newStorage[S primitives.Signer]() *storage {
	return &storage{Authorities: support.NewSimpleStorageValue(keyGrandpaAuthorities, primitives.DecodeVersionedAuthorityList[S])}
}
