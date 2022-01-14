package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(expected)
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

	t.Run("Get user by id", func(t *testing.T) {
		userId := 1
		expected := User{
			Id:        userId,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/api/user/%d", userId), func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				_ = json.NewEncoder(w).Encode(expected)
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

		us, err := c.GetUser(userId)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Create user", func(t *testing.T) {
		userToBeCreated := User{
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}
		expected := User{
			Id:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				_ = json.NewEncoder(w).Encode(expected)
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

		us, err := c.CreateUser(userToBeCreated)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Update user", func(t *testing.T) {
		userToBeUpdated := User{
			Id:        1,
			Email:     "test@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}
		expected := User{
			Id:        1,
			Email:     "test.updated@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPut:
				_ = json.NewEncoder(w).Encode(expected)
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

		us, err := c.UpdateUser(userToBeUpdated)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Delete user", func(t *testing.T) {
		userId := 1
		expected := DeleteSuccess{Success: true}

		mux := http.NewServeMux()
		mux.HandleFunc(fmt.Sprintf("/api/user/%d", userId), func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodDelete:
				_ = json.NewEncoder(w).Encode(expected)
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

		us, err := c.DeleteUser(userId)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})
}
