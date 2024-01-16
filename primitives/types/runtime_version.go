package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type RuntimeVersion struct {
	SpecName           sc.Str
	ImplName           sc.Str
	AuthoringVersion   sc.U32
	SpecVersion        sc.U32
	ImplVersion        sc.U32
	Apis               sc.Sequence[ApiItem]
	TransactionVersion sc.U32
	StateVersion       sc.U8
}

func (rv RuntimeVersion) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		rv.SpecName,
		rv.ImplName,
		rv.AuthoringVersion,
		rv.SpecVersion,
		rv.ImplVersion,
		rv.Apis,
		rv.TransactionVersion,
		rv.StateVersion,
	)
}

func (rv *RuntimeVersion) SetApis(apis sc.Sequence[ApiItem]) {
	rv.Apis = apis
}

func DecodeRuntimeVersion(buffer *bytes.Buffer) (RuntimeVersion, error) {
	var rv RuntimeVersion

	specName, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}
	implName, err := sc.DecodeStr(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}
	authVersion, err := sc.DecodeU32(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}
	specVersion, err := sc.DecodeU32(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}
	implVersion, err := sc.DecodeU32(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}

	rv.SpecName = specName
	rv.ImplName = implName
	rv.AuthoringVersion = authVersion
	rv.SpecVersion = specVersion
	rv.ImplVersion = implVersion

	compact, err := sc.DecodeCompact[sc.U64](buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}

	apisLength := compact.ToBigInt().Uint64()

	if apisLength != 0 {
		var apis []ApiItem
		for i := 0; i < int(apisLength); i++ {
			apiItem, err := DecodeApiItem(buffer)
			if err != nil {
				return RuntimeVersion{}, err
			}
			apis = append(apis, apiItem)
		}
		rv.Apis = apis
	}
	transactionVersion, err := sc.DecodeU32(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}
	stateVersion, err := sc.DecodeU8(buffer)
	if err != nil {
		return RuntimeVersion{}, err
	}

	rv.TransactionVersion = transactionVersion
	rv.StateVersion = stateVersion

	return rv, nil
}

func (rv RuntimeVersion) Bytes() []byte {
	return sc.EncodedBytes(rv)
}
