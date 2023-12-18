package types

import sc "github.com/LimeChain/goscale"

// ProvideInherent is an interface, implemented by modules in order to create and validate inherent calls/extrinsics.
type InherentProvider interface {
	// CreateInherent creates an inherent call based on InherentData.
	CreateInherent(inherent InherentData) (sc.Option[Call], error)
	// CheckInherent validates if the provided call is valid and exists in InherentData.
	CheckInherent(call Call, data InherentData) error
	// InherentIdentifier returns the identifier for the specific inherent call. Must be included in InherentData.
	InherentIdentifier() [8]byte
	// IsInherent checks if the call is from the given module.
	IsInherent(call Call) bool
}

// DefaultProvideInherent is an implementation of ProvideInherent and is used by modules, which do not have an
// implementation of ProvideInherent.
type DefaultInherentProvider struct {
}

func NewDefaultProvideInherent() DefaultInherentProvider {
	return DefaultInherentProvider{}
}

func (dp DefaultInherentProvider) CreateInherent(inherent InherentData) (sc.Option[Call], error) {
	return sc.NewOption[Call](nil), nil
}

func (dp DefaultInherentProvider) CheckInherent(call Call, data InherentData) error {
	return nil
}

func (dp DefaultInherentProvider) InherentIdentifier() [8]byte {
	return [8]byte{}
}

func (dp DefaultInherentProvider) IsInherent(call Call) bool {
	return false
}
