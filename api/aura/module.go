package aura

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/aura"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "AuraApi"
	apiVersion    = 1
)

type Module struct {
	aura     aura.AuraModule
	memUtils utils.WasmMemoryTranslator
}

func New(aura aura.AuraModule) Module {
	return Module{
		aura:     aura,
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash, err := hashing.MustBlake2b8([]byte(ApiModuleName))
	if err != nil {
		log.Critical(err.Error())
	}
	return primitives.NewApiItem(hash, apiVersion)
}

// Authorities returns current set of AuRa (Authority Round) authorities.
// Returns a pointer-size of the SCALE-encoded set of authorities.
func (m Module) Authorities() int64 {
	authorities, err := m.aura.GetAuthorities()
	if err != nil {
		log.Critical(err.Error())
	}

	if !authorities.HasValue {
		return m.memUtils.BytesToOffsetAndSize([]byte{0})
	}

	return m.memUtils.BytesToOffsetAndSize(sc.SequenceU8ToBytes(authorities.Value))
}

// SlotDuration returns the slot duration for AuRa.
// Returns a pointer-size of the SCALE-encoded slot duration
func (m Module) SlotDuration() int64 {
	slotDuration := m.aura.SlotDuration()
	return m.memUtils.BytesToOffsetAndSize(slotDuration.Bytes())
}
