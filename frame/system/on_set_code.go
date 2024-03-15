package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type defaultOnSetCode struct {
	module Module
}

func NewDefaultOnSetCode(module Module) defaultOnSetCode {
	return defaultOnSetCode{module}
}

// What to do if the runtime wants to change the code to something new.
//
// The default implementation is responsible for setting the correct storage
// entry and emitting corresponding event and log item. (see
// It's unlikely that this needs to be customized, unless you are writing a parachain using
// `Cumulus`, where the actual code change is deferred.
func (d defaultOnSetCode) SetCode(codeBlob sc.Sequence[sc.U8]) error {
	d.updateCodeInStorage(codeBlob)
	return nil
}

// Write code to the storage and emit related events and digest items.
//
// Note this function almost never should be used directly. It is exposed
// for `OnSetCode` implementations that defer actual code being written to
// the storage (for instance in case of parachains).
func (d defaultOnSetCode) updateCodeInStorage(codeBlob sc.Sequence[sc.U8]) {
	d.module.StorageCodeSet(codeBlob)
	d.module.DepositLog(primitives.NewDigestItemRuntimeEnvironmentUpgrade())
	d.module.DepositEvent(primitives.NewEvent(d.module.GetIndex(), EventCodeUpdated))
}
