package constants

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/types"
)

// If the runtime behavior changes, increment spec_version and set impl_version to 0.
// If only runtime implementation changes and behavior does not,
// then leave spec_version as is and increment impl_version.

const SPEC_NAME = "node-template"
const IMPL_NAME = "node-template"
const AUTHORING_VERSION = 1
const SPEC_VERSION = 100
const IMPL_VERSION = 1
const TRANSACTION_VERSION = 1
const STATE_VERSION = 1

var RuntimeVersion = types.VersionData{
	SpecName:         sc.Str(SPEC_NAME),
	ImplName:         sc.Str(IMPL_NAME),
	AuthoringVersion: sc.U32(AUTHORING_VERSION),
	SpecVersion:      sc.U32(SPEC_VERSION),
	ImplVersion:      sc.U32(IMPL_VERSION),
	Apis: sc.Sequence[types.ApiItem]{
		{
			Name:    sc.FixedSequence[sc.U8]{223, 106, 203, 104, 153, 7, 96, 155},
			Version: sc.U32(3),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{55, 227, 151, 252, 124, 145, 245, 228},
			Version: sc.U32(1),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{64, 254, 58, 212, 1, 248, 149, 154},
			Version: sc.U32(4),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{210, 188, 152, 151, 238, 208, 143, 21},
			Version: sc.U32(2),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{247, 139, 39, 139, 229, 63, 69, 76},
			Version: sc.U32(2),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{237, 153, 197, 172, 178, 94, 237, 245},
			Version: sc.U32(2),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{203, 202, 37, 227, 159, 20, 35, 135},
			Version: sc.U32(2),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{104, 122, 212, 74, 211, 127, 3, 194},
			Version: sc.U32(1),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{188, 157, 137, 144, 79, 91, 146, 63},
			Version: sc.U32(1),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{104, 182, 107, 161, 34, 201, 63, 167},
			Version: sc.U32(1),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{55, 200, 187, 19, 80, 169, 162, 168},
			Version: sc.U32(1),
		},
		{
			Name:    sc.FixedSequence[sc.U8]{171, 60, 5, 114, 41, 31, 235, 139},
			Version: sc.U32(1),
		},
	},
	TransactionVersion: sc.U32(TRANSACTION_VERSION),
	StateVersion:       sc.U8(STATE_VERSION),
}
