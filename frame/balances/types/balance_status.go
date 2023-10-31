package types

import (
	"bytes"

	"github.com/LimeChain/gosemble/primitives/log"

	sc "github.com/LimeChain/goscale"
)

const (
	BalanceStatusFree sc.U8 = iota
	BalanceStatusReserved
)

type BalanceStatus = sc.U8

func DecodeBalanceStatus(buffer *bytes.Buffer) (sc.U8, error) {
	value, err := sc.DecodeU8(buffer)
	if err != nil {
		return sc.U8(0), err
	}
	switch value {
	case BalanceStatusFree, BalanceStatusReserved:
		return value, nil
	default:
		log.Critical("invalid balance status type")
	}

	panic("unreachable")
}
