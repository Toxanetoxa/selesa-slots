package wallet

type Amount int64
type UserID int64

type Wallet struct {
	UserID  UserID
	Balance Amount
}

func NewWallet(userID UserID) *Wallet {
	return &Wallet{UserID: userID}
}
