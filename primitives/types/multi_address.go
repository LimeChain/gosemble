package types

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
)

// AccountId It's an account ID (pubkey).
type AccountId struct {
	Address32 // TODO: Varies depending on Signature (32 for ed25519 and sr25519, 33 for ecdsa)
}

func DecodeAccountId(buffer *bytes.Buffer) (AccountId, error) {
	addr32, err := DecodeAddress32(buffer)
	if err != nil {
		return AccountId{}, err
	}
	return AccountId{addr32}, nil // TODO: length 32 or 33 depending on algorithm
}

// AccountIndex It's an account index.
type AccountIndex = sc.U32

// AccountRaw It's some arbitrary raw bytes.
type AccountRaw struct {
	sc.Sequence[sc.U8]
}

func (a AccountRaw) Encode(buffer *bytes.Buffer) {
	a.Sequence.Encode(buffer)
}

func DecodeAccountRaw(buffer *bytes.Buffer) (AccountRaw, error) {
	seq, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return AccountRaw{}, err
	}
	return AccountRaw{seq}, nil
}

// Address32 It's a 32 byte representation.
type Address32 struct {
	sc.FixedSequence[sc.U8] // size 32
}

func NewAddress32(values ...sc.U8) (Address32, error) {
	if len(values) != 32 {
		return Address32{}, errors.New("Address32 should be of size 32")
	}
	return Address32{sc.NewFixedSequence(32, values...)}, nil
}

func DecodeAddress32(buffer *bytes.Buffer) (Address32, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](32, buffer)
	if err != nil {
		return Address32{}, err
	}
	return Address32{seq}, nil
}

// Address20 It's a 20 byte representation.
type Address20 struct {
	sc.FixedSequence[sc.U8] // size 20
}

func NewAddress20(values ...sc.U8) (Address20, error) {
	if len(values) != 20 {
		return Address20{}, errors.New("Address20 should be of size 20")
	}
	return Address20{sc.NewFixedSequence(20, values...)}, nil
}

func DecodeAddress20(buffer *bytes.Buffer) (Address20, error) {
	seq, err := sc.DecodeFixedSequence[sc.U8](20, buffer)
	if err != nil {
		return Address20{}, err
	}
	return Address20{seq}, nil
}

const (
	MultiAddressId sc.U8 = iota
	MultiAddressIndex
	MultiAddressRaw
	MultiAddress32
	MultiAddress20
)

type MultiAddress struct {
	sc.VaryingData
}

func NewMultiAddressId(id AccountId) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressId, id)}
}

func NewMultiAddressIndex(index AccountIndex) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressIndex, sc.ToCompact(index))}
}

func NewMultiAddressRaw(accountRaw AccountRaw) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddressRaw, accountRaw)}
}

func NewMultiAddress32(address Address32) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddress32, address)}
}

func NewMultiAddress20(address Address20) MultiAddress {
	return MultiAddress{sc.NewVaryingData(MultiAddress20, address)}
}

func DecodeMultiAddress(buffer *bytes.Buffer) (MultiAddress, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MultiAddress{}, err
	}

	switch b {
	case MultiAddressId:
		accId, err := DecodeAccountId(buffer)
		if err != nil {
			return MultiAddress{}, err
		}
		return NewMultiAddressId(accId), nil
	case MultiAddressIndex:
		compact, err := sc.DecodeCompact(buffer)
		if err != nil {
			return MultiAddress{}, err
		}
		index := sc.U32(compact.ToBigInt().Uint64())
		return NewMultiAddressIndex(index), nil
	case MultiAddressRaw:
		accRaw, err := DecodeAccountRaw(buffer)
		if err != nil {
			return MultiAddress{}, err
		}
		return NewMultiAddressRaw(accRaw), nil
	case MultiAddress32:
		addr32, err := DecodeAddress32(buffer)
		if err != nil {
			return MultiAddress{}, err
		}
		return NewMultiAddress32(addr32), nil
	case MultiAddress20:
		addr20, err := DecodeAddress20(buffer)
		if err != nil {
			return MultiAddress{}, err
		}
		return NewMultiAddress20(addr20), nil
	default:
		return MultiAddress{}, newTypeError("MultiAddress")
	}
}

func (a MultiAddress) IsAccountId() bool {
	switch a.VaryingData[0] {
	case MultiAddressId:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountId() (AccountId, error) {
	if a.IsAccountId() {
		return a.VaryingData[1].(AccountId), nil
	} else {
		return AccountId{}, newTypeError("AccountId")
	}
}

func (a MultiAddress) IsAccountIndex() bool {
	switch a.VaryingData[0] {
	case MultiAddressIndex:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAccountIndex() (AccountIndex, error) {
	if a.IsAccountIndex() {
		compact := a.VaryingData[1].(sc.Compact)
		return sc.U32(compact.ToBigInt().Uint64()), nil
	} else {
		return 0, newTypeError("AccountIndex")
	}
}

func (a MultiAddress) IsRaw() bool {
	switch a.VaryingData[0] {
	case MultiAddressRaw:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsRaw() (AccountRaw, error) {
	if a.IsRaw() {
		return a.VaryingData[1].(AccountRaw), nil
	} else {
		return AccountRaw{}, newTypeError("AccountRaw")
	}
}

func (a MultiAddress) IsAddress32() bool {
	switch a.VaryingData[0] {
	case MultiAddress32:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress32() (Address32, error) {
	if a.IsAddress32() {
		return a.VaryingData[1].(Address32), nil
	} else {
		return Address32{}, newTypeError("Address32")
	}
}

func (a MultiAddress) IsAddress20() bool {
	switch a.VaryingData[0] {
	case MultiAddress20:
		return true
	default:
		return false
	}
}

func (a MultiAddress) AsAddress20() (Address20, error) {
	if a.IsAddress20() {
		return a.VaryingData[1].(Address20), nil
	} else {
		return Address20{}, newTypeError("Address20")
	}
}
