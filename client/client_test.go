package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("1st time login", func(t *testing.T) {
		expected := "f4a53a0e-b1c9-4dda-a5b9-d3dce5caa959"
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{ "id": "%s" }`, expected)
		}))
		defer svr.Close()

		l := LoginDetails{
			Host:      svr.URL,
			Username:  "user",
			Password:  "pass",
			SessionId: "",
			UserAgent: "test",
		}

		success, err := NewClient(l)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if success.Client == nil {
			t.Fatalf("client is nil")
		}

		if success.SessionId != expected {
			t.Fatalf("unexpected session id: '%s'", success.SessionId)
		}
	})

	t.Run("Re-use sessionId", func(t *testing.T) {

	})

	t.Run("Login with username/password on expired sessionId", func(t *testing.T) {

	})
}
