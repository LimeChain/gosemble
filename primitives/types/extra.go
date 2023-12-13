package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/iancoleman/strcase"
)

const (
	additionalSignedTypeName = "typesInfoAdditionalSignedData"
	moduleTypeName           = "Module"
	hookOnChargeTypeName     = "OnChargeTransaction"
)

type SignedExtra interface {
	sc.Encodable

	Decode(buffer *bytes.Buffer)

	AdditionalSigned() (AdditionalSigned, error)
	Validate(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error)
	ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error)
	PreDispatch(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (sc.Sequence[Pre], error)
	PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error
	PostDispatch(pre sc.Option[sc.Sequence[Pre]], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) error

	Metadata(metadataIds map[string]int) (sc.Sequence[MetadataType], sc.Sequence[MetadataSignedExtension])
}

// signedExtra contains an array of SignedExtension, iterated through during extrinsic execution.
type signedExtra struct {
	extras []SignedExtension
}

func NewSignedExtra(checks []SignedExtension) SignedExtra {
	return signedExtra{
		extras: checks,
	}
}

func (e signedExtra) Encode(buffer *bytes.Buffer) error {
	for _, extra := range e.extras {
		err := extra.Encode(buffer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e signedExtra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e signedExtra) Decode(buffer *bytes.Buffer) {
	for _, extra := range e.extras {
		extra.Decode(buffer)
	}
}

func (e signedExtra) AdditionalSigned() (AdditionalSigned, error) {
	result := AdditionalSigned{}

	for _, extra := range e.extras {
		ok, err := extra.AdditionalSigned()
		if err != nil {
			return nil, err
		}
		result = append(result, ok...)
	}

	return result, nil
}

func (e signedExtra) Validate(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	valid := DefaultValidTransaction()

	for _, extra := range e.extras {
		v, err := extra.Validate(who, call, info, length)
		if err != nil {
			return ValidTransaction{}, err
		}
		valid = valid.CombineWith(v)
	}

	return valid, nil
}

func (e signedExtra) ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error) {
	valid := DefaultValidTransaction()

	for _, extra := range e.extras {
		v, err := extra.ValidateUnsigned(call, info, length)
		if err != nil {
			return ValidTransaction{}, err
		}
		valid = valid.CombineWith(v)
	}

	return valid, nil
}

func (e signedExtra) PreDispatch(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (sc.Sequence[Pre], error) {
	pre := sc.Sequence[Pre]{}

	for _, extra := range e.extras {
		p, err := extra.PreDispatch(who, call, info, length)
		if err != nil {
			return nil, err
		}

		pre = append(pre, p)
	}

	return pre, nil
}

func (e signedExtra) PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error {
	for _, extra := range e.extras {
		err := extra.PreDispatchUnsigned(call, info, length)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e signedExtra) PostDispatch(pre sc.Option[sc.Sequence[Pre]], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) error {
	if pre.HasValue {
		preValue := pre.Value
		for i, extra := range e.extras {
			err := extra.PostDispatch(sc.NewOption[Pre](preValue[i]), info, postInfo, length, result)
			if err != nil {
				return err
			}
		}
	} else {
		for _, extra := range e.extras {
			err := extra.PostDispatch(sc.NewOption[Pre](nil), info, postInfo, length, result)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e signedExtra) Metadata(metadataIds map[string]int) (sc.Sequence[MetadataType], sc.Sequence[MetadataSignedExtension]) {
	ids := sc.Sequence[sc.Compact]{}
	extraTypes := sc.Sequence[MetadataType]{}
	signedExtensions := sc.Sequence[MetadataSignedExtension]{}

	for _, extra := range e.extras {
		extraMetadataId := generateExtraMetadata(extra, metadataIds, &extraTypes, &signedExtensions)
		ids = append(ids, sc.ToCompact(extraMetadataId))
	}

	signedExtraType := NewMetadataType(metadata.SignedExtra, "SignedExtra", NewMetadataTypeDefinitionTuple(ids))

	return append(extraTypes, signedExtraType), signedExtensions
}

// generateExtraMetadata generates the metadata for a signed extension. It may generate some new metadata types in the process. Returns the metadata id for the extra
func generateExtraMetadata(extra SignedExtension, metadataIds map[string]int, metadataTypes *sc.Sequence[MetadataType], extensions *sc.Sequence[MetadataSignedExtension]) int {
	extraValue := reflect.ValueOf(extra)
	extraType := extraValue.Elem().Type()
	extraTypeName := extraType.Name()

	extraMetadataId, ok := metadataIds[extraTypeName]
	if !ok {
		extraMetadataId = buildMetadataTypeRecursively(extraType, metadataIds, metadataTypes, true, extra.ModulePath())
	} else {
		appendExtraMetadata(extraMetadataId, metadataTypes)
	}

	constructExtension(extraValue, extraMetadataId, extensions, metadataIds, metadataTypes, extra.ModulePath())

	return extraMetadataId
}

func appendExtraMetadata(id int, metadataTypes *sc.Sequence[MetadataType]) {
	for _, t := range *metadataTypes {
		if t.Id.ToBigInt().Int64() == int64(id) {
			*metadataTypes = append(*metadataTypes, t)
			break
		}
	}
}

// buildMetadataTypeRecursively build the metadata of the type recursively.
func buildMetadataTypeRecursively(t reflect.Type, metadataIds map[string]int, metadataTypes *sc.Sequence[MetadataType], isExtra bool, modulePath string) int {
	typeId := assignNewMetadataId(metadataIds, t.Name())

	typeName := t.Name()

	metadataFields := sc.Sequence[MetadataTypeDefinitionField]{}

	typeNumFields := 0

	if t.Kind() == reflect.Struct {
		typeNumFields = t.NumField()
	}

	for i := 0; i < typeNumFields; i++ {
		fieldName := t.Field(i).Name
		fieldTypeName := t.Field(i).Type.Name()
		if isIgnoredName(fieldName) || isIgnoredType(fieldTypeName) {
			continue
		}
		fieldId, ok := metadataIds[fieldTypeName]
		if !ok {
			fieldId = buildMetadataTypeRecursively(t.Field(i).Type, metadataIds, metadataTypes, false, modulePath)
		}
		metadataFields = append(metadataFields, NewMetadataTypeDefinitionFieldWithName(fieldId, sc.Str(fieldName)))
	}

	if isExtra {
		metadataType := NewMetadataTypeWithPath(typeId, typeName, sc.Sequence[sc.Str]{sc.Str(modulePath), "extensions", sc.Str(strcase.ToSnake(typeName)), sc.Str(typeName)}, NewMetadataTypeDefinitionComposite(metadataFields))
		*metadataTypes = append(*metadataTypes, metadataType)
		return typeId
	}

	*metadataTypes = append(*metadataTypes,
		NewMetadataType(
			typeId,
			typeName,
			NewMetadataTypeDefinitionComposite(metadataFields)))

	return typeId
}

func generateCompositeType(typeId int, typeName string, tupleIds sc.Sequence[sc.Compact]) MetadataType {
	return NewMetadataType(typeId, typeName, NewMetadataTypeDefinitionTuple(tupleIds))
}

func isIgnoredType(t string) bool {
	return t == moduleTypeName || t == hookOnChargeTypeName
}

func isIgnoredName(name string) bool {
	return name == additionalSignedTypeName
}

func assignNewMetadataId(metadataIds map[string]int, name string) int {
	lastIndex := len(metadataIds)
	newId := lastIndex + 1
	metadataIds[name] = newId
	return newId
}

// constructExtension Iterates through the elements of the typesInfoAdditionalSignedData slice and builds the extra extension. If an element in the slice is a type not present in the metadata map, it will also be generated.
func constructExtension(extra reflect.Value, extraMetadataId int, extensions *sc.Sequence[MetadataSignedExtension], metadataIds map[string]int, metadataTypes *sc.Sequence[MetadataType], modulePath string) {
	var resultTypeName string
	var resultTupleIds sc.Sequence[sc.Compact]

	extraType := extra.Elem().Type
	extraName := extraType().Name()

	additionalSignedField := extra.Elem().FieldByName(additionalSignedTypeName)

	if additionalSignedField.IsValid() {
		numAdditionalSignedTypes := additionalSignedField.Len()
		if numAdditionalSignedTypes == 0 {
			*extensions = append(*extensions, NewMetadataSignedExtension(sc.Str(extraName), extraMetadataId, metadata.TypesEmptyTuple))
			return
		}
		for i := 0; i < numAdditionalSignedTypes; i++ {
			currentType := additionalSignedField.Index(i).Elem().Type()
			currentTypeName := currentType.Name()
			currentTypeId, ok := metadataIds[currentTypeName]
			if !ok {
				currentTypeId = buildMetadataTypeRecursively(currentType, metadataIds, metadataTypes, false, modulePath)
			}
			resultTypeName = resultTypeName + currentTypeName
			resultTupleIds = append(resultTupleIds, sc.ToCompact(currentTypeId))
		}
		resultTypeId, ok := metadataIds[resultTypeName]
		if !ok {
			resultTypeId = assignNewMetadataId(metadataIds, resultTypeName)
			*metadataTypes = append(*metadataTypes, generateCompositeType(resultTypeId, resultTypeName, resultTupleIds))
		}
		*extensions = append(*extensions, NewMetadataSignedExtension(sc.Str(extraName), extraMetadataId, resultTypeId))
	}
}
