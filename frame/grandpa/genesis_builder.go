package grandpa

import (
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var (
	errAuthoritiesAlreadyInitialized = errors.New("Authorities are already initialized!") // todo the same with other errors in module.go and return them instead log.Critical
	errInvalidAddrValue              = errors.New("todo invalid address value")
	errInvalidWeightValue            = errors.New("todo invalid weight value")
)

type GenesisConfig struct {
	Authorities [][2]interface{} `json:"authorities"`
}

func (m Module[T]) CreateDefaultConfig() ([]byte, error) {
	gc := &GenesisConfig{Authorities: [][2]interface{}{}}
	return json.Marshal(gc)
}

func (m Module[T]) BuildConfig(config []byte) error {
	gc := GenesisConfig{}
	if err := json.Unmarshal(config, &gc); err != nil {
		return err
	}

	// todo missing
	// CurrentSetId::<T>::put(SetId::default());

	if len(gc.Authorities) == 0 {
		return nil
	}

	totalAuthorities, err := m.storage.Authorities.Get()
	if err != nil {
		return err
	}

	if len(totalAuthorities.AuthorityList) > 0 {
		return errAuthoritiesAlreadyInitialized
	}

	for _, a := range gc.Authorities {
		addrString, ok := a[0].(string)
		if !ok {
			return errInvalidAddrValue
		}

		_, publicKey, err := utils.SS58Decode(addrString)
		if err != nil {
			return err
		}

		ed25519Signer, err := primitives.NewEd25519PublicKey(sc.BytesToSequenceU8(publicKey)...)
		if err != nil {
			return err
		}

		who := primitives.NewAccountId[primitives.PublicKey](ed25519Signer)

		weightFloat, ok := a[1].(float64)
		if !ok {
			return errInvalidWeightValue
		}

		// todo missing
		// &BoundedAuthorityList::<T::MaxAuthorities>::try_from(authorities).expect(
		// 	"Grandpa: `Config::MaxAuthorities` is smaller than the number of genesis authorities!",
		// ),

		weight := sc.U64(uint64(weightFloat))

		totalAuthorities.AuthorityList = append(totalAuthorities.AuthorityList, primitives.Authority{Id: who, Weight: weight})
	}

	// TODO: VersionedList not encoded as expected
	totalAuthorities.Version = AuthorityVersion
	m.storage.Authorities.Put(totalAuthorities)

	// todo missing
	//// NOTE: initialize first session of first set. this is necessary for
	//// the genesis set and session since we only update the set -> session
	//// mapping whenever a new session starts, i.e. through `on_new_session`.
	// SetIdSession::<T>::insert(0, 0);

	return nil
}

func (m Module[T]) ConfigModuleKey() string {
	return "grandpa"
}
