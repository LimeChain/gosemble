package transaction_payment_call

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/transaction_payment"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "TransactionPaymentCallApi"
	apiVersion    = 3
)

type Module struct {
	decoder    types.RuntimeDecoder
	txPayments transaction_payment.Module
	memUtils   utils.WasmMemoryTranslator
}

func New(decoder types.RuntimeDecoder, txPayments transaction_payment.Module) Module {
	return Module{
		decoder:    decoder,
		txPayments: txPayments,
		memUtils:   utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash, err := hashing.Blake2b8([]byte(ApiModuleName))
	if err != nil {
		log.Critical(err.Error())
	}
	return primitives.NewApiItem(hash, apiVersion)
}

// QueryCallInfo queries the data of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded weight, dispatch class and partial fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-info)
func (m Module) QueryCallInfo(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call, err := m.decoder.DecodeCall(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	length, err := sc.DecodeU32(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	dispatchInfo := primitives.GetDispatchInfo(call)
	partialFee, err := m.txPayments.ComputeFee(length, dispatchInfo, constants.DefaultTip)
	if err != nil {
		log.Critical(err.Error())
	}

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return m.memUtils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

// QueryCallFeeDetails queries the detailed fee of a dispatch call.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded dispatch call and its length.
// Returns a pointer-size of the SCALE-encoded detailed fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentcallapi-query-call-fee-details)
func (m Module) QueryCallFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := m.memUtils.GetWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	call, err := m.decoder.DecodeCall(buffer)
	if err != nil {
		log.Critical(err.Error())
	}
	length, err := sc.DecodeU32(buffer)
	if err != nil {
		log.Critical(err.Error())
	}

	dispatchInfo := primitives.GetDispatchInfo(call)
	feeDetails, err := m.txPayments.ComputeFeeDetails(length, dispatchInfo, constants.DefaultTip)
	if err != nil {
		log.Critical(err.Error())
	}

	return m.memUtils.BytesToOffsetAndSize(feeDetails.Bytes())
}
