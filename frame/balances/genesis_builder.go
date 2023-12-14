package balances

import (
	"bytes"
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/vedhavyas/go-subkey"
)

var (
	errBalanceBelowExistentialDeposit = errors.New("the balance of any account should always be at least the existential deposit.")
	errDuplicateBalancesInGenesis     = errors.New("duplicate balances in genesis.")
	errInvalidBalanceValue            = errors.New("invalid balance in genesis config json")
	errInvalidAddrValue               = errors.New("invalid address in genesis config json")
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

	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	jsonDecoder.UseNumber()
	if err := jsonDecoder.Decode(&gcJson); err != nil {
		return err
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

		accId, err := types.NewAccountId(sc.BytesToSequenceU8(publicKey)...)
		if err != nil {
			return err
		}

		balance, ok := b[1].(json.Number)
		if !ok {
			return errInvalidBalanceValue
		}

		balanceU128, err := sc.NewU128FromString(balance.String())
		if err != nil {
			return err
		}

		gc.Balances = append(gc.Balances, gcAccountBalance{AccountId: accId, Balance: balanceU128})
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

		result, err := m.Config.StoredMap.TryMutateExists(
			b.AccountId,
			func(maybeAccount *types.AccountData) sc.Result[sc.Encodable] {
				return updateAccount(maybeAccount, b.Balance, sc.NewU128(0))
			},
		)
		if err != nil {
			return err
		}
		if result.HasError {
			err, _ := result.Value.(error)
			return err
		}
	}

	m.storage.TotalIssuance.Put(totalIssuance)

	return nil
}
