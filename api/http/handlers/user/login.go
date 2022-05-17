package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/api/http/request"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/helpers"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userLogin request.UserLogin

	if err := json.Unmarshal(body, &userLogin); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, errCode, err := config.App.Storage.UserLogin(r.Context(), userLogin.Login, userLogin.Password)
	if err != nil {
		errors.HTTPErrorGenerate(errCode, w)
		return
	}

	if helpers.CheckPasswordHash(user.UserPassword, userLogin.Password, user.UserSalt) {
		tokenString := middleware.MakeToken(user.UserID)

		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.Write([]byte("Bearer " + tokenString))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
