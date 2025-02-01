package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	cfg := config{
		redisCfg: redisConfig{
			enable: true,
		},
	}
	app := newTestApp(t, cfg)
	mux := app.mount()

	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated request", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		// request recorder
		rr := executeRequest(req, mux)

		checkResponseStatusCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated user request API", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseStatusCode(t, http.StatusOK, rr.Code)
	})
}
