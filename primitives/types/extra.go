package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/iancoleman/strcase"
)

const (
	additionalSignedTypeName = "VaryingData"
	moduleName               = "Module"
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
		metadataType, extension := constructExtensionMetadata(extra, constantsIdsMap)
		ids = append(ids, metadataType.Id)
		extraTypes = append(extraTypes, metadataType)
		signedExtensions = append(signedExtensions, extension)
	}

	signedExtraType := NewMetadataType(metadata.SignedExtra, "SignedExtra", NewMetadataTypeDefinitionTuple(ids))

	return append(extraTypes, signedExtraType), signedExtensions
}

func constructExtensionMetadata(extra SignedExtension, constantsMap map[string]int) (MetadataType, MetadataSignedExtension) {
	extraValue := reflect.ValueOf(extra)
	extraType := extraValue.Elem().Type()
	extraTypeName := extraType.Name()

	lastIndex := len(constantsMap)

	typeConstantId, exists := constantsMap[extraTypeName]
	if !exists {
		typeConstantId = lastIndex + 1
		lastIndex = lastIndex + 1
		constantsMap[extraTypeName] = typeConstantId
	}

	var extension MetadataSignedExtension
	var metadataTypeFields = sc.Sequence[MetadataTypeDefinitionField]{}

	typeNumOfFields := extraValue.Elem().NumField()

	for j := 0; j < typeNumOfFields; j++ {
		field := extraValue.Elem().Field(j)
		fieldName := field.Type().Name()
		// log.Info("Name: " + fieldName)
		switch fieldName {
		case moduleName, "Bool":
			continue
		case additionalSignedTypeName: // Process additionalSigned type(s) which define the MetadataSignedExtension
			var additionalSignedTypeId int
			numAdditionalSignedTypes := field.Len()
			for i := 0; i < numAdditionalSignedTypes; i++ { // Note: Making it work for only 1 item for now TODO: if varying type has more than 1 element, make a complex type
				additionalSignedType := field.Index(i).Elem()
				additionalSignedName := additionalSignedType.Type().Name() // e.g. Bool, U32, etc
				additionalSignedTypeId, exists = constantsMap[additionalSignedName]
				if !exists {
					additionalSignedTypeId = lastIndex + 1
					constantsMap[additionalSignedName] = additionalSignedTypeId
					log.Info("Addine new: " + additionalSignedName)
					lastIndex = lastIndex + 1
				}
			}
			if numAdditionalSignedTypes == 0 {
				extension = NewMetadataSignedExtension(sc.Str(extraTypeName), typeConstantId, metadata.TypesEmptyTuple)
			} else {
				extension = NewMetadataSignedExtension(sc.Str(extraTypeName), typeConstantId, additionalSignedTypeId)
			}
			continue // We have determined the additionalTypeId for the extension
		}
		fieldTypeId, exists := constantsMap[fieldName]
		if !exists {
			fieldTypeId = lastIndex + 1
			constantsMap[fieldName] = fieldTypeId
			log.Info("Adding type id new: " + fieldName)
			lastIndex = lastIndex + 1
		}
		metadataTypeFields = append(metadataTypeFields, NewMetadataTypeDefinitionFieldWithName(fieldTypeId, sc.Str(fieldName)))
	}

	return NewMetadataTypeWithPath(
		typeConstantId,
		extraTypeName,
		sc.Sequence[sc.Str]{"extensions", sc.Str(strcase.ToSnake(extraTypeName)), sc.Str(extraTypeName)},
		NewMetadataTypeDefinitionComposite(metadataTypeFields),
	), extension
}
