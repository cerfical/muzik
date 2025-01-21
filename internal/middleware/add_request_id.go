package middleware

import (
	"context"
	"net/http"
	"sync/atomic"
)

// AddRequestID assigns unique identifiers to incoming requests via [context.Context].
func AddRequestID(next http.Handler) http.Handler {
	var requestID atomic.Uint64

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := requestID.Add(1)
		ctx := context.WithValue(r.Context(), RequestID, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestID is the context key to retrieve the generated request ID from [context.Context].
var RequestID requestIDKey

type requestIDKey struct{}

func (requestIDKey) String() string {
	return "request_id"
}
