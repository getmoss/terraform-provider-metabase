package client

import (
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("1st time login", func(t *testing.T) {
		sessionId, _ := uuid.GenerateUUID()
		loginResponse := LoginResponse{Id: sessionId}
		url := "/api/session"
		httpMethod := http.MethodPost

		svr := server(url, httpMethod, loginResponse)
		defer svr.Close()

		l := LoginDetails{
			Host:      svr.URL,
			Username:  "user",
			Password:  "pass",
			SessionId: "",
			UserAgent: "test",
		}

		client, err := NewClient(l)

		assert.Nil(t, err)
		assert.Equal(t, sessionId, client.SessionId)
		assert.Equal(t, svr.URL, client.Client.BaseURL)
	})

	t.Run("Re-use sessionId", func(t *testing.T) {
		sessionId, _ := uuid.GenerateUUID()
		url := "/api/user/current"
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, nil)
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
			_ = json.NewEncoder(w).Encode(LoginResponse{Id: sessionId})
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
