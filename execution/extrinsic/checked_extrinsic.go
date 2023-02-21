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

type Extrinsic types.CheckedExtrinsic

func (xt Extrinsic) Validate(validator types.UnsignedValidator, source types.TransactionSource, info *types.DispatchInfo, length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.Extra
		_, _ = extra.Validate(&id, &xt.Function, info, length)
	} else {
		extra := &types.Extra{}
		valid, err := extra.ValidateUnsigned(&xt.Function, info, length)
		if err != nil {
			return ok, err
		}

		unsignedValidation, err := validator.ValidateUnsigned(source, &xt.Function)
		if err != nil {
			return ok, err
		}

		ok = valid.CombineWith(unsignedValidation)
	}

	return ok, err
}

func (xt Extrinsic) Apply(validator types.UnsignedValidator, info *types.DispatchInfo, length sc.Compact) (ok types.DispatchResultWithPostInfo[types.PostDispatchInfo], err types.TransactionValidityError) {
	var (
		maybeWho sc.Option[types.Address32]
		maybePre sc.Option[types.Pre]
	)

	if xt.Signed.HasValue {
		id, extra := xt.Signed.Value.Address32, xt.Signed.Value.Extra
		pre, err := types.Extra{}.PreDispatch(extra, &id, &xt.Function, info, length)
		if err != nil {
			return ok, err
		}
		maybeWho, maybePre = sc.NewOption[types.Address32](id), sc.NewOption[types.Pre](pre)
	} else {
		// Do any pre-flight stuff for an unsigned transaction.
		//
		// Note this function by default delegates to `ValidateUnsigned`, so that
		// all checks performed for the transaction queue are also performed during
		// the dispatch phase (applying the extrinsic).
		//
		// If you ever override this function, you need to make sure to always
		// perform the same validation as in `ValidateUnsigned`.
		_, err := types.Extra{}.PreDispatchUnsigned(&xt.Function, info, length)
		if err != nil {
			return ok, err
		}

		_, err = validator.PreDispatch(&xt.Function)
		if err != nil {
			return ok, err
		}

		maybeWho, maybePre = sc.NewOption[types.Address32](nil), sc.NewOption[types.Pre](nil)
	}

	postDispatchInfo, resWithInfo := Dispatch(xt.Function, types.RawOriginFrom(maybeWho))

	var postInfo types.PostDispatchInfo
	if resWithInfo.HasError {
		postInfo = resWithInfo.Err.PostInfo
	}
	postInfo = postDispatchInfo

	dispatchResult := types.NewDispatchResult(resWithInfo.Err)
	_, err = types.Extra{}.PostDispatch(maybePre, info, &postInfo, length, &dispatchResult)

	dispatchResultWithPostInfo := types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
	if resWithInfo.HasError {
		dispatchResultWithPostInfo.HasError = true
		dispatchResultWithPostInfo.Err = resWithInfo.Err
	} else {
		dispatchResultWithPostInfo.Ok = resWithInfo.Ok
	}

	return dispatchResultWithPostInfo, err
}

func Dispatch(call types.Call, maybeWho types.RuntimeOrigin) (ok types.PostDispatchInfo, err types.DispatchResultWithPostInfo[types.PostDispatchInfo]) {
	switch call.CallIndex.ModuleIndex {
	// TODO: Add more modules
	case system.Module.Index:
		switch call.CallIndex.FunctionIndex {
		// TODO: Add more functions
		case system.Module.Functions["remark"].Index:
			// TODO: Implement
		default:
			log.Critical("system.function with index " + string(call.CallIndex.ModuleIndex) + "not found")
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
			log.Critical("system.function with index " + string(call.CallIndex.ModuleIndex) + "not found")
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
		}
	default:
		log.Critical("module with index " + string(xt.Function.CallIndex.ModuleIndex) + "not found")
	}

	panic("unreachable")
}
