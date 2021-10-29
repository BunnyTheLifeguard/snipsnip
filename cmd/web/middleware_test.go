package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// Initialize a new ResponseRecorder & dummy http.Request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock HTTP handler that gets passed to secureHeaders middleware, writes a 200 status code and "OK" response body
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass mock HTTP handler to secureHeaders middleware, all ServeHTTP of returned http.Handler & pass in ResponseRecorder + dummy http.Request to exec
	secureHeaders(next).ServeHTTP(rr, r)

	// Get results of test
	rs := rr.Result()

	// Check if X-Frame-Options header was set correctly
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q got %q", "deny", frameOptions)
	}

	// Check if X-XSS-Protection header was set correctly
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q got %q", "1; mode=block", xssProtection)
	}

	// Check if middleware calls next handler in line
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
