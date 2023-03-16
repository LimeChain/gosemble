package extrinsic

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances"
	"github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

func GetDispatchInfo(xt types.CheckedExtrinsic) types.DispatchInfo {
	// TODO: add more module functions
	var function support.FunctionMetadata
	switch xt.Function.CallIndex.ModuleIndex {
	case system.Module.Index():
		function = system.Module.Functions()[xt.Function.CallIndex.FunctionIndex]
	case timestamp.Module.Index():
		function = timestamp.Module.Functions()[xt.Function.CallIndex.FunctionIndex]
	case balances.Module.Index():
		function = balances.Module.Functions()[xt.Function.CallIndex.FunctionIndex]
	default:
		log.Trace(fmt.Sprintf("module with index %d not found", xt.Function.CallIndex.ModuleIndex))
	}

	log.Trace(fmt.Sprintf("function with index %d not found", xt.Function.CallIndex.FunctionIndex))

	baseWeight := function.BaseWeight(xt.Function.Args)
	return types.DispatchInfo{
		Weight:  timestamp.Module.Set.WeightInfo(baseWeight, xt.Function.Args),
		Class:   timestamp.Module.Set.ClassifyDispatch(baseWeight, xt.Function.Args),
		PaysFee: timestamp.Module.Set.PaysFee(baseWeight, xt.Function.Args),
	}
}

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	// TODO: Add more modules and functions
	switch call.CallIndex.ModuleIndex {
	case system.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case system.Module.Remark.Index():
			return system.Module.Remark.Dispatch(maybeWho, call.Args)
		default:
			log.Trace(fmt.Sprintf("function index %d not found", call.CallIndex.FunctionIndex))
		}
	case timestamp.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case timestamp.Module.Set.Index():
			buffer := &bytes.Buffer{}
			buffer.Write(call.Args)
			compactTs := sc.DecodeCompact(buffer)
			ts := sc.U64(compactTs.ToBigInt().Uint64())
			timestamp.Module.Set.Dispatch(types.NewRawOriginNone(), ts)
		default:
			log.Trace(fmt.Sprintf("function index %d not found", call.CallIndex.FunctionIndex))
		}
	case balances.Module.Index():
		buffer := &bytes.Buffer{}
		buffer.Write(call.Args)

		var err types.DispatchError
		switch call.CallIndex.FunctionIndex {
		case constants.FunctionTransferIndex:
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			_, err = balances.Module.Transfer.Dispatch(maybeWho, to, sc.U128(value))
		case constants.FunctionSetBalanceIndex:
			to := types.DecodeMultiAddress(buffer)
			newFree := sc.DecodeCompact(buffer)
			newReserved := sc.DecodeCompact(buffer)

			_, err = balances.Module.SetBalance.Dispatch(maybeWho, to, newFree.ToBigInt(), newReserved.ToBigInt())
		case constants.FunctionForceTransferIndex:
			from := types.DecodeMultiAddress(buffer)
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			_, err = balances.Module.ForceTransfer.Dispatch(maybeWho, from, to, sc.U128(value))
		case constants.FunctionTransferKeepAliveIndex:
			destination := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			_, err = balances.Module.TransferKeepAlive.Dispatch(maybeWho, destination, sc.U128(value))
		case constants.FunctionTransferAllIndex:
			destination := types.DecodeMultiAddress(buffer)
			keepAlive := sc.DecodeBool(buffer)

			_, err = balances.Module.TransferAll.Dispatch(maybeWho, destination, bool(keepAlive))
		case constants.FunctionForceFreeIndex:
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			_, err = balances.Module.ForceFree.Dispatch(maybeWho, to, value.ToBigInt())
		}
		if err != nil {
			return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
				HasError: true,
				Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
					Error: err,
				},
			}
		}

	default:
		log.Trace(fmt.Sprintf("module with index %d not found", call.CallIndex.ModuleIndex))
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
}
