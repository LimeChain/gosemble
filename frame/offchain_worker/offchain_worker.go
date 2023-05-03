package offchain_worker

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

func OffchainWorker(dataPtr int32, dataLen int32) {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	header := types.DecodeHeader(buffer)

	system.Initialize(header.Number, header.ParentHash, header.Digest)

	hash := hashing.Blake256(header.Bytes())

	system.StorageSetBlockHash(header.Number, types.NewBlake2bHash(sc.BytesToSequenceU8(hash)...))

	// TODO:
	/*
		<AllPalletsWithSystem as OffchainWorker<System::BlockNumber>>::offchain_worker(*header.number(),)
	*/
}
