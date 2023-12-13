package system

import (
	"encoding/json"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type GenesisConfig struct{}

func (m module) CreateDefaultConfig() ([]byte, error) {
	gc := struct {
		SystemGc GenesisConfig `json:"system"`
	}{SystemGc: GenesisConfig{}}
	return json.Marshal(gc)
}

func (m module) BuildConfig(_ []byte) error {
	bytes69 := []byte{69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69}
	hash69, err := types.NewBlake2bHash(sc.BytesToFixedSequenceU8(bytes69)...)
	if err != nil {
		return err
	}

	m.StorageBlockHashSet(sc.U64(0), hash69)
	m.storage.ParentHash.Put(hash69)

	m.StorageLastRuntimeUpgradeSet(types.LastRuntimeUpgradeInfo{
		SpecVersion: m.Version().SpecVersion,
		SpecName:    m.Version().SpecName,
	})

	m.storage.ExtrinsicIndex.Put(sc.U32(0))

	// todo missing
	// <UpgradedToU32RefCount<T>>::put(true);
	// <UpgradedToTripleRefCount<T>>::put(true);

	return nil
}
