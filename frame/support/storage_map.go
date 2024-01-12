package support

import sc "github.com/LimeChain/goscale"

type StorageMap[K, V sc.Encodable] interface {
	Get(k K) (V, error)
	Exists(k K) bool
	Put(k K, value V)
	Append(k K, value V)
	TakeBytes(k K) ([]byte, error)
	Remove(k K)
	Clear(limit sc.U32)
	Mutate(k K, f func(v *V) (sc.Encodable, error)) (sc.Encodable, error)
	TryMutateExists(k K, f func(option *sc.Option[V]) (sc.Encodable, error)) (sc.Encodable, error)
}
