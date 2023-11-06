package extensions

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type CheckNonZeroAddress struct{}

func NewCheckNonZeroAddress() CheckNonZeroAddress {
	return CheckNonZeroAddress{}
}

func (c CheckNonZeroAddress) AdditionalSigned() (primitives.AdditionalSigned, primitives.TransactionValidityError) {
	return primitives.AdditionalSigned{}, nil
}

func (c CheckNonZeroAddress) Encode(*bytes.Buffer) error {
	return nil
}

func (c CheckNonZeroAddress) Decode(*bytes.Buffer) error { return nil }

func (c CheckNonZeroAddress) Bytes() []byte {
	return sc.EncodedBytes(c)
}

func (c CheckNonZeroAddress) Validate(who primitives.AccountId, _call primitives.Call, _info *primitives.DispatchInfo, _length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	// TODO:
	// Not sure when this is possible.
	// Checks signed transactions but will fail
	// before this check if the address is all zeros.
	if reflect.DeepEqual(who, constants.ZeroAddressAccountId) {
		// TODO https://github.com/LimeChain/gosemble/issues/271
		invalidTransactionBadSigner, _ := primitives.NewTransactionValidityError(primitives.NewInvalidTransactionBadSigner())
		return primitives.ValidTransaction{}, invalidTransactionBadSigner
	}

	return primitives.DefaultValidTransaction(), nil
}

func (c CheckNonZeroAddress) ValidateUnsigned(_call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (c CheckNonZeroAddress) PreDispatch(who primitives.AccountId, call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) (primitives.Pre, primitives.TransactionValidityError) {
	_, err := c.Validate(who, call, info, length)
	return primitives.Pre{}, err
}

func (c CheckNonZeroAddress) PreDispatchUnsigned(call primitives.Call, info *primitives.DispatchInfo, length sc.Compact) primitives.TransactionValidityError {
	_, err := c.ValidateUnsigned(call, info, length)
	return err
}

func (c CheckNonZeroAddress) PostDispatch(_pre sc.Option[primitives.Pre], info *primitives.DispatchInfo, postInfo *primitives.PostDispatchInfo, _length sc.Compact, _result *primitives.DispatchResult) primitives.TransactionValidityError {
	return nil
}

func (c CheckNonZeroAddress) Metadata() (primitives.MetadataType, primitives.MetadataSignedExtension) {
	return primitives.NewMetadataTypeWithPath(
			metadata.CheckNonZeroSender,
			"CheckNonZeroSender",
			sc.Sequence[sc.Str]{"frame_system", "extensions", "check_non_zero_sender", "CheckNonZeroSender"},
			primitives.NewMetadataTypeDefinitionComposite(sc.Sequence[primitives.MetadataTypeDefinitionField]{}),
		),
		primitives.NewMetadataSignedExtension("CheckNonZeroSender", metadata.CheckNonZeroSender, metadata.TypesEmptyTuple)
}
