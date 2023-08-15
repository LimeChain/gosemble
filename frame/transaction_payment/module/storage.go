package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
)

var defaultMultiplierValue = sc.NewU128FromUint64(1)

type storage struct {
	NextFeeMultiplier *support.StorageValue[sc.U128]
}

func newStorage() *storage {
	return &storage{NextFeeMultiplier: support.NewStorageValueWithDefault(constants.KeyTransactionPayment, constants.KeyNextFeeMultiplier, sc.DecodeU128, &defaultMultiplierValue)}
}
