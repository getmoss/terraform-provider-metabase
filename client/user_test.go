package client

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser(t *testing.T) {
	t.Run("Get users", func(t *testing.T) {
		expected := Users{
			Data: []User{{
				Id:    1,
				Email: "test@example.com",
			}},
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				json.NewEncoder(w).Encode(expected)
			default:
				w.WriteHeader(http.StatusBadRequest)
			}

		})
		svr := httptest.NewServer(mux)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		us, err := c.GetUsers()

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Create user", func(t *testing.T) {

	})

	t.Run("Delete user", func(t *testing.T) {

	})
}
