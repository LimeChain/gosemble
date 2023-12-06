package transaction_payment

import (
	"encoding/json"
	sc "github.com/LimeChain/goscale"
)

type GenesisConfig struct {
	Multiplier uint64 `json:"multiplier,string"` // todo test if can unmarshal from string
}

func (m module) CreateDefaultConfig() ([]byte, error) {
	gc := &GenesisConfig{Multiplier: defaultMultiplierValue.ToBigInt().Uint64()}
	return json.Marshal(gc)
}

func (m module) BuildConfig(config []byte) error {
	gc := GenesisConfig{}
	if err := json.Unmarshal(config, &gc); err != nil {
		return err
	}

	// todo missing
	// StorageVersion::<T>::put(Releases::V2);

	m.storage.NextFeeMultiplier.Put(sc.NewU128(gc.Multiplier))

	return nil
}

func (m module) ConfigModuleKey() string {
	return "transactionPayment"
}
