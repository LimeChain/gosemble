package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

// It's an account ID (pubkey).
type AccountId struct {
	sc.U64
}

func NewAccountId(value sc.U64) AccountId {
	return AccountId{U64: value}
}

func (a AccountId) Encode(buffer *bytes.Buffer) {
	a.U64.Encode(buffer)
}

func DecodeAccountId(buffer *bytes.Buffer) AccountId {
	return AccountId{sc.DecodeU64(buffer)}
}

func (a MultiAddress) IsAccountId() sc.Bool {
	switch a[0].(type) {
	case AccountId:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountId() AccountId {
	if a.IsAccountId() {
		return a[0].(AccountId)
	} else {
		log.Critical("not a AccountId type")
	}

	panic("unreachable")
}

// It's an account index.
type AccountIndex struct {
	sc.U64
}

func (a AccountIndex) Encode(buffer *bytes.Buffer) {
	a.U64.Encode(buffer)
}

func DecodeAccountIndex(buffer *bytes.Buffer) AccountIndex {
	return AccountIndex{sc.DecodeU64(buffer)}
}

func (a MultiAddress) IsAccountIndex() sc.Bool {
	switch a[0].(type) {
	case AccountIndex:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountIndex() AccountIndex {
	if a.IsAccountIndex() {
		return a[0].(AccountIndex)
	} else {
		log.Critical("not a AccountIndex type")
	}

	panic("unreachable")
}

// It's some arbitrary raw bytes.
type AccountRaw struct {
	sc.Sequence[sc.U8]
}

func (a AccountRaw) Encode(buffer *bytes.Buffer) {
	a.Sequence.Encode(buffer)
}

func DecodeAccountRaw(buffer *bytes.Buffer) AccountRaw {
	return AccountRaw{sc.DecodeSequence[sc.U8](buffer)}
}

func (a MultiAddress) IsRaw() sc.Bool {
	switch a[0].(type) {
	case AccountRaw:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsRaw() AccountRaw {
	if a.IsRaw() {
		return a[0].(AccountRaw)
	} else {
		log.Critical("not an AccountRaw type")
	}

	panic("unreachable")
}

// It's a 32 byte representation.
type Address32 struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewAddress32(values ...sc.U8) Address32 {
	if len(values) != 32 {
		log.Critical("Address32 should be of size 32")
	}
	return Address32{sc.NewFixedSequence(32, values...)}
}

func (a Address32) Encode(buffer *bytes.Buffer) {
	a.FixedSequence.Encode(buffer)
}

func DecodeAddress32(buffer *bytes.Buffer) Address32 {
	return Address32{sc.DecodeFixedSequence[sc.U8](32, buffer)}
}

func (a MultiAddress) IsAddress32() sc.Bool {
	switch a[0].(type) {
	case Address32:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress32() Address32 {
	if a.IsAddress32() {
		return a[0].(Address32)
	} else {
		log.Critical("not a Address32 type")
	}

	panic("unreachable")
}

func (a Address32) Validate() (ok ValidTransaction, err TransactionValidityError) {
	return ok, err
}

func (a Address32) PreDispatch() (ok Pre, err TransactionValidityError) {
	// TODO:
	ok = Pre{}
	return ok, err
}

// Its a 20 byte representation.
type Address20 struct {
	sc.FixedSequence[sc.U8] // size 20
}

func NewAddress20(values ...sc.U8) Address20 {
	if len(values) != 20 {
		log.Critical("Address20 should be of size 20")
	}
	return Address20{sc.NewFixedSequence(20, values...)}
}

func (a Address20) Encode(buffer *bytes.Buffer) {
	a.FixedSequence.Encode(buffer)
}

func DecodeAddress20(buffer *bytes.Buffer) Address20 {
	return Address20{sc.DecodeFixedSequence[sc.U8](20, buffer)}
}

func (a MultiAddress) IsAddress20() sc.Bool {
	switch a[0].(type) {
	case Address20:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress20() Address20 {
	if a.IsAddress20() {
		return a[0].(Address20)
	} else {
		log.Critical("not a Address20 type")
	}

	panic("unreachable")
}

type MultiAddress sc.VaryingData

func NewMultiAddress(value sc.Encodable) MultiAddress {
	switch value.(type) {
	case AccountId, AccountIndex, AccountRaw, Address32, Address20:
		return MultiAddress(sc.NewVaryingData(value))
	default:
		log.Critical("invalid Address type")
	}

	panic("unreachable")
}

func (a MultiAddress) Encode(buffer *bytes.Buffer) {
	if a.IsAccountId() {
		sc.U8(0).Encode(buffer)
		a.AsAccountId().Encode(buffer)
	} else if a.IsAccountIndex() {
		sc.U8(1).Encode(buffer)
		a.AsAccountIndex().Encode(buffer)
	} else if a.IsRaw() {
		sc.U8(2).Encode(buffer)
		a.AsRaw().Encode(buffer)
	} else if a.IsAddress32() {
		sc.U8(3).Encode(buffer)
		a.AsAddress32().Encode(buffer)
	} else if a.IsAddress20() {
		sc.U8(4).Encode(buffer)
		a.AsAddress20().Encode(buffer)
	} else {
		log.Critical("invalid MultiAddress type in Encode")
	}
}

func DecodeMultiAddress(buffer *bytes.Buffer) MultiAddress {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		return MultiAddress{DecodeAccountId(buffer)}
	case 1:
		return MultiAddress{DecodeAccountIndex(buffer)}
	case 2:
		return MultiAddress{DecodeAccountRaw(buffer)}
	case 3:
		return MultiAddress{DecodeAddress32(buffer)}
	case 4:
		return MultiAddress{DecodeAddress20(buffer)}
	default:
		log.Critical("invalid MultiAddress type in Decode")
	}

	panic("unreachable")
}
