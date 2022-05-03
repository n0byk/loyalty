package user

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/config"
)

func GetWithdraws(w http.ResponseWriter, r *http.Request) {

	orders, err := config.App.Storage.GetWithdraws(r.Context(), middleware.GetTokenClaims(r))
	if err != nil {
		errors.HTTPErrorGenerate("InternalServerError", w)
		return
	}

	if len(orders) > 0 {
		render.JSON(w, r, orders)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("StatusNoContent"))
}
