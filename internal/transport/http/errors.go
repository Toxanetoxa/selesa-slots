package http

import (
	"errors"
	"net/http"

	"github.com/toxanetoxa/selesa-slots/internal/wallet"
)

func HTTPError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, wallet.ErrInsufficientFunds):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, wallet.ErrNegativeAmount):
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
