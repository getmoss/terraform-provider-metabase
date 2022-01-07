package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	t.Run("1st time login", func(t *testing.T) {
		sessionId, _ := uuid.GenerateUUID()
		mux := http.NewServeMux()
		mux.HandleFunc("/api/session", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{ "id": "%s" }`, sessionId)
		})
		svr := httptest.NewServer(mux)
		defer svr.Close()

		l := LoginDetails{
			Host:      svr.URL,
			Username:  "user",
			Password:  "pass",
			SessionId: "",
			UserAgent: "test",
		}

		success, err := NewClient(l)

		assert.Nil(t, err)
		assert.NotNil(t, success.Client)
		assert.Equal(t, sessionId, success.SessionId)
	})

	t.Run("Re-use sessionId", func(t *testing.T) {
		sessionId, _ := uuid.GenerateUUID()
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user/current", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{ "id": 1 }`)
		})
		svr := httptest.NewServer(mux)
		defer svr.Close()

		l := LoginDetails{
			Host:      svr.URL,
			Username:  "irrelevant",
			Password:  "irrelevant",
			SessionId: sessionId,
			UserAgent: "test",
		}

		success, err := NewClient(l)

		assert.Nil(t, err)
		assert.NotNil(t, success.Client)
		assert.Equal(t, sessionId, success.SessionId)
	})

	t.Run("Login with username/password on expired sessionId", func(t *testing.T) {
		sessionId, _ := uuid.GenerateUUID()
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user/current", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		})
		mux.HandleFunc("/api/session", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, `{ "id": "%s" }`, sessionId)
		})
		svr := httptest.NewServer(mux)
		defer svr.Close()

		l := LoginDetails{
			Host:      svr.URL,
			Username:  "user",
			Password:  "pass",
			SessionId: sessionId,
			UserAgent: "test",
		}

		success, err := NewClient(l)

		assert.Nil(t, err)
		assert.NotNil(t, success.Client)
		assert.Equal(t, sessionId, success.SessionId)
	})
}
