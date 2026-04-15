package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"codeberg.org/d33trik/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()
	defer res.Body.Close()

	want := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), want)

	want = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("Referrer-Policy"), want)

	want = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), want)

	want = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), want)

	want = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), want)

	want = "Go"
	assert.Equal(t, res.Header.Get("Server"), want)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
