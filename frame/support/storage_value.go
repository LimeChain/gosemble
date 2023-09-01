package support

import sc "github.com/LimeChain/goscale"

type StorageValue[T sc.Encodable] interface {
	Get() T
	GetBytes() sc.Option[sc.Sequence[sc.U8]]
	Exists() bool
	Put(value T)
	Clear()
	Append(value T)
	TakeBytes() []byte
	Take() T
	DecodeLen() sc.Option[sc.U64]
}
