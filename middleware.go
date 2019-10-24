package main

import (
	"net/http"
	"net/http/httputil"
	"time"
)

func (a *application) loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		if a.Logger.DumpRequest {
			if reqDump, err := httputil.DumpRequest(r, true); err == nil {
				a.Logger.Printf("Recieved Request:\n%s\n", reqDump)
			}
		}

		h.ServeHTTP(w, r)

		t2 := time.Now()

		a.Logger.Printf("[%s] %s %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	})
}

func (a *application) recoverPanic(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				a.Logger.Printf("Panic recovered: %v", err)
				w.Header().Set("Connection", "close")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		h.ServeHTTP(w, r)
	})
}

func (a *application) securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")
		h.ServeHTTP(w, r)
	})
}
