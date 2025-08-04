package game

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type OutcomeRequest struct {
	GameID  string  `json:"game_id"  validate:"required"`
	UserID  int64   `json:"user_id"  validate:"required,gt=0"`
	Outcome Outcome `json:"outcome"  validate:"required,oneof=win lose"`
	Amount  int64   `json:"amount"   validate:"required,gte=0"`
}

func (r *OutcomeRequest) Bind() error { return validate.Struct(r) }
