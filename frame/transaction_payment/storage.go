package transaction_payment

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
)

var (
	keyTransactionPayment = []byte("TransactionPayment")
	keyNextFeeMultiplier  = []byte("NextFeeMultiplier")
)

var defaultMultiplierValue = sc.NewU128(1)

type Storage interface {
	GetNextFeeMultiplier() sc.U128
}

type storage struct {
	nextFeeMultiplier support.StorageValue[sc.U128]
}

func newStorage() Storage {
	return &storage{
		nextFeeMultiplier: support.NewHashStorageValueWithDefault(
			keyTransactionPayment,
			keyNextFeeMultiplier,
			sc.DecodeU128,
			&defaultMultiplierValue,
		),
	}
}

func (s storage) GetNextFeeMultiplier() sc.U128 {
	return s.nextFeeMultiplier.Get()
}
