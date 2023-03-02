package extrinsic

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) (ok types.PostDispatchInfo, err types.DispatchResultWithPostInfo[types.PostDispatchInfo]) {
	switch call.CallIndex.ModuleIndex {
	// TODO: Add more modules
	case system.Module.Index:
		switch call.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case system.Module.Functions["remark"].Index:
			// TODO: Implement
		default:
			log.Critical("system.function with index " + string(call.CallIndex.FunctionIndex) + "not found")
		}
	case timestamp.Module.Index:
		switch call.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case timestamp.Module.Functions["set"].Index:
			buffer := &bytes.Buffer{}
			buffer.Write(sc.SequenceU8ToBytes(call.Args))
			compactTs := sc.DecodeCompact(buffer)
			ts := sc.U64(compactTs.ToBigInt().Uint64())

			timestamp.Set(ts)
		default:
			log.Critical("timestamp.function with index " + string(call.CallIndex.FunctionIndex) + "not found")
		}

	default:
		log.Critical("module with index " + string(call.CallIndex.ModuleIndex) + "not found")
	}

	return ok, err
}

func GetDispatchInfo(xt types.CheckedExtrinsic) types.DispatchInfo {
	switch xt.Function.CallIndex.ModuleIndex {
	// TODO: Add more modules
	case system.Module.Index:
		// TODO: Implement
		return types.DispatchInfo{
			Weight:  types.WeightFromRefTime(sc.U64(len(xt.Bytes()))),
			Class:   types.NewDispatchClass(types.NormalDispatch),
			PaysFee: types.NewPays(types.PaysYes),
		}

	case timestamp.Module.Index:
		switch xt.Function.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case timestamp.Module.Functions["set"].Index:
			baseWeight := timestamp.Module.Functions["set"].BaseWeight
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
