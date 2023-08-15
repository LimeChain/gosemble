package types

type ApiModule interface {
	Name() string
	Item() ApiItem
}
