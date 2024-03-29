package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	ModuleVersion14 sc.U8 = 14
	ModuleVersion15 sc.U8 = 15
)

type MetadataModule struct {
	Version   sc.U8
	ModuleV14 MetadataModuleV14
	ModuleV15 MetadataModuleV15
}

func (m MetadataModule) Encode(buffer *bytes.Buffer) error {
	switch m.Version {
	case ModuleVersion14:
		return m.ModuleV14.Encode(buffer)
	case ModuleVersion15:
		return m.ModuleV15.Encode(buffer)
	}

	return newTypeError("ModuleVersion")
}

func (mm MetadataModule) Bytes() []byte {
	return sc.EncodedBytes(mm)
}

type MetadataModuleV15 struct {
	Name      sc.Str
	Storage   sc.Option[MetadataModuleStorage]
	Call      sc.Option[sc.Compact]
	CallDef   sc.Option[MetadataDefinitionVariant] // not encoded
	Event     sc.Option[sc.Compact]
	EventDef  sc.Option[MetadataDefinitionVariant] // not encoded
	Constants sc.Sequence[MetadataModuleConstant]
	Error     sc.Option[sc.Compact]
	ErrorDef  sc.Option[MetadataDefinitionVariant] // not encoded
	Index     sc.U8
	Docs      sc.Sequence[sc.Str]
}

func (mm MetadataModuleV15) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		mm.Name,
		mm.Storage,
		mm.Call,
		mm.Event,
		mm.Constants,
		mm.Error,
		mm.Index,
		mm.Docs,
	)
}

func DecodeMetadataModuleV15(buffer *bytes.Buffer) (MetadataModuleV15, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	storage, err := sc.DecodeOptionWith(buffer, DecodeMetadataModuleStorage)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	call, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	event, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	constants, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleConstant)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	e, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	index, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataModuleV15{}, err
	}
	return MetadataModuleV15{
		Name:      name,
		Storage:   storage,
		Call:      call,
		Event:     event,
		Constants: constants,
		Error:     e,
		Index:     index,
		Docs:      docs,
	}, nil
}

func (mm MetadataModuleV15) Bytes() []byte {
	return sc.EncodedBytes(mm)
}

type MetadataModuleV14 struct {
	Name      sc.Str
	Storage   sc.Option[MetadataModuleStorage]
	Call      sc.Option[sc.Compact]
	CallDef   sc.Option[MetadataDefinitionVariant] // not encoded
	Event     sc.Option[sc.Compact]
	EventDef  sc.Option[MetadataDefinitionVariant] // not encoded
	Constants sc.Sequence[MetadataModuleConstant]
	Error     sc.Option[sc.Compact]
	ErrorDef  sc.Option[MetadataDefinitionVariant] // not encoded
	Index     sc.U8
}

func (mm MetadataModuleV14) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		mm.Name,
		mm.Storage,
		mm.Call,
		mm.Event,
		mm.Constants,
		mm.Error,
		mm.Index,
	)
}

func DecodeMetadataModuleV14(buffer *bytes.Buffer) (MetadataModuleV14, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	storage, err := sc.DecodeOptionWith(buffer, DecodeMetadataModuleStorage)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	call, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	event, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	constants, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleConstant)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	e, err := sc.DecodeOption[sc.Compact](buffer)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	index, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataModuleV14{}, err
	}
	return MetadataModuleV14{
		Name:      name,
		Storage:   storage,
		Call:      call,
		Event:     event,
		Constants: constants,
		Error:     e,
		Index:     index,
	}, nil
}

func (mm MetadataModuleV14) Bytes() []byte {
	return sc.EncodedBytes(mm)
}

type MetadataModuleStorage struct {
	Prefix sc.Str
	Items  sc.Sequence[MetadataModuleStorageEntry]
}

func (mms MetadataModuleStorage) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		mms.Prefix,
		mms.Items,
	)
}

func DecodeMetadataModuleStorage(buffer *bytes.Buffer) (MetadataModuleStorage, error) {
	prefix, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataModuleStorage{}, err
	}
	items, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleStorageEntry)
	if err != nil {
		return MetadataModuleStorage{}, err
	}
	return MetadataModuleStorage{
		Prefix: prefix,
		Items:  items,
	}, nil
}

func (mms MetadataModuleStorage) Bytes() []byte {
	return sc.EncodedBytes(mms)
}

