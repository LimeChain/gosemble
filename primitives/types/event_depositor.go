package types

type EventDepositor interface {
	DepositEvent(event Event)
}
