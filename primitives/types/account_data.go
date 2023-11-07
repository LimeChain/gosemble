package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Balance = sc.U128

type AccountData struct {
	Free       Balance
	Reserved   Balance
	MiscFrozen Balance
	FeeFrozen  Balance
}

func (ad AccountData) Encode(buffer *bytes.Buffer) error {
	err := ad.Free.Encode(buffer)
	if err != nil {
		return err
	}
	err = ad.Reserved.Encode(buffer)
	if err != nil {
		return err
	}
	err = ad.MiscFrozen.Encode(buffer)
	if err != nil {
		return err
	}
	return ad.FeeFrozen.Encode(buffer)
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
