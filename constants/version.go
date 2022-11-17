package constants

import "github.com/LimeChain/gosemble/types"

// Per convention: if the runtime behavior changes, increment spec_version
// and set impl_version to 0. If only runtime
// implementation changes and behavior does not, then leave spec_version as
// is and increment impl_version.

const SPEC_NAME = "polkadot"
const IMPL_NAME = "parity-polkadot"
const AUTHORING_VERSION = 0
const SPEC_VERSION = 9160
const IMPL_VERSION = 0
const TRANSACTION_VERSION = 0
const STATE_VERSION = 0

var RuntimeVersion = types.VersionData{
	SpecName:         []byte(SPEC_NAME),
	ImplName:         []byte(IMPL_NAME),
	AuthoringVersion: uint32(AUTHORING_VERSION),
	SpecVersion:      uint32(SPEC_VERSION),
	ImplVersion:      uint32(IMPL_VERSION),
	Apis: []types.ApiItem{
		{Name: [8]byte{223, 106, 203, 104, 153, 7, 96, 155}, Version: 3},
		{Name: [8]byte{55, 227, 151, 252, 124, 145, 245, 228}, Version: 1},
		{Name: [8]byte{64, 254, 58, 212, 1, 248, 149, 154}, Version: 4},
		{Name: [8]byte{210, 188, 152, 151, 238, 208, 143, 21}, Version: 2},
		{Name: [8]byte{247, 139, 39, 139, 229, 63, 69, 76}, Version: 2},
		{Name: [8]byte{237, 153, 197, 172, 178, 94, 237, 245}, Version: 2},
		{Name: [8]byte{203, 202, 37, 227, 159, 20, 35, 135}, Version: 2},
		{Name: [8]byte{104, 122, 212, 74, 211, 127, 3, 194}, Version: 1},
		{Name: [8]byte{188, 157, 137, 144, 79, 91, 146, 63}, Version: 1},
		{Name: [8]byte{104, 182, 107, 161, 34, 201, 63, 167}, Version: 1},
		{Name: [8]byte{55, 200, 187, 19, 80, 169, 162, 168}, Version: 1},
		{Name: [8]byte{171, 60, 5, 114, 41, 31, 235, 139}, Version: 1},
	},
	TransactionVersion: uint32(TRANSACTION_VERSION),
	StateVersion:       uint8(STATE_VERSION),
}
