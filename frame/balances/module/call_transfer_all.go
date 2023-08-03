package module

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type transferAllCall struct {
	primitives.Callable
	transfer
}

func newTransferAllCall(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts) primitives.Call {
	call := transferAllCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants),
	}

	return call
}

func (c transferAllCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeBool(buffer),
	)
	return c
}

func (c transferAllCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c transferAllCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c transferAllCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c transferAllCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c transferAllCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ transferAllCall) IsInherent() bool {
	return false
}

func (_ transferAllCall) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 34_878 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(35_121_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ transferAllCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ transferAllCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ transferAllCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c transferAllCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := c.transferAll(origin, args[0].(types.MultiAddress), bool(args[1].(sc.Bool)))
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}

// transferAll transfers the entire transferable balance from `origin` to `dest`.
// By transferable it means that any locked or reserved amounts will not be transferred.
// `keepAlive`: A boolean to determine if the `transfer_all` operation should send all
// the funds the account has, causing the sender account to be killed (false), or
// transfer everything except at least the existential deposit, which will guarantee to
// keep the sender account alive (true).
func (c transferAllCall) transferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	transactor := origin.AsSigned()
	reducibleBalance := c.reducibleBalance(transactor, keepAlive)

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		log.Debug("Failed to lookup [" + string(dest.Bytes()) + "]") // fmt.Sprintf("Failed to lookup [%s]", dest.Bytes())
		return types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return c.transfer.trans(transactor, to, reducibleBalance, keep)
}
