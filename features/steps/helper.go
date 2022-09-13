package steps

import (
	"net/http"

	"github.com/maxcnunes/httpfake"
)

func healthCheckStatusHandle(status int) httpfake.Responder {
	return func(w http.ResponseWriter, r *http.Request, rh *httpfake.Request) {
		w.WriteHeader(status)
	}
}
