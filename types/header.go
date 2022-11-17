package types

type BlockNumber uint32

type Header struct {
	ParentHash     Blake2bHash
	Number         BlockNumber
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
