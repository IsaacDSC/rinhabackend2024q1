package middleware

import (
	"context"
	"net/http"
	"time"
)

func TimeoutLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//ctx := context.WithValue(r.Context(), "user", "123")
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			next.ServeHTTP(w, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(`{"error": "process timeout"}`))
		case <-processDone:
		}

	})
}
