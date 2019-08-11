package middlewares

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// RequestLogger is a simple middleware used to log each request made to the server
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, false)

		if err != nil {
			log.Printf("WARN: error while trying to dump the request: %v", err)
			next.ServeHTTP(w, r)

			return
		}

		log.Print(string(dump))

		next.ServeHTTP(w, r)
	})
}
