package middleware

import (
	"log"
	"net/http"
	"time"
)

func TimerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		next.ServeHTTP(w, r)
		t1 := time.Now()
		log.Printf("%s cost %v", r.URL.Path, t1.Sub(t0))
	})
}
