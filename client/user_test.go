package client

import (
	"fmt"
	"net/http"
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
		httpMethod := http.MethodGet
		svr := server("/api/user", httpMethod, expected)
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

		url := fmt.Sprintf("/api/user/%d", userId)
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, expected)
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
		httpMethod := http.MethodPost
		url := "/api/user"
		svr := server(url, httpMethod, expected)
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
		httpMethod := http.MethodPut
		svr := server("/api/user", httpMethod, expected)
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

		url := fmt.Sprintf("/api/user/%d", userId)
		httpMethod := http.MethodDelete
		svr := server(url, httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		us, err := c.DeleteUser(userId)

		assert.Nil(t, err)
		assert.Equal(t, expected, us)
	})

	t.Run("Get users & cache them so that subsequent calls are not making network call", func(t *testing.T) {
		email := "test@example.com"
		expected := Users{
			Data: []User{{
				Id:    1,
				Email: email,
			}},
		}
		httpMethod := http.MethodGet
		svr := server("/api/user", httpMethod, expected)

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}
		orig, errOrig := c.GetUsers()
		svr.Close() // close so the mock server is not running
		later, errLater := c.GetUsers()

		assert.Nil(t, errOrig)
		assert.Nil(t, errLater)
		assert.Equal(t, orig, later)
	})
}
