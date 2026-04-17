package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codeberg.org/d33trik/snippetbox/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	healthCheck(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func TestHealthCheckE2E(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	res := ts.get(t, "/healthcheck")
	assert.Equal(t, res.status, http.StatusOK)
	assert.Equal(t, res.body, "OK")
}

func TestSnippetViewE2E(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := map[string]struct {
		urlPath    string
		wantStatus int
		wantBody   string
	}{
		"Valid ID": {
			urlPath:    "/snippet/view/1",
			wantStatus: http.StatusOK,
			wantBody:   "An old silent pond...",
		},
		"Non-existent ID": {
			urlPath:    "/snippet/view/2",
			wantStatus: http.StatusNotFound,
		},
		"Negative ID": {
			urlPath:    "/snippet/view/-1",
			wantStatus: http.StatusNotFound,
		},
		"Decimal ID": {
			urlPath:    "/snippet/view/1.23",
			wantStatus: http.StatusNotFound,
		},
		"String ID": {
			urlPath:    "/snippet/view/foo",
			wantStatus: http.StatusNotFound,
		},
		"Empty ID": {
			urlPath:    "/snippet/view/",
			wantStatus: http.StatusNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, tc.urlPath)
			assert.Equal(t, res.status, tc.wantStatus)
			assert.True(t, strings.Contains(res.body, tc.wantBody))
		})
	}
}
