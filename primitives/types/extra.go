package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/iancoleman/strcase"
)

const (
	additionalSignedTypeName = "additionalSignedData"
	moduleTypeName           = "Module"
	varyingDataType          = "VaryingData"
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

	Metadata(constantsIdsMap map[string]int) (sc.Sequence[MetadataType], sc.Sequence[MetadataSignedExtension])
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

func (e signedExtra) Metadata(constantsIdsMap map[string]int) (sc.Sequence[MetadataType], sc.Sequence[MetadataSignedExtension]) {
	ids := sc.Sequence[sc.Compact]{}
	extraTypes := sc.Sequence[MetadataType]{}
	signedExtensions := sc.Sequence[MetadataSignedExtension]{}

	for _, extra := range e.extras {
		extraType, extension := generateExtensionMetadata(extra, constantsIdsMap, &extraTypes)
		ids = append(ids, extraType.Id)
		extraTypes = append(extraTypes, extraType)
		signedExtensions = append(signedExtensions, extension)
	}

	signedExtraType := NewMetadataType(metadata.SignedExtra, "SignedExtra", NewMetadataTypeDefinitionTuple(ids))

	return append(extraTypes, signedExtraType), signedExtensions
}

// generateExtensionMetadata generates the metadata for a signed extension. It may generate some new metadata types in the process.
func generateExtensionMetadata(extra SignedExtension, metadataIds map[string]int, metadataTypes *sc.Sequence[MetadataType]) (MetadataType, MetadataSignedExtension) {
	extraValue := reflect.ValueOf(extra)
	extraType := extraValue.Elem().Type()
	extraTypeName := extraType.Name()

	lastIndex := len(metadataIds)

	typeConstantId, exists := metadataIds[extraTypeName]
	if !exists {
		typeConstantId = assignNewMetadataId(metadataIds, &lastIndex, &extraTypeName)
	}

	var extension MetadataSignedExtension
	var metadataTypeFields sc.Sequence[MetadataTypeDefinitionField]

	extraNumOfFields := extraType.NumField()

	for j := 0; j < extraNumOfFields; j++ {
		field := extraValue.Elem().Field(j)
		fieldTypeName := field.Type().Name()
		fieldName := extraType.Field(j).Name
		if fieldTypeName == moduleTypeName {
			continue
		}
		if fieldName != additionalSignedTypeName {
			fieldTypeId, exists := metadataIds[fieldTypeName]
			if !exists {
				fieldTypeId = assignNewMetadataId(metadataIds, &lastIndex, &fieldTypeName)
				newType := generateNewType(fieldTypeId, field.Type())
				*metadataTypes = append(*metadataTypes, newType)
			}
			metadataTypeFields = append(metadataTypeFields, NewMetadataTypeDefinitionFieldWithName(fieldTypeId, sc.Str(fieldTypeName)))
			continue
		}
		var resultTypeId int
		var resultTypeName string
		var resultTupleIds sc.Sequence[sc.Compact]
		numAdditionalSignedTypes := field.Len()
		if numAdditionalSignedTypes == 0 {
			extension = NewMetadataSignedExtension(sc.Str(extraTypeName), typeConstantId, metadata.TypesEmptyTuple)
			continue
		}
		for i := 0; i < numAdditionalSignedTypes; i++ {
			currentType := field.Index(i).Elem().Type()
			currentTypeName := currentType.Name()
			currentTypeId, exists := metadataIds[currentTypeName]
			if !exists {
				currentTypeId = assignNewMetadataId(metadataIds, &lastIndex, &currentTypeName)
				newType := generateNewType(currentTypeId, currentType)
				*metadataTypes = append(*metadataTypes, newType)
			}
			resultTypeName = resultTypeName + currentTypeName
			resultTupleIds = append(resultTupleIds, sc.ToCompact(currentTypeId))
		}
		resultTypeId, exists = metadataIds[resultTypeName]
		if !exists {
			resultTypeId = assignNewMetadataId(metadataIds, &lastIndex, &resultTypeName)
			*metadataTypes = append(*metadataTypes, generateCompositeType(resultTypeId, resultTypeName, resultTupleIds))
		}
		extension = NewMetadataSignedExtension(sc.Str(extraTypeName), typeConstantId, resultTypeId)
	}

	return NewMetadataTypeWithPath(
		typeConstantId,
		extraTypeName,
		sc.Sequence[sc.Str]{"extensions", sc.Str(strcase.ToSnake(extraTypeName)), sc.Str(extraTypeName)},
		NewMetadataTypeDefinitionComposite(metadataTypeFields)), extension
}

func generateNewType(id int, t reflect.Type) MetadataType {
	typeFields := sc.Sequence[MetadataTypeDefinitionField]{}

	typeName := t.Name()

	typeNumFields := t.NumField()

	for i := 0; i < typeNumFields; i++ {
		fieldName := t.Field(i).Name
		typeFields = append(typeFields, NewMetadataTypeDefinitionFieldWithName(id, sc.Str(fieldName)))
	}

	return NewMetadataType(
		id,
		typeName,
		NewMetadataTypeDefinitionComposite(typeFields),
	)
}

func generateCompositeType(typeId int, typeName string, tupleIds sc.Sequence[sc.Compact]) MetadataType {
	return NewMetadataType(typeId, typeName, NewMetadataTypeDefinitionTuple(tupleIds))
}

func assignNewMetadataId(metadataIds map[string]int, lastIndex *int, name *string) int {
	newId := *lastIndex + 1
	*lastIndex = *lastIndex + 1
	metadataIds[*name] = newId
	return newId
}
