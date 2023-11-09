package constants

import primitives "github.com/LimeChain/gosemble/primitives/types"

var (
	ed25519SignerZero, _ = primitives.NewEd25519PublicKey(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	ed25519SignerOne, _  = primitives.NewEd25519PublicKey(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)
	ed25519SignerTwo, _  = primitives.NewEd25519PublicKey(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2)
	ZeroAddressAccountId = primitives.New[primitives.PublicKey](ed25519SignerZero)
	OneAddressAccountId  = primitives.New[primitives.PublicKey](ed25519SignerOne)
	TwoAddressAccountId  = primitives.New[primitives.PublicKey](ed25519SignerTwo)
)
