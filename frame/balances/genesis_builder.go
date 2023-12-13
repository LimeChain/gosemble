package balances

import (
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/vedhavyas/go-subkey"
)

var (
	errBalanceBelowExistentialDeposit = errors.New("the balance of any account should always be at least the existential deposit.")
	errDuplicateBalancesInGenesis     = errors.New("duplicate balances in genesis.")
	errInvalidGenesisConfig           = errors.New("Invalid balances genesis config")
	errInvalidBalanceValue            = errors.New("todo invalid balance value")
	errInvalidAddrValue               = errors.New("todo invalid address value")
)

type gcAccountBalance struct {
	AccountId types.AccountId
	Balance   types.Balance
}
type GenesisConfig struct {
	Balances []gcAccountBalance
}

type gcJsonStruct struct {
	BalancesGc struct {
		Balances [][2]interface{} `json:"balances"`
	} `json:"balances"`
}

func (gc *GenesisConfig) UnmarshalJSON(data []byte) error {
	gcJson := gcJsonStruct{}

	if err := json.Unmarshal(data, &gcJson); err != nil {
		return err
	}

	if len(gcJson.BalancesGc.Balances) == 0 {
		return nil
	}
	addrExists := map[string]bool{}
	for _, b := range gcJson.BalancesGc.Balances {
		addrString, ok := b[0].(string)
		if !ok {
			return errInvalidAddrValue
		}

		if addrExists[addrString] {
			return errDuplicateBalancesInGenesis
		}

		_, publicKey, err := subkey.SS58Decode(addrString)
		if err != nil {
			return err
		}

		// ed25519Signer, err := types.NewEd25519PublicKey(sc.BytesToSequenceU8(publicKey)...)
		// if err != nil {
		// 	return err
		// }

		accId, err := types.NewAccountId(sc.BytesToSequenceU8(publicKey)...)
		if err != nil {
			return err
		}
		balanceFloat, ok := b[1].(float64)
		if !ok {
			return errInvalidBalanceValue
		}

		gc.Balances = append(gc.Balances, gcAccountBalance{AccountId: accId, Balance: sc.NewU128(uint64(balanceFloat))})
		addrExists[addrString] = true
	}

	return nil
}
func (m Module) CreateDefaultConfig() ([]byte, error) {
	gc := &gcJsonStruct{}
	gc.BalancesGc.Balances = [][2]interface{}{}

	return json.Marshal(gc)
}

func (m Module) BuildConfig(config []byte) error {
	gc := GenesisConfig{}
	if err := json.Unmarshal(config, &gc); err != nil {
		return err
	}

	if len(gc.Balances) == 0 {
		return nil
	}

	totalIssuance := sc.NewU128(0)
	for _, b := range gc.Balances {
		if b.Balance.Lt(m.Config.ExistentialDeposit) {
			return errBalanceBelowExistentialDeposit
		}

		totalIssuance = totalIssuance.Add(b.Balance)

		if _, err := m.Config.StoredMap.IncProviders(b.AccountId); err != nil {
			return err
		}

		m.Config.StoredMap.Put(b.AccountId, types.AccountInfo{Data: types.AccountData{Free: b.Balance}})
	}

	m.storage.TotalIssuance.Put(totalIssuance)

	return nil
}
