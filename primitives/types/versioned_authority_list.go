package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type VersionedAuthorityList struct {
	Version       sc.U8
	AuthorityList sc.Sequence[Authority]
}

func (val VersionedAuthorityList) Encode(buffer *bytes.Buffer) {
	val.Version.Encode(buffer)
	val.AuthorityList.Encode(buffer)
}

func DecodeVersionedAuthorityList(buffer *bytes.Buffer) (VersionedAuthorityList, error) {
	version, err := sc.DecodeU8(buffer)
	if err != nil {
		return VersionedAuthorityList{}, err
	}
	authList, err := sc.DecodeSequenceWith(buffer, DecodeAuthority)
	if err != nil {
		return VersionedAuthorityList{}, err
	}
	return VersionedAuthorityList{
		Version:       version,
		AuthorityList: authList,
	}, nil
}

func (val VersionedAuthorityList) Bytes() []byte {
	return sc.EncodedBytes(val)
}
