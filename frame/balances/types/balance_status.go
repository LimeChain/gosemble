package types

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
)

const (
	BalanceStatusFree sc.U8 = iota
	BalanceStatusReserved
)

type BalanceStatus = sc.U8

var (
	errInvalidBalanceStatusType = errors.New("invalid balance status type")
)

func DecodeBalanceStatus(buffer *bytes.Buffer) (sc.U8, error) {
	value, err := sc.DecodeU8(buffer)
	if err != nil {
		return sc.U8(0), err
	}
	switch value {
	case BalanceStatusFree, BalanceStatusReserved:
		return value, nil
	default:
		return sc.U8(0), errInvalidBalanceStatusType
	}
}
