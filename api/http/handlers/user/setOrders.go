package user

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/n0byk/loyalty/api/http/errors"
	"github.com/n0byk/loyalty/api/http/middleware"
	"github.com/n0byk/loyalty/api/validator"
	"github.com/n0byk/loyalty/config"
)

func SetOrders(w http.ResponseWriter, r *http.Request) {

	order, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(string(order))
	if validator.ValidateLuhn(string(order)) {
		userID := middleware.GetTokenClaims(r)

		orderOwner, err := config.App.Storage.CheckOrder(r.Context(), string(order))
		if err != nil {
			fmt.Println("SetOrders - db", err)
			errors.HTTPErrorGenerate("InternalError", w)
			return
		}

		if userID == orderOwner {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(orderOwner))
			return
		}

		responce, err := config.App.Storage.SetOrder(r.Context(), string(order), userID)
		if err != nil {
			fmt.Println("SetOrders", err)
			errors.HTTPErrorGenerate(responce, w)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(responce))
		return

	}
	w.WriteHeader(http.StatusUnprocessableEntity)

}
