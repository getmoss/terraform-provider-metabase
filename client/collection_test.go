package client

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection(t *testing.T) {
	t.Run("Get collection", func(t *testing.T) {
		expected := Collections{
			{
				Id:    1,
				Name:  "TestCollection",
				Color: "#FFFFFF",
			},
			{
				Id:    2,
				Name:  "TestCollection",
				Color: "#000000",
			},
		}

		httpMethod := http.MethodGet
		svr := server("/api/collection", httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		col, err := c.GetCollections()

		assert.Nil(t, err)
		assert.Equal(t, expected, col)
	})

	t.Run("Get collection by id", func(t *testing.T) {
		collectionId := "root"
		expected := Collection{
			Id:    1,
			Name:  "TestCollection",
			Color: "#FFFFFF",
		}

		url := fmt.Sprintf("/api/collection/%s", collectionId)
		httpMethod := http.MethodGet
		svr := server(url, httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		col, err := c.GetCollection(collectionId)

		assert.Nil(t, err)
		assert.Equal(t, expected, col)
	})

	t.Run("Create collection", func(t *testing.T) {
		collectionToBeCreated := Collection{
			Name:  "TestCollection",
			Color: "#FFFFFF",
		}
		expected := Collection{
			Id:    1,
			Name:  "TestCollection",
			Color: "#FFFFFF",
		}
		httpMethod := http.MethodPost
		url := "/api/collection"
		svr := server(url, httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		col, err := c.CreateCollection(collectionToBeCreated)

		assert.Nil(t, err)
		assert.Equal(t, expected, col)
	})

	t.Run("Update collection", func(t *testing.T) {
		collectionToBeUpdated := Collection{
			Id:    1,
			Name:  "TestRed",
			Color: "#FF0000",
		}
		expected := Collection{
			Id:    1,
			Name:  "TestRed",
			Color: "#FF0000",
		}
		httpMethod := http.MethodPut
		svr := server("/api/collection/"+fmt.Sprintf("%d", collectionToBeUpdated.Id), httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		col, err := c.UpdateCollection(collectionToBeUpdated)

		assert.Nil(t, err)
		assert.Equal(t, expected, col)
	})

	t.Run("Archive collection", func(t *testing.T) {
		collectionToBeArchived := Collection{
			Id: 1,
		}
		expected := Collection{
			Id:       1,
			Name:     "TestRed",
			Color:    "#FF0000",
			Archived: true,
		}
		httpMethod := http.MethodPut
		svr := server("/api/collection/"+fmt.Sprintf("%d", collectionToBeArchived.Id), httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		col, err := c.UpdateCollection(collectionToBeArchived)

		assert.Nil(t, err)
		assert.Equal(t, expected, col)
	})

	t.Run("Get collections & cache them so that subsequent calls are not making network call", func(t *testing.T) {
		expected := Collections{
			{
				Id:    1,
				Name:  "TestCollection",
				Color: "#000000",
			},
			{
				Id:    2,
				Name:  "TestCollection",
				Color: "#000000",
			},
		}
		svr := server("/api/collection", http.MethodGet, expected)
		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		orig, errOrig := c.GetCollections()
		svr.Close() // close so the mock server is not running
		later, errLater := c.GetCollections()

		assert.Nil(t, errOrig)
		assert.Nil(t, errLater)
		assert.Equal(t, orig, later)
	})
}
