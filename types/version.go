package types

import "github.com/ChainSafe/gossamer/lib/scale"

type ApiItem struct {
	Name    [8]byte
	Version uint32
}

type VersionData struct {
	SpecName           []byte
	ImplName           []byte
	AuthoringVersion   uint32
	SpecVersion        uint32
	ImplVersion        uint32
	Apis               []ApiItem
	TransactionVersion uint32
	StateVersion       uint32
}

func (v *VersionData) Encode() ([]byte, error) {
	enc, err := scale.Encode(v)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func (v *VersionData) Decode(enc []byte) error {
	var data VersionData

	_, err := scale.Decode(enc, &data)
	if err != nil {
		return err
	}

	v.SpecName = data.SpecName
	v.ImplName = data.ImplName
	v.AuthoringVersion = data.AuthoringVersion
	v.SpecVersion = data.SpecVersion
	v.ImplVersion = data.ImplVersion
	v.Apis = data.Apis
	v.TransactionVersion = data.TransactionVersion
	v.StateVersion = data.StateVersion

	return nil
}
