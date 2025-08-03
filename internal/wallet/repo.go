package wallet

import "context"

type Repo interface {
	Deposit(ctx context.Context, userID UserID, amount Amount) error
	Withdraw(ctx context.Context, userID UserID, amount Amount) error
	GetBalance(ctx context.Context, userID UserID) (Amount, error)
}
