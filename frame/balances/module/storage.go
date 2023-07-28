package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
)

type storage struct {
	TotalIssuance *support.StorageValue[sc.U128]
	//InactiveIssuance *support.StorageValue[sc.U128]
}

func newStorage() *storage {
	return &storage{
		TotalIssuance: support.NewStorageValue(constants.KeyBalances, constants.KeyTotalIssuance, sc.DecodeU128),
		//InactiveIssuance: support.NewStorageValue(constants.KeyBalances, constants.KeyInactiveIssuance, sc.DecodeU128),
	}
}
