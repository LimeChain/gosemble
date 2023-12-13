package main

import (
	"bytes"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/pkg/scale"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_CreateDefaultConfig(t *testing.T) {
	rt, _ := newTestRuntime(t)
	wantGc := []byte("{\"system\":{},\"aura\":{\"authorities\":[]},\"grandpa\":{\"authorities\":[]},\"balances\":{\"balances\":[]},\"transactionPayment\":{\"multiplier\":\"1\"}}")

	res, err := rt.Exec("GenesisBuilder_create_default_config", []byte{})
	assert.NoError(t, err)

	resDecoded, err := sc.DecodeSequence[sc.U8](bytes.NewBuffer(res))
	assert.Equal(t, wantGc, sc.SequenceU8ToBytes(resDecoded))
}

func Test_BuildConfig(t *testing.T) {
	rt, storage := newTestRuntime(t)

	gc := []byte("{\"system\":{},\"aura\":{\"authorities\":[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\"]},\"grandpa\":{\"authorities\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1]]},\"balances\":{\"balances\":[[\"5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY\",1000000000000000000]]},\"transactionPayment\":{\"multiplier\":\"2\"}}")

	_, err := rt.Exec("GenesisBuilder_build_config", sc.BytesToSequenceU8(gc).Bytes())
	assert.NoError(t, err)

	// assert BlockHash
	encBlockNumber, _ := scale.Marshal(uint64(0))
	blockNumHash, _ := common.Twox64(encBlockNumber)
	blockHashKey := append(keySystemHash, keyBlockHash...)
	blockHashKey = append(blockHashKey, blockNumHash...)
	blockHashKey = append(blockHashKey, encBlockNumber...)
	zeroBlockHash := (*storage).Get(blockHashKey)
	bytes69 := []byte{69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69, 69}
	wantBlockHash, _ := types.NewBlake2bHash(sc.BytesToFixedSequenceU8(bytes69)...)
	assert.Equal(t, wantBlockHash.Bytes(), zeroBlockHash)

	// assert ParentHash
	parentHash := (*storage).Get(append(keySystemHash, keyParentHash...))
	assert.Equal(t, wantBlockHash.Bytes(), parentHash)

	// assert LastRuntimeUpgradeSet
	lrui := (*storage).Get(append(keySystemHash, keyLastRuntime...))
	wantLrui := types.LastRuntimeUpgradeInfo{SpecVersion: 100, SpecName: "node-template"}
	assert.Equal(t, wantLrui.Bytes(), lrui)

	// assert ExtrinsicIndex
	extrinsicIndex := (*storage).Get(keyExtrinsicIndex)
	wantExtrinsicIndex := sc.U32(0)
	assert.Equal(t, wantExtrinsicIndex.Bytes(), extrinsicIndex)

	// assert aura authorities
	auraAuthorities := (*storage).Get(append(keyAuraHash, keyAuthoritiesHash...))
	wantPubKey := sc.BytesToSequenceU8([]byte{212, 53, 147, 199, 21, 253, 211, 28, 97, 20, 26, 189, 4, 169, 159, 214, 130, 44, 133, 88, 133, 76, 205, 227, 154, 86, 132, 231, 165, 109, 162, 125})
	wantAuraAuthorityPubKey, _ := types.NewSr25519PublicKey(wantPubKey...)
	wantAuraAuthorities := sc.Sequence[types.Sr25519PublicKey]{wantAuraAuthorityPubKey}
	assert.Equal(t, wantAuraAuthorities.Bytes(), auraAuthorities)

	// assert grandpa authorities
	grandpaAuthorities := (*storage).Get(keyGrandpaAuthorities)
	wantGrandpaAuthorityPubKey, _ := types.NewEd25519PublicKey(wantPubKey...)
	accId := types.NewAccountId[types.PublicKey](wantGrandpaAuthorityPubKey)
	authorities := sc.Sequence[types.Authority]{{Id: accId, Weight: sc.U64(1)}}
	wantGrandpaAuthorities := types.VersionedAuthorityList{AuthorityList: authorities, Version: grandpa.AuthorityVersion}
	assert.Equal(t, wantGrandpaAuthorities.Bytes(), grandpaAuthorities)

	// assert balance
	accHash, _ := common.Blake2b128(accId.Bytes())
	keyStorageAccount := append(keySystemHash, keyAccountHash...)
	keyStorageAccount = append(keyStorageAccount, accHash...)
	keyStorageAccount = append(keyStorageAccount, accId.Bytes()...)
	accInfo := (*storage).Get(keyStorageAccount)
	wantBalance := sc.NewU128(uint64(1000000000000000000))
	wantAccInfo := types.AccountInfo{Data: types.AccountData{Free: wantBalance}}
	assert.Equal(t, wantAccInfo.Bytes(), accInfo)

	// assert total issuance
	totalIssuance := (*storage).Get(append(keyBalancesHash, keyTotalIssuanceHash...))
	assert.Equal(t, wantBalance.Bytes(), totalIssuance)

	// assert next fee multiplier
	nextFeeMultiplier := (*storage).Get(append(keyTransactionPaymentHash, keyNextFeeMultiplierHash...))
	wantNextFeeMultiplier := sc.NewU128(2)
	assert.Equal(t, wantNextFeeMultiplier.Bytes(), nextFeeMultiplier)
}