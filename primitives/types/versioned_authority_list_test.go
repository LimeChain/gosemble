package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	ed25519SignerOne, _ = NewEd25519PublicKey(sc.BytesToFixedSequenceU8(
		common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
	)...)
	ed25519SignerTwo, _ = NewEd25519PublicKey(sc.BytesToFixedSequenceU8(
		common.MustHexToHash("0x88dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff").ToBytes(),
	)...)
	idOne                                = NewAccountId[PublicKey](ed25519SignerOne)
	idTwo                                = NewAccountId[PublicKey](ed25519SignerTwo)
	expectBytesVersionedAuthorityList, _ = hex.DecodeString("010888dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff020000000000000088dc3417d5058ec4b4503e0c12ea1a0a89be200fe98922423d4334014fa6b0ff0300000000000000")
)

var (
	targetVersionedAuthorityList = VersionedAuthorityList{
		Version: sc.U8(1),
		AuthorityList: sc.Sequence[Authority]{
			{
				Id:     idOne,
				Weight: sc.U64(2),
			},
			{
				Id:     idTwo,
				Weight: sc.U64(3),
			},
		},
	}
)

func Test_VersionedAuthorityList_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetVersionedAuthorityList.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectBytesVersionedAuthorityList, buffer.Bytes())
}

func Test_DecodeVersionedAuthorityList(t *testing.T) {
	buffer := bytes.NewBuffer(expectBytesVersionedAuthorityList)

	result, err := DecodeVersionedAuthorityList[Ed25519PublicKey](buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetVersionedAuthorityList, result)
}

func Test_VersionedAuthorityList_Bytes(t *testing.T) {
	assert.Equal(t, expectBytesVersionedAuthorityList, targetVersionedAuthorityList.Bytes())
}
