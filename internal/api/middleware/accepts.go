package middleware

import (
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/cerfical/muzik/internal/api/errors"
)

func Accepts(mediaType string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if checkAcceptHeader(mediaType, r.Header.Get("Accept")) {
				next(w, r)
				return
			}

			details := fmt.Sprintf("accept header doesn't specify '%s'", mediaType)
			e := errors.Error{
				Title:  "media type is not acceptable",
				Detail: details,
				Status: http.StatusNotAcceptable,
			}
			e.Write(w)
		}
	}
}

func checkAcceptHeader(supportedType, acceptHeader string) bool {
	if acceptHeader == "" {
		// Ignore empty Accept headers
		return true
	}

	supMain, supSub := splitMediaType(supportedType)
	accTypes := strings.Split(acceptHeader, ",")

	for _, accType := range accTypes {
		accType, params, err := mime.ParseMediaType(accType)
		if err != nil || len(params) > 0 {
			continue
		}

		accMain, accSub := splitMediaType(accType)
		if (accMain != supMain && accMain != "*") || (accSub != supSub && accSub != "*") {
			continue
		}

		return true
	}

	return false
}

func splitMediaType(mediaType string) (string, string) {
	mainType, subType, _ := strings.Cut(mediaType, "/")
	return mainType, subType
}
