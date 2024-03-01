package parachain_info

import (
	"encoding/json"
	sc "github.com/LimeChain/goscale"
)

type GenesisConfig struct {
	ParachainId sc.U32
}

type genesisConfigJsonStruct struct {
	ParachainInfoGenesisConfig struct {
		ParachainId uint32 `json:"parachainId"`
	} `json:"parachainInfo"`
}

func (gc *GenesisConfig) UnmarshalJSON(data []byte) error {
	gcJson := genesisConfigJsonStruct{}

	if err := json.Unmarshal(data, &gcJson); err != nil {
		return err
	}

	gc.ParachainId = sc.U32(gcJson.ParachainInfoGenesisConfig.ParachainId)

	return nil
}

func (m Module) CreateDefaultConfig() ([]byte, error) {
	gc := genesisConfigJsonStruct{}

	gc.ParachainInfoGenesisConfig.ParachainId = uint32(defaultParachainId)

	return json.Marshal(gc)
}

func (m Module) BuildConfig(config []byte) error {
	gc := GenesisConfig{}
	if err := json.Unmarshal(config, &gc); err != nil {
		return err
	}

	m.storage.ParachainId.Put(gc.ParachainId)

	return nil
}
