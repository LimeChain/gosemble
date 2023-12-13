package transaction_payment

import (
	"encoding/json"
	sc "github.com/LimeChain/goscale"
)

type GenesisConfig struct {
	Multiplier sc.U128
}

type gcJsonStruct struct {
	TransactionPaymentGc struct {
		Multiplier string `json:"multiplier"`
	} `json:"transactionPayment"`
}

func (gc *GenesisConfig) UnmarshalJSON(data []byte) error {
	gcJson := gcJsonStruct{}

	if err := json.Unmarshal(data, &gcJson); err != nil {
		return err
	}

	multiplier, err := sc.NewU128FromString(gcJson.TransactionPaymentGc.Multiplier)
	if err != nil {
		return err
	}

	gc.Multiplier = multiplier
	return nil
}

func (m module) CreateDefaultConfig() ([]byte, error) {
	gc := &gcJsonStruct{}
	gc.TransactionPaymentGc.Multiplier = defaultMultiplierValue.ToBigInt().String()

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
