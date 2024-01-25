package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type PerDispatchClassWeightsPerClass struct {
	// Value for `Normal` extrinsics.
	Normal WeightsPerClass
	// Value for `Operational` extrinsics.
	Operational WeightsPerClass
	// Value for `Mandatory` extrinsics.
	Mandatory WeightsPerClass
}

func (pdc PerDispatchClassWeightsPerClass) GetMetadata(typeId int, generator *MetadataTypeGenerator) MetadataType {
	typesWeightsPerClassId, _ := generator.GetId("WeightsPerClass")

	return NewMetadataTypeWithPath(typeId, "PerDispatchClass[WeightPerClass]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(typesWeightsPerClassId, "normal"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightsPerClassId, "operational"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightsPerClassId, "mandatory"),
		}))

	//return NewMetadataTypeWithParam(typeId, "PerDispatchClass[WeightPerClass]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, NewMetadataTypeDefinitionComposite(
	//	sc.Sequence[MetadataTypeDefinitionField]{
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightsPerClassId, "normal", "T"),
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightsPerClassId, "operational", "T"),
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightsPerClassId, "mandatory", "T"),
	//	}),
	//	NewMetadataTypeParameter(typesWeightsPerClassId, "T"))
}

func (pdc PerDispatchClassWeightsPerClass) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		pdc.Normal,
		pdc.Operational,
		pdc.Mandatory,
	)
}

func DecodePerDispatchClassWeightPerClass(buffer *bytes.Buffer, decodeWeightPerClass func(buffer *bytes.Buffer) (WeightsPerClass, error)) (PerDispatchClassWeightsPerClass, error) {
	normal, err := decodeWeightPerClass(buffer)
	if err != nil {
		return PerDispatchClassWeightsPerClass{}, err
	}
	operational, err := decodeWeightPerClass(buffer)
	if err != nil {
		return PerDispatchClassWeightsPerClass{}, err
	}
	mandatory, err := decodeWeightPerClass(buffer)
	if err != nil {
		return PerDispatchClassWeightsPerClass{}, err
	}
	return PerDispatchClassWeightsPerClass{
		Normal:      normal,
		Operational: operational,
		Mandatory:   mandatory,
	}, nil
}

func (pdc PerDispatchClassWeightsPerClass) Bytes() []byte {
	return sc.EncodedBytes(pdc)
}

// Get current value for given class.
func (pdc *PerDispatchClassWeightsPerClass) Get(class DispatchClass) (*WeightsPerClass, error) {
	switch class.VaryingData[0] {
	case DispatchClassNormal:
		return &pdc.Normal, nil
	case DispatchClassOperational:
		return &pdc.Operational, nil
	case DispatchClassMandatory:
		return &pdc.Mandatory, nil
	default:
		return nil, newTypeError("DispatchClass")
	}
}