type MetadataModuleStorageEntry struct {
	Name       sc.Str
	Modifier   MetadataModuleStorageEntryModifier
	Definition MetadataModuleStorageEntryDefinition
	Fallback   sc.Sequence[sc.U8]
	Docs       sc.Sequence[sc.Str]
}

func NewMetadataModuleStorageEntry(name string, modifier MetadataModuleStorageEntryModifier, definition MetadataModuleStorageEntryDefinition, docs string) MetadataModuleStorageEntry {
	return MetadataModuleStorageEntry{
		Name:       sc.Str(name),
		Modifier:   modifier,
		Definition: definition,
		Fallback:   sc.Sequence[sc.U8]{},
		Docs:       sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mmse MetadataModuleStorageEntry) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		mmse.Name,
		mmse.Modifier,
		mmse.Definition,
		mmse.Fallback,
		mmse.Docs,
	)
}

func DecodeMetadataModuleStorageEntry(buffer *bytes.Buffer) (MetadataModuleStorageEntry, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataModuleStorageEntry{}, err
	}
	mod, err := DecodeMetadataModuleStorageEntryModifier(buffer)
	if err != nil {
		return MetadataModuleStorageEntry{}, err
	}
	def, err := DecodeMetadataModuleStorageEntryDefinition(buffer)
	if err != nil {
		return MetadataModuleStorageEntry{}, err
	}
	fallback, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return MetadataModuleStorageEntry{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataModuleStorageEntry{}, err
	}
	return MetadataModuleStorageEntry{
		Name:       name,
		Modifier:   mod,
		Definition: def,
		Fallback:   fallback,
		Docs:       docs,
	}, nil
}

func (mmse MetadataModuleStorageEntry) Bytes() []byte {
	return sc.EncodedBytes(mmse)
}

type MetadataModuleStorageEntryModifier sc.U8

const (
	MetadataModuleStorageEntryModifierOptional MetadataModuleStorageEntryModifier = iota
	MetadataModuleStorageEntryModifierDefault                                     = 1
)

func (mmsem MetadataModuleStorageEntryModifier) Encode(buffer *bytes.Buffer) error {
	return sc.U8(mmsem).Encode(buffer)
}

func DecodeMetadataModuleStorageEntryModifier(buffer *bytes.Buffer) (MetadataModuleStorageEntryModifier, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataModuleStorageEntryModifier(0), err
	}

	switch MetadataModuleStorageEntryModifier(b) {
	case MetadataModuleStorageEntryModifierOptional:
		return MetadataModuleStorageEntryModifierOptional, nil
	case MetadataModuleStorageEntryModifierDefault:
		return MetadataModuleStorageEntryModifierDefault, nil
	default:
		return MetadataModuleStorageEntryModifier(0), newTypeError("MetadataModuleStorageEntryModifier")
	}
}

func (mmsem MetadataModuleStorageEntryModifier) Bytes() []byte {
	return sc.EncodedBytes(mmsem)
}

const (
	MetadataModuleStorageEntryDefinitionPlain sc.U8 = iota
	MetadataModuleStorageEntryDefinitionMap
)

type MetadataModuleStorageEntryDefinition struct {
	sc.VaryingData
}

func NewMetadataModuleStorageEntryDefinitionPlain(key sc.Compact) MetadataModuleStorageEntryDefinition {
	return MetadataModuleStorageEntryDefinition{sc.NewVaryingData(MetadataModuleStorageEntryDefinitionPlain, key)}
}

func NewMetadataModuleStorageEntryDefinitionMap(storageHashFuncs sc.Sequence[MetadataModuleStorageHashFunc], key, value sc.Compact) MetadataModuleStorageEntryDefinition {
	return MetadataModuleStorageEntryDefinition{sc.NewVaryingData(MetadataModuleStorageEntryDefinitionMap, storageHashFuncs, key, value)}
}

