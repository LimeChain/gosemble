package main

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

func Test_Grandpa_Authorities_Empty(t *testing.T) {
	rt, _ := newTestRuntime(t)

	result, err := rt.Exec("GrandpaApi_grandpa_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, []byte{0}, result)
}

func Test_Grandpa_Authorities(t *testing.T) {
	rt, storage := newTestRuntime(t)
	pubKey1 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ee")
	pubKey2 := common.MustHexToBytes("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ef")
	signerOne, e := types.NewEd25519PublicKey(sc.BytesToSequenceU8(pubKey1)...)
	assert.Nil(t, e)
	signerTwo, e := types.NewEd25519PublicKey(sc.BytesToSequenceU8(pubKey2)...)
	assert.Nil(t, e)
	weight := sc.U64(1)

	storageAuthorityList := types.VersionedAuthorityList{
		Version: grandpa.AuthorityVersion,
		AuthorityList: sc.Sequence[types.Authority]{
			{
				Id:     types.New[types.PublicKey](signerOne),
				Weight: weight,
			},
			{
				Id:     types.New[types.PublicKey](signerTwo),
				Weight: weight,
			},
		},
	}
	err := (*storage).Put([]byte(":grandpa_authorities"), storageAuthorityList.Bytes())
	assert.NoError(t, err)

	result, err := rt.Exec("GrandpaApi_grandpa_authorities", []byte{})
	assert.NoError(t, err)

	assert.Equal(t, storageAuthorityList.AuthorityList.Bytes(), result)
}
