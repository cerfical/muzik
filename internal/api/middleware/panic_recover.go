package middleware

import (
	"net/http"

	apierrs "github.com/cerfical/muzik/internal/api/errors"
	"github.com/cerfical/muzik/internal/log"

	"github.com/pkg/errors"
)

func PanicRecover(log *log.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					log.Error("recovered from panic", errors.Errorf("%v", e))
					apierrs.InternalError(w, r)
				}
			}()

			next(w, r)
		}
	}
}