func DecodeMetadataModuleStorageEntryDefinition(buffer *bytes.Buffer) (MetadataModuleStorageEntryDefinition, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataModuleStorageEntryDefinition{}, err
	}

	switch b {
	case MetadataModuleStorageEntryDefinitionPlain:
		key, err := sc.DecodeCompact[sc.U128](buffer)
		if err != nil {
			return MetadataModuleStorageEntryDefinition{}, err
		}
		return NewMetadataModuleStorageEntryDefinitionPlain(key), nil
	case MetadataModuleStorageEntryDefinitionMap:
		storageHashFuncs, err := sc.DecodeSequenceWith(buffer, DecodeMetadataModuleStorageHashFunc)
		if err != nil {
			return MetadataModuleStorageEntryDefinition{}, err
		}
		key, err := sc.DecodeCompact[sc.U128](buffer)
		if err != nil {
			return MetadataModuleStorageEntryDefinition{}, err
		}
		value, err := sc.DecodeCompact[sc.U128](buffer)
		if err != nil {
			return MetadataModuleStorageEntryDefinition{}, err
		}
		return NewMetadataModuleStorageEntryDefinitionMap(storageHashFuncs, key, value), nil
	default:
		return MetadataModuleStorageEntryDefinition{}, newTypeError("MetadataModuleStorageEntryDefinition")
	}
}

type MetadataModuleConstant struct {
	Name  sc.Str
	Type  sc.Compact
	Value sc.Sequence[sc.U8]
	Docs  sc.Sequence[sc.Str]
}

func NewMetadataModuleConstant(name string, id sc.Compact, value sc.Sequence[sc.U8], docs string) MetadataModuleConstant {
	return MetadataModuleConstant{
		Name:  sc.Str(name),
		Type:  id,
		Value: value,
		Docs:  sc.Sequence[sc.Str]{sc.Str(docs)},
	}
}

func (mmc MetadataModuleConstant) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		mmc.Name,
		mmc.Type,
		mmc.Value,
		mmc.Docs,
	)
}

func DecodeMetadataModuleConstant(buffer *bytes.Buffer) (MetadataModuleConstant, error) {
	name, err := sc.DecodeStr(buffer)
	if err != nil {
		return MetadataModuleConstant{}, err
	}
	t, err := sc.DecodeCompact[sc.U128](buffer)
	if err != nil {
		return MetadataModuleConstant{}, err
	}
	val, err := sc.DecodeSequence[sc.U8](buffer)
	if err != nil {
		return MetadataModuleConstant{}, err
	}
	docs, err := sc.DecodeSequence[sc.Str](buffer)
	if err != nil {
		return MetadataModuleConstant{}, err
	}
	return MetadataModuleConstant{
		Name:  name,
		Type:  t,
		Value: val,
		Docs:  docs,
	}, nil
}

func (mmc MetadataModuleConstant) Bytes() []byte {
	return sc.EncodedBytes(mmc)
}

type MetadataModuleStorageHashFunc sc.U8

const (
	MetadataModuleStorageHashFuncBlake128 MetadataModuleStorageHashFunc = iota
	MetadataModuleStorageHashFuncBlake256
	MetadataModuleStorageHashFuncMultiBlake128Concat
	MetadataModuleStorageHashFuncXX128
	MetadataModuleStorageHashFuncXX256
	MetadataModuleStorageHashFuncMultiXX64
	MetadataModuleStorageHashFuncIdentity
)

func (mmshf MetadataModuleStorageHashFunc) Encode(buffer *bytes.Buffer) error {
	return sc.U8(mmshf).Encode(buffer)
}

func DecodeMetadataModuleStorageHashFunc(buffer *bytes.Buffer) (MetadataModuleStorageHashFunc, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return MetadataModuleStorageHashFunc(0), err
	}

	switch MetadataModuleStorageHashFunc(b) {
	case MetadataModuleStorageHashFuncBlake128:
		return MetadataModuleStorageHashFuncBlake128, nil
	case MetadataModuleStorageHashFuncBlake256:
		return MetadataModuleStorageHashFuncBlake256, nil
	case MetadataModuleStorageHashFuncMultiBlake128Concat:
		return MetadataModuleStorageHashFuncMultiBlake128Concat, nil
	case MetadataModuleStorageHashFuncXX128:
		return MetadataModuleStorageHashFuncXX128, nil
	case MetadataModuleStorageHashFuncXX256:
		return MetadataModuleStorageHashFuncXX256, nil
	case MetadataModuleStorageHashFuncMultiXX64:
		return MetadataModuleStorageHashFuncMultiXX64, nil
	case MetadataModuleStorageHashFuncIdentity:
		return MetadataModuleStorageHashFuncIdentity, nil
	default:
		return MetadataModuleStorageHashFunc(0), newTypeError("MetadataModuleStorageHashFunc")
	}
}

func (mmshf MetadataModuleStorageHashFunc) Bytes() []byte {
	return sc.EncodedBytes(mmshf)
}
