package extrinsic

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	balances_constants "github.com/LimeChain/gosemble/constants/balances"
	system_constants "github.com/LimeChain/gosemble/constants/system"
	timestamp_constants "github.com/LimeChain/gosemble/constants/timestamp"
	"github.com/LimeChain/gosemble/frame/balances"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	switch call.CallIndex.ModuleIndex {
	// TODO: Add more modules
	case system.Module.Index:
		switch call.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case system_constants.FunctionRemarkIndex:
			// TODO: Implement
		default:
			log.Critical("system.function with index " + string(call.CallIndex.FunctionIndex) + "not found")
		}
	case timestamp.Module.Index:
		switch call.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case timestamp_constants.FunctionSetIndex:
			buffer := &bytes.Buffer{}
			buffer.Write(sc.SequenceU8ToBytes(call.Args))
			compactTs := sc.DecodeCompact(buffer)
			ts := sc.U64(compactTs.ToBigInt().Uint64())

			timestamp.Set(ts)
		default:
			log.Critical("timestamp.function with index " + string(call.CallIndex.FunctionIndex) + "not found")
		}
	case balances.Module.Index:
		buffer := &bytes.Buffer{}
		buffer.Write(sc.SequenceU8ToBytes(call.Args))

		switch call.CallIndex.FunctionIndex {
		case balances_constants.FunctionTransferIndex:
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			return balances.Transfer(maybeWho, to, sc.U128(value))
		case balances_constants.FunctionSetBalanceIndex:
			to := types.DecodeMultiAddress(buffer)
			newFree := sc.DecodeCompact(buffer)
			newReserved := sc.DecodeCompact(buffer)

			return balances.SetBalance(maybeWho, to, newFree.ToBigInt(), newReserved.ToBigInt())
		case balances_constants.FunctionForceTransferIndex:
			from := types.DecodeMultiAddress(buffer)
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			return balances.ForceTransfer(maybeWho, from, to, sc.U128(value))
		case balances_constants.FunctionTransferKeepAliveIndex:
			destination := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			return balances.TransferKeepAlive(maybeWho, destination, sc.U128(value))
		case balances_constants.FunctionTransferAllIndex:
			destination := types.DecodeMultiAddress(buffer)
			keepAlive := sc.DecodeBool(buffer)

			result := balances.TransferAll(maybeWho, destination, bool(keepAlive))
			if result != nil {
				return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
					HasError: true,
					Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
						PostInfo: types.PostDispatchInfo{
							ActualWeight: sc.Option[types.Weight]{
								HasValue: false,
							},
							PaysFee: 0,
						},
						Error: result,
					},
				}
			}
		case balances_constants.FunctionForceFreeIndex:
			to := types.DecodeMultiAddress(buffer)
			value := sc.DecodeCompact(buffer)

			result := balances.ForceFree(maybeWho, to, value.ToBigInt())
			if result != nil {
				return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
					HasError: true,
					Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
						PostInfo: types.PostDispatchInfo{
							ActualWeight: sc.Option[types.Weight]{
								HasValue: false,
							},
							PaysFee: 0,
						},
						Error: result,
					},
				}
			}
		}

	default:
		log.Critical("module with index " + string(call.CallIndex.ModuleIndex) + "not found")
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
}

func GetDispatchInfo(xt types.CheckedExtrinsic) types.DispatchInfo {
	switch xt.Function.CallIndex.ModuleIndex {
	// TODO: Add more modules
	case system.Module.Index:
		// TODO: Implement
		return types.DispatchInfo{
			Weight:  types.WeightFromRefTime(sc.U64(len(xt.Bytes()))),
			Class:   types.NewDispatchClassNormal(),
			PaysFee: types.NewPaysYes(),
		}

	case timestamp.Module.Index:
		switch xt.Function.CallIndex.FunctionIndex {
		case timestamp_constants.FunctionSetIndex:
			baseWeight := timestamp.Module.Functions[xt.Function.CallIndex.FunctionIndex].BaseWeight
			weight := support.WeighData(baseWeight, xt.Function.Args)
			class := support.ClassifyDispatch(baseWeight)
			paysFee := support.PaysFee(baseWeight)

			return types.DispatchInfo{
				Weight:  weight,
				Class:   class,
				PaysFee: paysFee,
			}
		default:
			log.Critical("system.function with index " + string(xt.Function.CallIndex.ModuleIndex) + "not found")
		}
	case balances.Module.Index:
		switch xt.Function.CallIndex.FunctionIndex {
		case balances_constants.FunctionTransferIndex,
			balances_constants.FunctionSetBalanceIndex,
			balances_constants.FunctionForceTransferIndex,
			balances_constants.FunctionTransferKeepAliveIndex,
			balances_constants.FunctionTransferAllIndex,
			balances_constants.FunctionForceFreeIndex:

			baseWeight := timestamp.Module.Functions[xt.Function.CallIndex.FunctionIndex].BaseWeight
			weight := support.WeighData(baseWeight, xt.Function.Args)
			class := support.ClassifyDispatch(baseWeight)
			paysFee := support.PaysFee(baseWeight)

			return types.DispatchInfo{
				Weight:  weight,
				Class:   class,
				PaysFee: paysFee,
			}
		default:
			log.Critical("system.function with index " + string(xt.Function.CallIndex.ModuleIndex) + "not found")
		}
	default:
		log.Critical("module with index " + string(xt.Function.CallIndex.ModuleIndex) + "not found")
	}

	panic("unreachable")
}
