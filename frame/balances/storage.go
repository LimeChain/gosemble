package balances

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
)

var (
	keyBalances      = []byte("Balances")
	keyTotalIssuance = []byte("TotalIssuance")
)

type storage struct {
	TotalIssuance support.StorageValue[sc.U128]
}

func newStorage() *storage {
	return &storage{
		TotalIssuance: support.NewHashStorageValue(keyBalances, keyTotalIssuance, sc.DecodeU128),
	}
}
