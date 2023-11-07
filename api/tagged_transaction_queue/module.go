package tagged_transaction_queue

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
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
	decoder   types.RuntimeDecoder
	memUtils  utils.WasmMemoryTranslator
}

func New(executive executive.Module, decoder types.RuntimeDecoder) Module {
	return Module{
		executive: executive,
		decoder:   decoder,
		memUtils:  utils.NewMemoryTranslator(),
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
	data := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	txSource, err := primitives.DecodeTransactionSource(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	tx, err := m.decoder.DecodeUncheckedExtrinsic(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	blockHash, err := primitives.DecodeBlake2bHash(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	ok, errTx := m.executive.ValidateTransaction(txSource, tx, blockHash)

	var res primitives.TransactionValidityResult
	if errTx != nil {
		res, err = primitives.NewTransactionValidityResult(errTx)
		if err != nil {
			log.Critical(err.Error())
		}
	} else {
		res, err = primitives.NewTransactionValidityResult(ok)
		if err != nil {
			log.Critical(err.Error())
		}
	}

	return m.memUtils.BytesToOffsetAndSize(res.Bytes())
}
