package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/dataservice/models/request"
	"github.com/n0byk/loyalty/helpers"
	"github.com/n0byk/loyalty/helpers/jwtauth"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userLogin request.UserLogin

	if err := json.Unmarshal(body, &userLogin); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := config.App.Storage.UserLogin(r.Context(), userLogin.Login, userLogin.Password)
	if err != nil {
		errors.HTTPErrorGenerate(err, w)
		return
	}

	if helpers.CheckPasswordHash(user.UserPassword, userLogin.Password, user.UserSalt) {
		tokenString := jwtauth.MakeToken(user.UserID)

		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.Write([]byte("Bearer " + tokenString))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
}
