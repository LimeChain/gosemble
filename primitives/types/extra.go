package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/iancoleman/strcase"
)

type SignedExtra interface {
	sc.Encodable

	Decode(buffer *bytes.Buffer)

	AdditionalSigned() (AdditionalSigned, error)
	Validate(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error)
	ValidateUnsigned(call Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, error)
	PreDispatch(who AccountId, call Call, info *DispatchInfo, length sc.Compact) (sc.Sequence[Pre], error)
	PreDispatchUnsigned(call Call, info *DispatchInfo, length sc.Compact) error
	PostDispatch(pre sc.Option[sc.Sequence[Pre]], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, dispatchErr error) error

	Metadata() sc.Sequence[MetadataSignedExtension]
}

// signedExtra contains an array of SignedExtension, iterated through during extrinsic execution.
type signedExtra struct {
	extras      []SignedExtension
	mdGenerator *MetadataTypeGenerator
}

func NewSignedExtra(checks []SignedExtension, mdGenerator *MetadataTypeGenerator) SignedExtra {
	return signedExtra{
		extras:      checks,
		mdGenerator: mdGenerator,
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

func (e signedExtra) PostDispatch(pre sc.Option[sc.Sequence[Pre]], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, dispatchErr error) error {
	if pre.HasValue {
		preValue := pre.Value
		for i, extra := range e.extras {
			err := extra.PostDispatch(sc.NewOption[Pre](preValue[i]), info, postInfo, length, dispatchErr)
			if err != nil {
				return err
			}
		}
	} else {
		for _, extra := range e.extras {
			err := extra.PostDispatch(sc.NewOption[Pre](nil), info, postInfo, length, dispatchErr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e signedExtra) Metadata() sc.Sequence[MetadataSignedExtension] {
	ids := sc.Sequence[sc.Compact]{}
	extraTypes := sc.Sequence[MetadataType]{}
	signedExtensions := sc.Sequence[MetadataSignedExtension]{}

	for _, extra := range e.extras {
		extraMetadataId := generateExtraMetadata(extra, e.mdGenerator, &extraTypes, &signedExtensions)
		ids = append(ids, sc.ToCompact(extraMetadataId))
	}

	e.mdGenerator.AppendMetadataTypes(extraTypes)

	signedExtraType := NewMetadataType(metadata.SignedExtra, "SignedExtra", NewMetadataTypeDefinitionTuple(ids))
	e.mdGenerator.AppendMetadataTypes(sc.Sequence[MetadataType]{signedExtraType})

	return signedExtensions
}

// generateExtraMetadata generates the metadata for a signed extension. It may generate some new metadata types in the process. Returns the metadata id for the extra
func generateExtraMetadata(extra SignedExtension, metadataGenerator *MetadataTypeGenerator, metadataTypes *sc.Sequence[MetadataType], extensions *sc.Sequence[MetadataSignedExtension]) int {
	extraValue := reflect.ValueOf(extra)
	extraTypeName := extraValue.Elem().Type().Name()
	extraMetadataId := metadataGenerator.BuildMetadataTypeRecursively(extraValue.Elem(), &sc.Sequence[sc.Str]{sc.Str(extra.ModulePath()), "extensions", sc.Str(strcase.ToSnake(extraTypeName)), sc.Str(extraTypeName)}, nil, nil)
	constructExtension(extraValue, extraMetadataId, extensions, metadataGenerator, metadataTypes)

	return extraMetadataId
}

func generateCompositeType(typeId int, typeName string, tupleIds sc.Sequence[sc.Compact]) MetadataType {
	return NewMetadataType(typeId, typeName, NewMetadataTypeDefinitionTuple(tupleIds))
}

// constructExtension Iterates through the elements of the typesInfoAdditionalSignedData slice and builds the extra extension. If an element in the slice is a type not present in the metadata map, it will also be generated.
func constructExtension(extra reflect.Value, extraMetadataId int, extensions *sc.Sequence[MetadataSignedExtension], metadataGenerator *MetadataTypeGenerator, metadataTypes *sc.Sequence[MetadataType]) {
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
			currentType := additionalSignedField.Index(i).Elem()
			currentTypeName := currentType.Type().Name()
			currentTypeId, ok := metadataGenerator.GetId(currentTypeName)
			if !ok {
				currentTypeId = metadataGenerator.BuildMetadataTypeRecursively(currentType, nil, nil, nil)
			}
			resultTypeName = resultTypeName + currentTypeName
			resultTupleIds = append(resultTupleIds, sc.ToCompact(currentTypeId))
		}
		resultTypeId, ok := metadataGenerator.GetId(resultTypeName)
		if !ok {
			resultTypeId = metadataGenerator.assignNewMetadataId(resultTypeName)
			*metadataTypes = append(*metadataTypes, generateCompositeType(resultTypeId, resultTypeName, resultTupleIds))
		}
		*extensions = append(*extensions, NewMetadataSignedExtension(sc.Str(extraName), extraMetadataId, resultTypeId))
	}
}
