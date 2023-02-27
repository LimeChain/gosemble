package constants

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

// If the runtime behavior changes, increment spec_version and set impl_version to 0.
// If only runtime implementation changes and behavior does not,
// then leave spec_version as is and increment impl_version.

const SpecName = "node-template"
const ImplName = "node-template"
const AuthoringVersion = 1
const SpecVersion = 100
const ImplVersion = 1
const TransactionVersion = 1
const StateVersion = 1
const StorageVersion = 0

const BlockHashCount = sc.U32(2400)

var RuntimeVersion = types.RuntimeVersion{
	SpecName:         sc.Str(SpecName),
	ImplName:         sc.Str(ImplName),
	AuthoringVersion: sc.U32(AuthoringVersion),
	SpecVersion:      sc.U32(SpecVersion),
	ImplVersion:      sc.U32(ImplVersion),
	Apis: sc.Sequence[types.ApiItem]{
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 223, 106, 203, 104, 153, 7, 96, 155),
			Version: sc.U32(3),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 55, 227, 151, 252, 124, 145, 245, 228),
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 64, 254, 58, 212, 1, 248, 149, 154),
			Version: sc.U32(4),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 210, 188, 152, 151, 238, 208, 143, 21),
			Version: sc.U32(2),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 247, 139, 39, 139, 229, 63, 69, 76),
			Version: sc.U32(2),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 237, 153, 197, 172, 178, 94, 237, 245),
			Version: sc.U32(2),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 203, 202, 37, 227, 159, 20, 35, 135),
			Version: sc.U32(2),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 104, 122, 212, 74, 211, 127, 3, 194),
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 188, 157, 137, 144, 79, 91, 146, 63),
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 104, 182, 107, 161, 34, 201, 63, 167),
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 55, 200, 187, 19, 80, 169, 162, 168),
			Version: sc.U32(1),
		},
		{
			Name:    sc.NewFixedSequence[sc.U8](8, 171, 60, 5, 114, 41, 31, 235, 139),
			Version: sc.U32(1),
		},
	},
	TransactionVersion: sc.U32(TransactionVersion),
	StateVersion:       sc.U8(StateVersion),
}
