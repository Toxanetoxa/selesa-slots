package wallet

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryRepo(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryWallet()
	const uid UserID = 444444

	if err := repo.Deposit(ctx, uid, 100); err != nil {
		t.Fatalf("deposit: %v", err)
	}

	balance, _ := repo.GetBalance(ctx, uid)
	if balance != 100 {
		t.Fatalf("balance: %v, want 100", balance)
	}

	if err := repo.Withdraw(ctx, uid, 30); err != nil {
		t.Fatalf("withdraw: %v", err)
	}

	balance, _ = repo.GetBalance(ctx, uid)
	if balance != 70 {
		t.Fatalf("balance: %v, want 70", balance)
	}

	if err := repo.Withdraw(ctx, uid, 100); !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("expected insufficient funds error, got %v", err)
	}
}
