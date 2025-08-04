package leaderboard

import "github.com/go-playground/validator/v10"

var validate = validator.New()

type UpdateRequest struct {
	UserID   int64 `json:"user_id"  validate:"required,gt=0"`
	Position int   `json:"position" validate:"required,gt=0"`
	Score    int64 `json:"score"    validate:"required,gte=0"`
}

func (r *UpdateRequest) Bind() error { return validate.Struct(r) }
