package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{ResponseWriter: w}

		next.ServeHTTP(wrapper, r)
		log.Printf("METHOD: %s URI:%s STATUS CODE:%s TIME:%s", r.Method, r.RequestURI, wrapper.StatusCode, time.Since(start))
	})
}
