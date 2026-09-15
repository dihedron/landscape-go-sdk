package landscape

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWithTokenAuthInitializesToken(t *testing.T) {
	api := New("https://example.com", WithTokenAuth("jwt-token"))
	if got := api.Token(); got != "jwt-token" {
		t.Fatalf("expected token to be initialized from option, got %q", got)
	}
}

func TestNewWithBasicAuthLogsInAndStoresToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/login" {
			http.Error(w, fmt.Sprintf("unexpected path: %s", r.URL.Path), http.StatusNotFound)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("unexpected method: %s", r.Method), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"token":"server-jwt"}`))
	}))
	defer server.Close()

	api := New(server.URL, WithBasicAuth("user@example.com", "secret", "default"))
	if got := api.Token(); got != "server-jwt" {
		t.Fatalf("expected token to be initialized after login, got %q", got)
	}
}
