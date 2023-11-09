package constants

import primitives "github.com/LimeChain/gosemble/primitives/types"

var (
	ed25519SignerZero, _ = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	ed25519SignerOne, _  = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)
	ed25519SignerTwo, _  = primitives.NewEd25519Signer(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2)
	ZeroAddressAccountId = primitives.New[primitives.SignerAddress](ed25519SignerZero)
	OneAddressAccountId  = primitives.New[primitives.SignerAddress](ed25519SignerOne)
	TwoAddressAccountId  = primitives.New[primitives.SignerAddress](ed25519SignerTwo)
)
