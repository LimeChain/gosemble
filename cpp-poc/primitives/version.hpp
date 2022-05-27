#include <stdint.h>

namespace primitives
{
	struct ApiItem
	{
		int8_t Name[8];
		uint32_t Version;
	};

	struct VersionData
	{
		// int8_t SpecName[];
		// int8_t ImplName[];
		uint32_t AuthoringVersion;
		uint32_t SpecVersion;
		uint32_t ImplVersion;
		// ApiItem Apis[];
		uint32_t TransactionVersion;
		uint32_t StateVersion;
	};
};

// func(v *VersionData) Encode()([] byte, error)
// {
// 	enc, err : = scale.Encode(v) if err != nil
// 	{
// 		return nil, err
// 	}

// 	return enc, nil
// }

// func(v *VersionData) Decode(enc[] byte) error
// {
// 	var data VersionData

// 			_,
// 			err : = scale.Decode(enc, &data) if err != nil{
// 																										 return err}

// 																								 v.SpecName = data.SpecName v.ImplName = data.ImplName v.AuthoringVersion = data.AuthoringVersion v.SpecVersion = data.SpecVersion v.ImplVersion = data.ImplVersion v.Apis = data.Apis v.TransactionVersion = data.TransactionVersion v.StateVersion = data.StateVersion

// 																																																																																																																																																			 return nil
// };
