package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/api/http/request"
	"github.com/n0byk/loyalty/config"
)

func PostWithdraw(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var UserWithdraw request.UserWithdraw

	if err := json.Unmarshal(body, &UserWithdraw); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = config.App.Storage.UpsertOrder(r.Context(), UserWithdraw.Order, middleware.GetTokenClaims(r))

	errCode, err := config.App.Storage.PostWithdraw(r.Context(), UserWithdraw.Order, float32(UserWithdraw.Sum))
	if err != nil {
		errors.HTTPErrorGenerate(errCode, w)
		return
	}

	w.WriteHeader(http.StatusOK)

}
