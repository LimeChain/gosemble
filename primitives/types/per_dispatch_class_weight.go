package types

import (
	"strconv"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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
	log.NewLogger().Info("Weight in PerDispatchClass: " + strconv.Itoa(typesWeightId))

	return NewMetadataTypeWithPath(typeId, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, NewMetadataTypeDefinitionComposite(
		sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "normal"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "operational"),
			NewMetadataTypeDefinitionFieldWithName(typesWeightId, "mandatory"),
		},
	),
	)

	//return NewMetadataTypeWithParam(typeId, "PerDispatchClass[Weight]", sc.Sequence[sc.Str]{"frame_support", "dispatch", "PerDispatchClass"}, NewMetadataTypeDefinitionComposite(
	//	sc.Sequence[MetadataTypeDefinitionField]{
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightId, "normal", "T"),
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightId, "operational", "T"),
	//		NewMetadataTypeDefinitionFieldWithNames(typesWeightId, "mandatory", "T"),
	//	},
	//),
	//	NewMetadataTypeParameter(typesWeightId, "T"),
	//)
}
