package aura

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "AuraApi"
	apiVersion    = 1
)

type Module[N sc.Numeric] struct {
	aura aura.Module[N]
}

func New[N sc.Numeric](aura aura.Module[N]) Module[N] {
	return Module[N]{
		aura: aura,
	}
}

func (m Module[N]) Name() string {
	return ApiModuleName
}

func (m Module[N]) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// Authorities returns current set of AuRa (Authority Round) authorities.
// Returns a pointer-size of the SCALE-encoded set of authorities.
func (m Module[N]) Authorities() int64 {
	authorities := m.aura.Storage.Authorities.GetBytes()

	if !authorities.HasValue {
		return utils.BytesToOffsetAndSize([]byte{0})
	}

	return utils.BytesToOffsetAndSize(sc.SequenceU8ToBytes(authorities.Value))
}

// SlotDuration returns the slot duration for AuRa.
// Returns a pointer-size of the SCALE-encoded slot duration
func (m Module[N]) SlotDuration() int64 {
	slotDuration := m.aura.SlotDuration()
	return utils.BytesToOffsetAndSize(slotDuration.Bytes())
}
