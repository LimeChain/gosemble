package tagged_transaction_queue

import (
	"bytes"

	"github.com/LimeChain/gosemble/frame/executive"
	"github.com/LimeChain/gosemble/primitives/types"
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
	buffer := &bytes.Buffer{}
	buffer.Write(data)

	txSource := types.DecodeTransactionSource(buffer)
	tx := types.DecodeUncheckedExtrinsic(buffer)
	blockHash := types.DecodeBlake2bHash(buffer)
	buffer.Reset()

	ok, err := executive.ValidateTransaction(txSource, tx, blockHash)

	var res types.TransactionValidityResult
	if err != nil {
		res = types.NewTransactionValidityResult(err)
	} else {
		res = types.NewTransactionValidityResult(ok)
	}

	return utils.BytesToOffsetAndSize(res.Bytes())
}
