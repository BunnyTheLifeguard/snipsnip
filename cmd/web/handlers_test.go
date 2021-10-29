package main

import (
	"bytes"
	"net/http"
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
