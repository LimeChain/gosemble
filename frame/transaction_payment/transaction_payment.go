package transaction_payment

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/frame/transaction_payment/module"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

var DefaultMultiplierValue = sc.NewU128FromUint64(1)
var DefaultTip = sc.NewU128FromUint64(0)

const (
	ApiModuleName = "TransactionPaymentApi"
	apiVersion    = 3
)

type Module struct {
	decoder    types.ModuleDecoder
	txPayments module.TransactionPaymentModule
}

func New(decoder types.ModuleDecoder, txPayments module.TransactionPaymentModule) Module {
	return Module{
		decoder:    decoder,
		txPayments: txPayments,
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

// QueryInfo queries the data of an extrinsic.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded extrinsic and its length.
// Returns a pointer-size of the SCALE-encoded weight, dispatch class and partial fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentapi-query-info)
func (m Module) QueryInfo(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := m.decoder.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	partialFee := sc.NewU128FromUint64(0)
	if ext.IsSigned() {
		partialFee = m.txPayments.ComputeFee(length, dispatchInfo, DefaultTip)
	}

	runtimeDispatchInfo := primitives.RuntimeDispatchInfo{
		Weight:     dispatchInfo.Weight,
		Class:      dispatchInfo.Class,
		PartialFee: partialFee,
	}

	return utils.BytesToOffsetAndSize(runtimeDispatchInfo.Bytes())
}

// QueryFeeDetails queries the detailed fee of an extrinsic.
// It takes two arguments:
// - dataPtr: Pointer to the data in the Wasm memory.
// - dataLen: Length of the data.
// which represent the SCALE-encoded extrinsic and its length.
// Returns a pointer-size of the SCALE-encoded detailed fee.
// [Specification](https://spec.polkadot.network/chap-runtime-api#sect-rte-transactionpaymentapi-query-fee-details)
func (m Module) QueryFeeDetails(dataPtr int32, dataLen int32) int64 {
	b := utils.ToWasmMemorySlice(dataPtr, dataLen)
	buffer := bytes.NewBuffer(b)

	ext := m.decoder.DecodeUncheckedExtrinsic(buffer)
	length := sc.DecodeU32(buffer)

	dispatchInfo := primitives.GetDispatchInfo(ext.Function)

	var feeDetails primitives.FeeDetails
	if ext.IsSigned() {
		feeDetails = m.txPayments.ComputeFeeDetails(length, dispatchInfo, DefaultTip)
	} else {
		feeDetails = primitives.FeeDetails{
			InclusionFee: sc.NewOption[primitives.InclusionFee](nil),
		}
	}

	return utils.BytesToOffsetAndSize(feeDetails.Bytes())
}
