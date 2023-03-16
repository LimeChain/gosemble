package extrinsic

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

func GetDispatchInfo(xt types.CheckedExtrinsic) types.DispatchInfo {
	// TODO: add more module functions
	switch xt.Function.CallIndex.ModuleIndex {
	case system.Module.Index():
		switch xt.Function.CallIndex.FunctionIndex {
		case system.Module.Remark.Index():
			baseWeight := system.Module.Remark.BaseWeight(xt.Function.Args)

			return types.DispatchInfo{
				Weight:  system.Module.Remark.WeightInfo(baseWeight, xt.Function.Args),
				Class:   system.Module.Remark.ClassifyDispatch(baseWeight, xt.Function.Args),
				PaysFee: system.Module.Remark.PaysFee(baseWeight, xt.Function.Args),
			}
		}

	case timestamp.Module.Index():
		switch xt.Function.CallIndex.FunctionIndex {
		case timestamp.Module.Set.Index():
			baseWeight := timestamp.Module.Set.BaseWeight(xt.Function.Args)

			return types.DispatchInfo{
				Weight:  timestamp.Module.Set.WeightInfo(baseWeight, xt.Function.Args),
				Class:   timestamp.Module.Set.ClassifyDispatch(baseWeight, xt.Function.Args),
				PaysFee: timestamp.Module.Set.PaysFee(baseWeight, xt.Function.Args),
			}
		}

	default:
		log.Trace(fmt.Sprintf("module with index %d not found", xt.Function.CallIndex.ModuleIndex))
	}

	log.Trace(fmt.Sprintf("function with index %d not found", xt.Function.CallIndex.FunctionIndex))
	return types.DispatchInfo{
		Weight:  types.WeightFromParts(sc.U64(len(xt.Bytes())), sc.U64(0)),
		Class:   types.NewDispatchClassNormal(),
		PaysFee: types.NewPaysYes(),
	}
}

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) (ok types.PostDispatchInfo, err types.DispatchResultWithPostInfo[types.PostDispatchInfo]) {
	// TODO: Add more modules and functions
	switch call.CallIndex.ModuleIndex {
	case system.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case system.Module.Remark.Index():
			res := system.Module.Remark.Dispatch(maybeWho, call.Args)
			if res.HasError {
				err = res
				return ok, err
			}
			ok = res.Ok
		default:
			log.Trace(fmt.Sprintf("function index %d not found", call.CallIndex.FunctionIndex))
		}
	case timestamp.Module.Index():
		switch call.CallIndex.FunctionIndex {
		case timestamp.Module.Set.Index():
			buffer := &bytes.Buffer{}
			buffer.Write(sc.SequenceU8ToBytes(call.Args))
			compactTs := sc.DecodeCompact(buffer)
			ts := sc.U64(compactTs.ToBigInt().Uint64())
			timestamp.Module.Set.Dispatch(types.NewRawOriginNone(), ts)
		default:
			log.Trace(fmt.Sprintf("function index %d not found", call.CallIndex.FunctionIndex))
		}

	default:
		log.Trace(fmt.Sprintf("module with index %d not found", call.CallIndex.ModuleIndex))
	}

	return ok, err
}
