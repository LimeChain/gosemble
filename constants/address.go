package constants

import primitives "github.com/LimeChain/gosemble/primitives/types"

var (
	ed25519SignerZero, _ = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	ed25519SignerOne, _  = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)
	ed25519SignerTwo, _  = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2)
	ZeroAddressAccountId = primitives.AccountId{Ed25519Signer: ed25519SignerZero}
	OneAddressAccountId  = primitives.AccountId{Ed25519Signer: ed25519SignerOne}
	TwoAddressAccountId  = primitives.AccountId{Ed25519Signer: ed25519SignerTwo}
)
