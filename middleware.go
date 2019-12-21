package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/makeict/MESSforMakers/models"
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

type contextKey int

const contextkeyUser contextKey = 1

func (a *application) authenticationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check if cookie exists
		if a.Session.Exists(r, "user") && a.Session.Exists(r, "authKey") {

			idStr := a.Session.GetString(r, "user")
			auth := a.Session.GetString(r, "authKey")

			//if it exists, check user ID and token for validity.
			if idStr != "" && auth != "" {

				id, err := strconv.Atoi(idStr)
				if err != nil {
					a.Session.Remove(r, "user")
					a.Session.Remove(r, "authKey")
					h.ServeHTTP(w, r)
					return
				}

				user, err := a.UserC.Users.SessionLookup(id, auth)

				//if username and token are valid, store the username, and id in the context for use by later handlers
				if err == nil {
					user.Authorized = true
					// context key should be a custom type and a const NOT a string
					ctx := context.WithValue(r.Context(), contextkeyUser, user)

					// TODO if the user is authenticated, the LastSeenTime of the session should be updated

					h.ServeHTTP(w, r.WithContext(ctx))
					return
				} else if err == models.ErrNoRecord {
					// user is either not logged in, doesn't exist, or session has expired.
					// Remove the user id from the session and proceed as if there is no user logged in.
					a.Session.Remove(r, "user")
					a.Session.Remove(r, "authKey")
					h.ServeHTTP(w, r)
					return
				} else {
					a.Logger.Debugf("%v", err)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			}
		}

		h.ServeHTTP(w, r)

	})
}

func (a *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(contextkeyUser).(*models.User)
		if !ok {
			http.Redirect(w, r, fmt.Sprintf("http://%s:%d/login", a.Config.App.Host, a.Config.App.Port), http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
