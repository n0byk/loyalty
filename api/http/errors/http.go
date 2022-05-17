package errors

import (
	"net/http"
)

func HTTPErrorGenerate(err string, w http.ResponseWriter) {
	switch err {
	case "UniqueViolation":
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("UniqueViolation"))
		return

	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("InternalServerError"))

	}

}
