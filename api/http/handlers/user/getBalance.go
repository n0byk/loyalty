package user

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/config"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := config.App.Storage.GetBalance(r.Context(), middleware.GetTokenClaims(r))
	if err != nil {
		fmt.Println(err)
		errors.HTTPErrorGenerate("InternalServerError", w)
		return
	}
	fmt.Println("GetBalance", balance)

	if balance.Current < 0 {
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	render.JSON(w, r, balance)

}
