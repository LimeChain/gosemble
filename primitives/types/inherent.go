package types

import sc "github.com/LimeChain/goscale"

// ProvideInherent is an interface, implemented by modules in order to create and validate inherent calls/extrinsics.
type ProvideInherent interface {
	// CreateInherent creates an inherent call based on InherentData.
	CreateInherent(inherent InherentData) sc.Option[Call]
	// CheckInherent validates if the provided call is valid and exists in InherentData.
	CheckInherent(call Call, data InherentData) error
	// InherentIdentifier returns the identifier for the specific inherent call. Must be included in InherentData.
	InherentIdentifier() [8]byte
	// IsInherent checks if the call is from the given module.
	IsInherent(call Call) bool
}

// DefaultProvideInherent is an implementation of ProvideInherent and is used by modules, which do not have an
// implementation of ProvideInherent.
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
