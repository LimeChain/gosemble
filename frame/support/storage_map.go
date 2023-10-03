package support

import sc "github.com/LimeChain/goscale"

type StorageMap[K, V sc.Encodable] interface {
	Get(k K) V
	Exists(k K) bool
	Put(k K, value V)
	Append(k K, value V)
	TakeBytes(k K) []byte
	Remove(k K)
	Clear(limit sc.U32)
	Mutate(k K, f func(v *V) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
	TryMutateExists(k K, f func(option *sc.Option[V]) sc.Result[sc.Encodable]) sc.Result[sc.Encodable]
}
