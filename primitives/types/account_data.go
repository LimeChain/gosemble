package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
)

type Balance = sc.U128

type AccountData struct {
	Free       Balance
	Reserved   Balance
	MiscFrozen Balance
	FeeFrozen  Balance
}

func (ad AccountData) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		ad.Free,
		ad.Reserved,
		ad.MiscFrozen,
		ad.FeeFrozen,
	)
}

func (ad AccountData) Bytes() []byte {
	return sc.EncodedBytes(ad)
}

func DecodeAccountData(buffer *bytes.Buffer) (AccountData, error) {
	free, err := sc.DecodeU128(buffer)
	if err != nil {
		return AccountData{}, err
	}
	reserved, err := sc.DecodeU128(buffer)
	if err != nil {
		return AccountData{}, err
	}
	misc, err := sc.DecodeU128(buffer)
	if err != nil {
		return AccountData{}, err
	}
	fee, err := sc.DecodeU128(buffer)
	if err != nil {
		return AccountData{}, err
	}
	return AccountData{
		Free:       free,
		Reserved:   reserved,
		MiscFrozen: misc,
		FeeFrozen:  fee,
	}, nil
}

func (ad AccountData) Total() sc.U128 {
	return ad.Free.Add(ad.Reserved)
}
