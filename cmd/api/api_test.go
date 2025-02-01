package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jancewicz/social/internal/ratelimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		ratelimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Second * 5,
			Enabled:              true,
		},
		addr: ":8080",
	}

	app := newTestApp(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIPAddr := "192.168.1.1"
	marginOfErr := 5

	for i := 0; i < cfg.ratelimiter.RequestsPerTimeFrame+marginOfErr; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("couldn't create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIPAddr)

		res, err := client.Do(req)
		if err != nil {
			t.Fatalf("couldn't send request: %v", err)
		}
		defer res.Body.Close()

		if i < cfg.ratelimiter.RequestsPerTimeFrame {
			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK; got %v", res.Status)
			}
		} else {
			if res.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status Too Many Requests, got: %v", res.Status)
			}
		}
	}
}
