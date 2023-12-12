package balances

import (
	"encoding/json"
	"errors"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var (
	errBalanceBelowExistentialDeposit = errors.New("the balance of any account should always be at least the existential deposit.")
	errDuplicateBalancesInGenesis     = errors.New("duplicate balances in genesis.")
	errInvalidGenesisConfig           = errors.New("Invalid balances genesis config")
	errInvalidBalanceValue            = errors.New("todo invalid balance value")
	errInvalidAddrValue               = errors.New("todo invalid address value")
)

type GenesisConfig struct {
	Balances [][2]interface{} `json:"balances"`
}

func (m Module) CreateDefaultConfig() ([]byte, error) {
	gc := &GenesisConfig{Balances: [][2]interface{}{}}

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

	accExist := make(map[string]bool)
	totalIssuance := sc.NewU128(0)
	for _, b := range gc.Balances {
		balanceFloat, ok := b[1].(float64)
		if !ok {
			return errInvalidBalanceValue
		}

		balance := sc.NewU128(uint64(balanceFloat))

		if balance.Lt(m.Config.ExistentialDeposit) {
			return errBalanceBelowExistentialDeposit
		}

		addrString, ok := b[0].(string)
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

		if accExist[addrString] {
			return errDuplicateBalancesInGenesis
		}

		accExist[addrString] = true
		totalIssuance.Add(balance)

		// TODO: this creates events as well, but should not - call directly the specific functions
		result := m.tryMutateAccount(
			who,
			func(account *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
				return updateAccount(account, balance, sc.NewU128(0))
			},
		)

		if result.HasError {
			return errors.New("todo DispatchError")
			// return result.Value.(primitives.DispatchError) // todo after merge latest develop
		}
	}

	m.storage.TotalIssuance.Put(totalIssuance)

	return nil
}

func (m Module) ConfigModuleKey() string {
	return "balances"
}
