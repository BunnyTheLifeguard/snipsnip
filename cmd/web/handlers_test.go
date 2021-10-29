package main

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	if code != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestShowSnip(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snip/617ab3734902e994b0d646ec", http.StatusOK, []byte("0 Test Snip")},
		{"Non-existent ID", "/snip/617ab3734902e994b0000000", http.StatusNotFound, nil},
		{"Int ID", "/snip/1", http.StatusNotFound, nil},
		{"Decimal ID", "/snip/1.23", http.StatusNotFound, nil},
		{"String ID", "/snip/abc", http.StatusNotFound, nil},
		{"Empty ID", "/snip/", http.StatusNotFound, nil},
		{"Trayling slash", "/snip/617ab3734902e994b0d646ec/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Errorf("want %d got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q", tt.wantBody)
			}
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	t.Log(csrfToken)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Sir Test", "sirtest@snipsnip.com", "validPa$$word", csrfToken, http.StatusSeeOther, nil},
		{"Empty name", "", "sirtest@snipsnip.com", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank.")},
		{"Empty email", "Sir Test", "", "validPa$$word", csrfToken, http.StatusOK, []byte("This field cannot be blank.")},
		{"Empty password", "Sir Test", "sirtest@snipsnip.com", "", csrfToken, http.StatusOK, []byte("This field cannot be blank.")},
		{"Invalid email (incomplete domain)", "Sir Test", "sirtest@snipsnip.", "validPa$$word", csrfToken, http.StatusOK, []byte("The email address is invalid.")},
		{"Invalid email (missing @)", "Sir Test", "sirtestsnipsnip.com", "validPa$$word", csrfToken, http.StatusOK, []byte("The email address is invalid.")},
		{"Invalid email (missing local part)", "Sir Test", "@snipsnip.com", "validPa$$word", csrfToken, http.StatusOK, []byte("The email address is invalid.")},
		{"Short password", "Sir Test", "sirtest@snipsnip.com", "pa$$word", csrfToken, http.StatusOK, []byte("This field must have at least 10 characters")},
		{"Duplicate email", "Sir Test", "dupe@snipsnip.com", "validPa$$word", csrfToken, http.StatusOK, []byte("Address already in use.")},
		{"Invalid CSRF Token)", "", "", "", "wrongToken", http.StatusBadRequest, []byte("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			if code != tt.wantCode {
				t.Errorf("want %d got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body %s to contain %q", body, tt.wantBody)
			}
		})
	}

}
