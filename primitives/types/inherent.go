package types

import sc "github.com/LimeChain/goscale"

type ProvideInherent interface {
	CreateInherent(inherent InherentData) sc.Option[Call]
	CheckInherent(call Call, data InherentData) error
	InherentIdentifier() [8]byte
	IsInherent(call Call) bool
}

type DefaultProvideInherent struct {
}

func NewDefaultProvideInherent() DefaultProvideInherent {
	return DefaultProvideInherent{}
}

func (dpi DefaultProvideInherent) CreateInherent(inherent InherentData) sc.Option[Call] {
	return sc.NewOption[Call](nil)
}

func (dpi DefaultProvideInherent) CheckInherent(call Call, data InherentData) error {
	return nil
}

func (dpi DefaultProvideInherent) InherentIdentifier() [8]byte {
	return [8]byte{}
}

func (dpi DefaultProvideInherent) IsInherent(call Call) bool {
	return false
}
