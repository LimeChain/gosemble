package types

import (
	sc "github.com/LimeChain/goscale"
)

type PerDispatchClassWeight struct {
	// Value for `Normal` extrinsics.
	Normal Weight
	// Value for `Operational` extrinsics.
	Operational Weight
	// Value for `Mandatory` extrinsics.
	Mandatory Weight
}

func (pw PerDispatchClassWeight) GetMetadata(typeId int, generator *MetadataTypeGenerator) MetadataType {
	typesWeightId, _ := generator.GetId("Weight")

	return NewMetadataTypeWithPath(typeId, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "normal"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "operational"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "mandatory"),
		},
	),
	)
}
