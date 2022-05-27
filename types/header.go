package types

type Header struct {
	ParentHash     Blake2bHash
	Number         uint64
	StateRoot      Hash
	ExtrinsicsRoot Hash
	Digest         Digest
}
