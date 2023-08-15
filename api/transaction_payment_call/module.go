package transaction_payment_call

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "TransactionPaymentCallApi"
	apiVersion    = 3
)

var (
	DefaultTip = sc.NewU128FromUint64(0)
)

type Module[N sc.Numeric] struct {
	decoder    types.ModuleDecoder[N]
	txPayments transaction_payment.Module[N]
}

func NewCallApi[N sc.Numeric](decoder types.ModuleDecoder[N], txPayments transaction_payment.Module[N]) Module[N] {
	return Module[N]{
		decoder:    decoder,
		txPayments: txPayments,
	}
}

func (m Module[N]) Name() string {
	return ApiModuleName
}

func (m Module[N]) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// QueryCallInfo queries the data of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded weight, dispatch class and partial fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-info)
func (m Module[N]) QueryCallInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := m.decoder.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	partialFee := m.txPayments.ComputeFee(length, dispatchInfo, DefaultTip)

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

// QueryCallFeeDetails queries the detailed fee of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded detailed fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-fee-details)
func (m Module[N]) QueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call := m.decoder.DecodeCall(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(call)
	feeDetails := m.txPayments.ComputeFeeDetails(length, dispatchInfo, DefaultTip)

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}
