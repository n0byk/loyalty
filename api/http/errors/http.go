package errors

import (
	"net/http"

	"github.com/jackc/pgerrcode"
)

func HTTPErrorGenerate(err error, w http.ResponseWriter) {
	if err.Error() == pgerrcode.UniqueViolation {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("UniqueViolation"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("InternalServerError"))
}
