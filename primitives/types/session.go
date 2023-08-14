package types

// Session provides the key type and id of a module, which has a session.
type Session interface {
	KeyType() PublicKeyType
	KeyTypeId() [4]byte
}
