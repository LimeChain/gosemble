package module

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/types"
)

type storage struct {
	Account            *support.StorageMap[types.PublicKey, types.AccountInfo]
	BlockWeight        *support.StorageValue[types.ConsumedWeight]
	BlockHash          *support.StorageMap[sc.U32, types.Blake2bHash]
	BlockNumber        *support.StorageValue[sc.U32]
	AllExtrinsicsLen   *support.StorageValue[sc.U32]
	ExtrinsicData      *support.StorageMap[sc.U32, sc.Sequence[sc.U8]]
	ExtrinsicCount     *support.StorageValue[sc.U32]
	ParentHash         *support.StorageValue[types.Blake2bHash]
	Digest             *support.StorageValue[types.Digest]
	Events             *support.StorageValue[types.EventRecord] // This calls only Append and Kill
	EventCount         *support.StorageValue[sc.U32]
	EventTopics        *support.StorageMap[types.H256, sc.VaryingData]
	LastRuntimeUpgrade *support.StorageValue[types.LastRuntimeUpgradeInfo]
	ExecutionPhase     *support.StorageValue[types.ExtrinsicPhase]
}

func newStorage() *storage {
	return &storage{
		Account:            support.NewStorageMap[types.PublicKey, types.AccountInfo](constants.KeySystem, constants.KeyAccount, hashing.Blake128, types.DecodeAccountInfo),
		BlockWeight:        support.NewStorageValue(constants.KeySystem, constants.KeyBlockWeight, types.DecodeConsumedWeight),
		BlockHash:          support.NewStorageMap[sc.U32, types.Blake2bHash](constants.KeySystem, constants.KeyBlockHash, hashing.Twox64, types.DecodeBlake2bHash),
		BlockNumber:        support.NewStorageValue(constants.KeySystem, constants.KeyNumber, sc.DecodeU32),
		AllExtrinsicsLen:   support.NewStorageValue(constants.KeySystem, constants.KeyAllExtrinsicsLen, sc.DecodeU32),
		ExtrinsicData:      support.NewStorageMap[sc.U32, sc.Sequence[sc.U8]](constants.KeySystem, constants.KeyExtrinsicData, hashing.Twox64, sc.DecodeSequence[sc.U8]),
		ExtrinsicCount:     support.NewStorageValue(constants.KeySystem, constants.KeyExtrinsicCount, sc.DecodeU32),
		ParentHash:         support.NewStorageValue(constants.KeySystem, constants.KeyParentHash, types.DecodeBlake2bHash),
		Digest:             support.NewStorageValue(constants.KeySystem, constants.KeyDigest, types.DecodeDigest),
		Events:             support.NewStorageValue(constants.KeySystem, constants.KeyEvents, types.DecodeEventRecord),
		EventCount:         support.NewStorageValue(constants.KeySystem, constants.KeyEventCount, sc.DecodeU32),
		EventTopics:        support.NewStorageMap[types.H256, sc.VaryingData](constants.KeySystem, constants.KeyEventTopics, hashing.Blake128, func(buffer *bytes.Buffer) sc.VaryingData { return sc.NewVaryingData() }),
		LastRuntimeUpgrade: support.NewStorageValue(constants.KeySystem, constants.KeyLastRuntimeUpgrade, types.DecodeLastRuntimeUpgradeInfo),
		ExecutionPhase:     support.NewStorageValue(constants.KeySystem, constants.KeyExecutionPhase, types.DecodeExtrinsicPhase),
	}
}
