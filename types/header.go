package types

type Header struct {
	ParentHash     Blake2bHash
	Number         uint64
	StateRoot      Hash
	ExtrinsicsRoot Hash
	Digest         Digest
}

func (v *Header) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (v *Header) Decode(enc []byte) error {
	return nil
}
