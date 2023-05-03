package grandpa

import (
	"fmt"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/grandpa"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

func Authorities() int64 {
	versionedAuthorityList := storage.GetDecode(constants.KeyGrandpaAuthorities, types.DecodeVersionedAuthorityList)

	authorities := versionedAuthorityList.AuthorityList
	if versionedAuthorityList.Version != grandpa.AuthorityVersion {
		log.Warn(fmt.Sprintf("unknown Grandpa authorities version: [%d]", versionedAuthorityList.Version))
		authorities = sc.Sequence[types.Authority]{}
	}

	return utils.BytesToOffsetAndSize(authorities.Bytes())
}
