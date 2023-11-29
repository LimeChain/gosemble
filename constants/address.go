package constants

import primitives "github.com/LimeChain/gosemble/primitives/types"

var (
	ZeroAddress, _ = primitives.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	OneAddress, _  = primitives.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)
	TwoAddress, _  = primitives.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2)

	ZeroAccountId = primitives.NewAccountIdFromAddress32(ZeroAddress)
	OneAccountId  = primitives.NewAccountIdFromAddress32(OneAddress)
	TwoAccountId  = primitives.NewAccountIdFromAddress32(TwoAddress)
)
