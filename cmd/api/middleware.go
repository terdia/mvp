package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/tomasen/realip"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/validator"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		app.logger.Printf("incoming request: %+v", map[string]string{
			"ip":     realip.FromRequest(r),
			"proto":  r.Proto,
			"method": r.Method,
			"uri":    r.URL.RequestURI(),
		})

		next.ServeHTTP(rw, r)
	})
}

func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(rw, r)

			return
		}

		parts := strings.Split(authorizationHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(rw, r)
			return
		}

		token := parts[1]

		v := validator.New()
		v.Check(token != "", "token", "must be provided")
		v.Check(len(token) == 26, "token", "must be 26 bytes long")
		if !v.Valid() {
			app.invalidAuthenticationTokenResponse(rw, r)
			return
		}

		user, err := app.userService.GetUserByToken(token, data.TokenScopeAuthentication)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(rw, r)
			default:
				app.serverErrorResponse(rw, r, err)
			}

			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(rw, r)
	}
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user == nil || user.IsAnonymous() {
			app.authenticationRequiredResponse(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
	}

	return app.authenticate(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {

	fn := func(rw http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.userService.GetPermissions(user.ID)
		if err != nil {
			app.serverErrorResponse(rw, r, err)
			return
		}

		if !permissions.Includes(code) {
			app.notPermittedRResponse(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
	}

	return app.requireAuthenticatedUser(fn)
}

func (app *application) enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw.Header().Set("Vary", "Origin")
		rw.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.config.Cors.TrustedOrigins {
				if origin == app.config.Cors.TrustedOrigins[i] {
					rw.Header().Set("Access-Control-Allow-Origin", origin)

					// handle prefight
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						rw.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						rw.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						rw.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(rw, r)
	})
}
