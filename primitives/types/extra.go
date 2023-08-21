package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

// SignedExtra contains an array of SignedExtension, iterated through during extrinsic execution.
type SignedExtra struct {
	extras []SignedExtension
}

func NewSignedExtra(checks []SignedExtension) SignedExtra {
	return SignedExtra{
		extras: checks,
	}
}

func (e SignedExtra) Encode(buffer *bytes.Buffer) {
	for _, extra := range e.extras {
		extra.Encode(buffer)
	}
}

func (e *SignedExtra) Decode(buffer *bytes.Buffer) {
	for _, extra := range e.extras {
		extra.Decode(buffer)
	}
}

func (e SignedExtra) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e SignedExtra) AdditionalSigned() (AdditionalSigned, TransactionValidityError) {
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

func (e SignedExtra) Validate(who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
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

func (e SignedExtra) ValidateUnsigned(call *Call, info *DispatchInfo, length sc.Compact) (ValidTransaction, TransactionValidityError) {
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

func (e SignedExtra) PreDispatch(who *Address32, call *Call, info *DispatchInfo, length sc.Compact) (sc.Sequence[Pre], TransactionValidityError) {
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

func (e SignedExtra) PreDispatchUnsigned(call *Call, info *DispatchInfo, length sc.Compact) TransactionValidityError {
	for _, extra := range e.extras {
		err := extra.PreDispatchUnsigned(call, info, length)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e SignedExtra) PostDispatch(pre sc.Option[sc.Sequence[Pre]], info *DispatchInfo, postInfo *PostDispatchInfo, length sc.Compact, result *DispatchResult) TransactionValidityError {
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

func (e SignedExtra) Metadata() (sc.Sequence[MetadataType], sc.Sequence[MetadataSignedExtension]) {
	ids := sc.Sequence[sc.Compact]{}
	extraTypes := sc.Sequence[MetadataType]{}
	signedExtensions := sc.Sequence[MetadataSignedExtension]{}

	for _, extra := range e.extras {
		metadataType, extension := extra.Metadata()

		ids = append(ids, metadataType.Id)
		extraTypes = append(extraTypes, metadataType)
		signedExtensions = append(signedExtensions, extension)
	}

	signedExtraType := NewMetadataType(metadata.SignedExtra, "SignedExtra", NewMetadataTypeDefinitionTuple(ids))

	return append(extraTypes, signedExtraType), signedExtensions
}
