package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jancewicz/social/internal/auth"
	"github.com/jancewicz/social/internal/ratelimiter"
	"github.com/jancewicz/social/internal/store"
	"github.com/jancewicz/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApp(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	mockAuth := &auth.TestAuthenticator{}

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.ratelimiter.RequestsPerTimeFrame,
		cfg.ratelimiter.TimeFrame,
	)

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStore:    mockCacheStore,
		authenticator: mockAuth,
		config:        cfg,
		rateLimiter:   rateLimiter,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	// request recorder
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseStatusCode(t *testing.T, want, got int) {
	if want != got {
		t.Errorf("expected status code to be %d, got %d", want, got)
	}
}
