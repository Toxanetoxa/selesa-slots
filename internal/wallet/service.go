package wallet

import (
	"context"
	"time"
)

type Service struct {
	repo    Repo
	emitter Emitter
}

type noEmitter struct{}

func (noEmitter) Publish(Event) {}

func NewService(r Repo, e Emitter) *Service {
	if e == nil {
		e = noEmitter{}
	}

	return &Service{repo: r, emitter: e}
}

func (s *Service) Deposit(ctx context.Context, uID int64, a int64) error {
	userID, amount := UserID(uID), Amount(a)

	if err := s.repo.Deposit(ctx, userID, amount); err != nil {
		return err
	}
	balance, _ := s.repo.GetBalance(ctx, userID)

	s.emitter.Publish(Event{
		EventDeposit,
		userID,
		amount,
		balance,
		time.Now().UTC(),
	})
	return nil
}

func (s *Service) Withdraw(ctx context.Context, uID int64, a int64) error {
	userID, amount := UserID(uID), Amount(a)

	if err := s.repo.Withdraw(ctx, userID, amount); err != nil {
		return err
	}

	balance, _ := s.repo.GetBalance(ctx, userID)

	s.emitter.Publish(Event{
		EventWithdrawal,
		userID,
		amount,
		balance,
		time.Now().UTC(),
	})

	return nil
}

func (s *Service) GetBalance(ctx context.Context, uID int64) (int64, error) {
	userID := UserID(uID)

	amt, err := s.repo.GetBalance(ctx, userID)
	if err != nil {
		return 0, err
	}

	return int64(amt), nil
}
