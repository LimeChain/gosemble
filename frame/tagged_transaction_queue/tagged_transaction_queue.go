package tagged_transaction_queue

import (
	"bytes"

	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/executive"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

type TaggedTransactionQueue interface {
	ValidateTransaction(dataPtr int32, dataLen int32) int64
}

/*
https://spec.polkadot.network/#sect-rte-validate-transaction

	dataPtr - i32 pointer to the memory location.
	dataLen - i32 length (in bytes) of the encoded arguments.
	returns a pointer-size to the SCALE-encoded ([]byte) data.
*/
func ValidateTransaction(dataPtr int32, dataLen int32) int64 {
	data := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(data)

	txSource := primitives.DecodeTransactionSource(buffer)
	tx := types.DecodeUncheckedExtrinsic(buffer)
	blockHash := primitives.DecodeBlake2bHash(buffer)

	ok, err := executive.ValidateTransaction(txSource, tx, blockHash)

	var res primitives.TransactionValidityResult
	if err != nil {
		res = primitives.NewTransactionValidityResult(err)
	} else {
		res = primitives.NewTransactionValidityResult(ok)
	}

	return utils.BytesToOffsetAndSize(res.Bytes())
}
