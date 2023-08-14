package grandpa

import (
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type storage struct {
	Authorities *support.SimpleStorageValue[primitives.VersionedAuthorityList]
}

func newStorage() *storage {
	return &storage{Authorities: support.NewSimpleStorageValue(constants.KeyGrandpaAuthorities, primitives.DecodeVersionedAuthorityList)}
}
