package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// RequestLogger is a simple middleware used to log each request made to the server
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := fmt.Sprintf(`%s - - [%s]" "%s %s %s" - -`,
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL,
			r.Proto,
		)

		fmt.Println(log)

		next.ServeHTTP(w, r)
	})
}

// ResponseLogger is a simple middleware used to log each response the server sends to client
func ResponseLogger(w http.ResponseWriter, r *http.Request) {
	log := fmt.Sprintf(`%s - - [%s]" "%s %s %s" %s %s`,
		r.RemoteAddr,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.URL,
		r.Proto,
		r.Response.StatusCode,
		r.Response.ContentLength,
	)

	fmt.Println(log)
}
