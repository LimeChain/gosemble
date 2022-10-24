package constants

import "github.com/LimeChain/gosemble/types"

const SPEC_NAME = "gosemble"
const IMPL_NAME = "go"
const AUTHORING_VERSION = 1
const SPEC_VERSION = 2
const IMPL_VERSION = 3
const TRANSACTION_VERSION = 4
const STATE_VERSION = 5

var VersionDataConfig = types.VersionData{
	SpecName:         []byte(SPEC_NAME),
	ImplName:         []byte(IMPL_NAME),
	AuthoringVersion: uint32(AUTHORING_VERSION),
	SpecVersion:      uint32(SPEC_VERSION),
	ImplVersion:      uint32(IMPL_VERSION),
	Apis: []types.ApiItem{
		{Name: [8]byte{1, 1, 1, 1, 1, 1, 1, 1}, Version: 0},
	},
	TransactionVersion: uint32(TRANSACTION_VERSION),
	StateVersion:       uint32(STATE_VERSION),
}
