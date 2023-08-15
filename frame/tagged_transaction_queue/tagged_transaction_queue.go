package tagged_transaction_queue

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "TaggedTransactionQueue"
	apiVersion    = 3
)

type TaggedTransactionQueue interface {
	ValidateTransaction(dataPtr int32, dataLen int32) int64
}

type Module struct {
	executive executive.Module
	decoder   types.ModuleDecoder
}

func New(executive executive.Module, decoder types.ModuleDecoder) Module {
	return Module{
		executive: executive,
		decoder:   decoder,
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// ValidateTransaction validates an extrinsic at a given block.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded tx source, extrinsic and block hash.
// Returns a pointer-size of the SCALE-encoded result whether the extrinsic is valid.
// [Specification](https://spec.polkadot.network/#sect-rte-validate-transaction)
func (m Module) ValidateTransaction(dataPtr int32, dataLen int32) int64 {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	txSource := primitives.DecodeTransactionSource(buffer)
	tx := m.decoder.DecodeUncheckedExtrinsic(buffer)
	blockHash := primitives.DecodeBlake2bHash(buffer)

	ok, err := m.executive.ValidateTransaction(txSource, tx, blockHash)

	var res primitives.TransactionValidityResult
	if err != nil {
		res = primitives.NewTransactionValidityResult(err)
	} else {
		res = primitives.NewTransactionValidityResult(ok)
	}

	return utils.BytesToOffsetAndSize(res.Bytes())
}
