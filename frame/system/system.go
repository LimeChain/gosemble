package system

import (
	"math"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/types"
)

func Mutate(who types.Address32, f func(who *types.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	accountInfo := StorageGetAccount(who.FixedSequence)

	result := f(&accountInfo)
	if !result.HasError {
		systemHash := hashing.Twox128(constants.KeySystem)
		accountHash := hashing.Twox128(constants.KeyAccount)

		whoBytes := sc.FixedSequenceU8ToBytes(who.FixedSequence)

		key := append(systemHash, accountHash...)
		key = append(key, hashing.Blake128(whoBytes)...)
		key = append(key, whoBytes...)

		storage.Set(key, accountInfo.Bytes())
	}

	return result
}

func TryMutateExists(who types.Address32, f func(who *types.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := StorageGetAccount(who.FixedSequence)
	wasProviding := false
	if !reflect.DeepEqual(account.Data, types.AccountData{}) {
		wasProviding = true
	}

	someData := &types.AccountData{}
	if wasProviding {
		someData = &account.Data
	}

	result := f(someData)
	if result.HasError {
		return result
	}

	isProviding := !reflect.DeepEqual(someData, types.AccountData{})

	if !wasProviding && isProviding {
		incProviders(who)
	} else if wasProviding && !isProviding {
		status, err := decProviders(who)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}
		if status == types.DecRefStatusExists {
			return result
		}
	} else if !wasProviding && !isProviding {
		return result
	}

	Mutate(who, func(a *types.AccountInfo) sc.Result[sc.Encodable] {
		if someData != nil {
			a.Data = *someData
		} else {
			a.Data = types.AccountData{}
		}

		return sc.Result[sc.Encodable]{}
	})

	return result
}

func AccountTryMutateExists(who types.Address32, f func(who *types.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	account := StorageGetAccount(who.FixedSequence)

	result := f(&account)

	if !result.HasError {
		StorageSetAccount(who.FixedSequence, account)
	}

	return result
}

func incProviders(who types.Address32) types.IncRefStatus {
	result := Mutate(who, func(a *types.AccountInfo) sc.Result[sc.Encodable] {
		if a.Providers == 0 && a.Sufficients == 0 {
			a.Providers = 1
			onCreatedAccount(who)

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.IncRefStatusCreated,
			}
		} else {
			// saturating_add
			newProviders := a.Providers + 1
			if newProviders < a.Providers {
				newProviders = math.MaxUint32
			}

			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.IncRefStatusExisted,
			}
		}
	})

	return result.Value.(types.IncRefStatus)
}

func decProviders(who types.Address32) (types.DecRefStatus, types.DispatchError) {
	result := AccountTryMutateExists(who, func(account *types.AccountInfo) sc.Result[sc.Encodable] {
		if account.Providers == 0 {
			log.Warn("Logic error: Unexpected underflow in reducing provider")

			account.Providers = 1
		}

		if account.Providers == 1 && account.Consumers == 0 && account.Sufficients == 0 {
			return sc.Result[sc.Encodable]{
				HasError: false,
				Value:    types.DecRefStatusReaped,
			}
		}

		if account.Providers == 1 && account.Consumers > 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    types.NewDispatchErrorConsumerRemaining(),
			}
		}

		account.Providers -= 1
		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    types.DecRefStatusExists,
		}
	})

	if result.HasError {
		return sc.U8(0), result.Value.(types.DispatchError)
	}

	return result.Value.(types.DecRefStatus), nil
}
