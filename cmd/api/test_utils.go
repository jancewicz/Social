package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jancewicz/social/internal/store"
	"github.com/jancewicz/social/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApp(t *testing.T) *application {
	t.Helper()

	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	return &application{
		logger:     logger,
		store:      mockStore,
		cacheStore: mockCacheStore,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	// request recorder
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}
