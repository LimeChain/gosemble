package aura

import (
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/itering/subscan/util/ss58"
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

	for _, a := range gc.Authorities {
		m.storage.Authorities.Append(sc.BytesToSequenceU8([]byte(ss58.Decode(a, 42)))) // todo ensure handled properly
	}

	return nil
}

func (m Module) ConfigModuleKey() string {
	return "aura"
}
