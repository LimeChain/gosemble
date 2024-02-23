package system

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type EventDepositor interface {
	DepositEvent(event primitives.Event)
}

type LogDepositor interface {
	DepositLog(item primitives.DigestItem)
}

type CodeUpgrader interface {
	CanSetCode(codeBlob sc.Sequence[sc.U8]) error
	DoAuthorizeUpgrade(codeHash primitives.H256, checkVersion sc.Bool)
	DoApplyAuthorizeUpgrade(codeBlob sc.Sequence[sc.U8]) (primitives.PostDispatchInfo, error)
}

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

// type Key = sc.Sequence[sc.U8]

type KeyValue struct {
	Key   sc.Sequence[sc.U8]
	Value sc.Sequence[sc.U8]
}

func (pair KeyValue) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer, pair.Key, pair.Value)
}

func (pair KeyValue) Bytes() []byte {
	return sc.EncodedBytes(pair)
}

func DecodeKeyValue(buffer *bytes.Buffer) (KeyValue, error) {
	key, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return KeyValue{}, err
	}

	value, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return KeyValue{}, err
	}

	return KeyValue{key, value}, nil
}

// Information needed when a new runtime binary is submitted and needs to be authorized before
// replacing the current runtime.
type CodeUpgradeAuthorization struct {
	// Hash of the new runtime binary.
	CodeHash primitives.H256
	// Whether or not to carry out version checks.
	CheckVersion sc.Bool
}

func (c CodeUpgradeAuthorization) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer, c.CodeHash, c.CheckVersion)
}

func (c CodeUpgradeAuthorization) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func DecodeCodeUpgradeAuthorization(buffer *bytes.Buffer) (CodeUpgradeAuthorization, error) {
	codeHash, err := primitives.DecodeH256(buffer)
	if err != nil {
		return CodeUpgradeAuthorization{}, err
	}
	checkVersion, err := sc.DecodeBool(buffer)
	if err != nil {
		return CodeUpgradeAuthorization{}, err
	}
	return CodeUpgradeAuthorization{codeHash, checkVersion}, nil
}

func DecodeCodeUpgradeAuthorizationOption(buffer *bytes.Buffer) (sc.Option[CodeUpgradeAuthorization], error) {
	return sc.DecodeOptionWith(buffer, DecodeCodeUpgradeAuthorization)
}
