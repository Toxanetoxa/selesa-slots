package wallet

import (
	"context"
	"sync"
)

type MemoryWallet struct {
	mu   sync.RWMutex
	data map[UserID]*Wallet
}

func NewMemoryWallet() *MemoryWallet {
	return &MemoryWallet{
		data: make(map[UserID]*Wallet),
	}
}

func (m *MemoryWallet) Deposit(_ context.Context, userID UserID, amount Amount) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.data[userID]
	if !ok {
		w = NewWallet(userID)
		m.data[userID] = w
	}
	w.Balance += amount

	return nil
}

func (m *MemoryWallet) Withdraw(_ context.Context, userID UserID, amount Amount) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	w, ok := m.data[userID]
	if !ok || w.Balance < amount {
		return ErrInsufficientFunds
	}

	w.Balance -= amount
	return nil
}

func (m *MemoryWallet) GetBalance(_ context.Context, userID UserID) (Amount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if w, ok := m.data[userID]; ok {
		return w.Balance, nil
	}

	return 0, nil
}
