package http

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type amountReq struct {
	UserId int64 `json:"user_id" validate:"required,gt=0"`
	Amount int64 `json:"amount" validate:"required,gt=0"`
}

type balanceResp struct {
	UserId  int64 `json:"user_id"`
	Balance int64 `json:"balance"`
}

func (r *amountReq) Bind() error {
	return validate.Struct(r)
}
