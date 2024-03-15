package support

import sc "github.com/LimeChain/goscale"

type StorageRawValue interface {
	Get() (sc.Sequence[sc.U8], error)
	Put(value sc.Sequence[sc.U8])
	Clear()
	ClearPrefix(limit sc.U32)
}
