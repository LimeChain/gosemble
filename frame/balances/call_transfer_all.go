package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTransferAll struct {
	primitives.Callable
	transfer
}

func newCallTransferAll(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callTransferAll{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callTransferAll) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	dest, err := types.DecodeMultiAddress[testKeyType](buffer)
	if err != nil {
		return nil, err
	}
	keepAlive, err := sc.DecodeBool(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		dest,
		keepAlive,
	)
	return c, nil
}

func (c callTransferAll) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callTransferAll) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTransferAll) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTransferAll) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTransferAll) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callTransferAll) BaseWeight() types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 34_878 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(35_121_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callTransferAll) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransferAll) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callTransferAll) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callTransferAll) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
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
func (c callTransferAll) transferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	transactor, err := origin.AsSigned()
	if err != nil {
		log.Critical(err.Error())
	}

	reducibleBalance, err := c.reducibleBalance(transactor, keepAlive)
	if err != nil {
		return primitives.NewDispatchErrorOther(sc.Str(err.Error()))
	}

	to, errLookup := types.Lookup(dest)
	if errLookup != nil {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		log.Debug("Failed to lookup [" + string(dest.Bytes()) + "]")
		return types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return c.transfer.trans(transactor, to, reducibleBalance, keep)
}
