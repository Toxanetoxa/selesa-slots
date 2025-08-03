package wallet

import "time"

type Type string

const (
	EventDeposit    Type = "deposit"
	EventWithdrawal Type = "withdrawal"
)

type Event struct {
	Type      Type      `json:"type"`
	UserID    UserID    `json:"user_id"`
	Amount    Amount    `json:"amount"`
	Balance   Amount    `json:"balance"`
	Timestamp time.Time `json:"timestamp"`
}

type Emitter interface {
	Publish(e Event)
}
