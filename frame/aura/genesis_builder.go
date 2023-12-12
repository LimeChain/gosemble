package aura

import (
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var (
	errAuthoritiesAlreadyInitialized   = errors.New("Authorities are already initialized!") // todo the same with other errors in module.go and return them instead log.Critical
	errAuthoritiesExceedMaxAuthorities = errors.New("Initial authority set must be less than MaxAuthorities")
	errInvalidGenesisConfig            = errors.New("Invalid aura genesis config")
)

type GenesisConfig struct {
	Authorities []string `json:"authorities"`
}

func (m Module) CreateDefaultConfig() ([]byte, error) {
	gc := &GenesisConfig{Authorities: []string{}}
	return json.Marshal(gc)
}

func (m Module) BuildConfig(config []byte) error {
	gc := GenesisConfig{}
	if err := json.Unmarshal(config, &gc); err != nil {
		return err
	}

	if len(gc.Authorities) == 0 {
		return nil
	}

	totalAuthorities, err := m.storage.Authorities.DecodeLen()
	if err != nil {
		return err
	}

	if totalAuthorities.HasValue {
		return errAuthoritiesAlreadyInitialized
	}

	if len(gc.Authorities) > int(m.config.MaxAuthorities) {
		return errAuthoritiesExceedMaxAuthorities
	}

	authorities := sc.Sequence[types.Sr25519PublicKey]{}
	for _, a := range gc.Authorities {
		_, publicKey, err := utils.SS58Decode(a)
		if err != nil {
			return err
		}
		key, err := types.NewSr25519PublicKey(sc.BytesToSequenceU8(publicKey)...)
		if err != nil {
			return err
		}
		authorities = append(authorities, key)
	}
	m.storage.Authorities.Put(authorities)

	return nil
}

func (m Module) ConfigModuleKey() string {
	return "aura"
}
