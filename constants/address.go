package constants

import primitives "github.com/LimeChain/gosemble/primitives/types"

var (
	ZeroAddressAccountId = primitives.AccountId{Ed25519Signer: primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)}
	OneAddressAccountId  = primitives.AccountId{Ed25519Signer: primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)}
	TwoAddressAccountId  = primitives.AccountId{Ed25519Signer: primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2)}
)
