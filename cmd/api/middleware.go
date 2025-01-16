package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is missing"))
				return
			}

			// parse header
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicError(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// decode header, read base64
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicError(w, r, err)
				return
			}

			// check user credentials
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.password

			credentials := strings.SplitN(string(decoded), ":", 2)
			if len(credentials) != 2 || credentials[0] != username || credentials[1] != pass {
				app.unauthorizedBasicError(w, r, fmt.Errorf("invalid user's credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
