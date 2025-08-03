package wallet

import "errors"

var (
	ErrNegativeAmount    = errors.New("Amount cannot be negative")
	ErrInsufficientFunds = errors.New("Insufficient funds")
)
