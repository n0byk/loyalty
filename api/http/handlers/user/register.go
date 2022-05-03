package user

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/api/http/request"
	"github.com/n0byk/loyalty/config"
	"github.com/n0byk/loyalty/helpers"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userRegistration request.UserRegistration

	if err := json.Unmarshal(body, &userRegistration); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var salt = uuid.New().String()
	password, err := helpers.HashPassword(userRegistration.Password, salt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, errCode, err := config.App.Storage.UserRegister(r.Context(), userRegistration.Login, password, salt)
	if err != nil {
		errors.HTTPErrorGenerate(errCode, w)
		return
	}
	tokenString := middleware.MakeToken(userID.String())

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.Write([]byte("Bearer " + tokenString))
	w.WriteHeader(http.StatusOK)
}
