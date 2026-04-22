package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestSnippetCreateE2E(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		ts.resetClientCookieJar(t)

		res := ts.get(t, "/snippet/create")
		assert.Equal(t, res.status, http.StatusSeeOther)
		assert.Equal(t, res.headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		ts.resetClientCookieJar(t)

		res := ts.get(t, "/user/login")

		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", extractCSRFToken(t, res.body))

		ts.postForm(t, "/user/authenticate", form)

		res = ts.get(t, "/snippet/create")
		assert.Equal(t, res.status, http.StatusOK)
		assert.True(t, strings.Contains(res.body, `<form action="/snippet/save" method="POST">`))
	})
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

func TestUserSignupE2E(t *testing.T) {
	const (
		validName     = "Bob"
		validPassword = "validPa$$word"
		validEmail    = "bob@test.com"
		formTag       = `<form action="/user/save" method="POST" novalidate>`
	)

	tests := map[string]struct {
		userName           string
		userEmail          string
		userPassword       string
		userValidCSRFToken bool
		wantStatus         int
		wantFormTag        string
	}{
		"Valid submission": {
			userName:           validName,
			userEmail:          validEmail,
			userPassword:       validPassword,
			userValidCSRFToken: true,
			wantStatus:         http.StatusSeeOther,
		},
		"Invalid CSRF Token": {
			userName:           validName,
			userEmail:          validEmail,
			userPassword:       validPassword,
			userValidCSRFToken: false,
			wantStatus:         http.StatusBadRequest,
		},
		"Empty name": {
			userName:           "",
			userEmail:          validEmail,
			userPassword:       validPassword,
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
		"Empty email": {
			userName:           validName,
			userEmail:          "",
			userPassword:       validPassword,
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
		"Empty password": {
			userName:           validName,
			userEmail:          validEmail,
			userPassword:       "",
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
		"Invalid email": {
			userName:           validName,
			userEmail:          "bob@test.",
			userPassword:       validPassword,
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
		"Short password": {
			userName:           validName,
			userEmail:          validEmail,
			userPassword:       "pa$$",
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
		"Duplicate email": {
			userName:           validName,
			userEmail:          "dupe@example.com",
			userPassword:       validPassword,
			userValidCSRFToken: true,
			wantStatus:         http.StatusUnprocessableEntity,
			wantFormTag:        formTag,
		},
	}

	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ts.resetClientCookieJar(t)

			res := ts.get(t, "/user/signup")

			form := url.Values{}
			form.Add("name", tc.userName)
			form.Add("email", tc.userEmail)
			form.Add("password", tc.userPassword)
			if tc.userValidCSRFToken {
				form.Add("csrf_token", extractCSRFToken(t, res.body))
			}

			res = ts.postForm(t, "/user/save", form)

			assert.Equal(t, res.status, tc.wantStatus)
			assert.True(t, strings.Contains(res.body, tc.wantFormTag))
		})
	}
}
