package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionGraph(t *testing.T) {
	t.Run("Get collection graph", func(t *testing.T) {
		expected := CollectionGraph{
			Revision: 1,
			Groups: map[string]map[string]string{
				"1": {
					"1":  "read",
					"10": "write",
					"9":  "none",
				},
				"4": {
					"1":  "none",
					"10": "none",
					"9":  "none"},
			},
		}
		httpMethod := http.MethodGet
		svr := server("/api/collection/graph", httpMethod, expected)
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		pg, _ := c.GetCollectionGraph()

		assert.Equal(t, expected, pg)
	})
}

func TestUpdateCollectionGraph(t *testing.T) {
	t.Run("Update collection graph", func(t *testing.T) {
		updatedExpected := CollectionGraph{
			Revision: 2,
			Groups: map[string]map[string]string{
				"1": {
					"1":  "none",
					"10": "read",
					"9":  "write",
				},
				"4": {
					"1":  "none",
					"10": "none",
					"9":  "none"},
			},
		}

		svr := graphServer()
		defer svr.Close()

		c := Client{
			BaseURL:    svr.URL,
			HTTPClient: &http.Client{},
		}

		// Update the graph
		_, err := c.UpdateCollectionGraph(updatedExpected)
		if err != nil {
			t.Fatalf("Failed to update collection graph: %v", err)
		}

		// Get the graph after update
		actualUpdatedCG, err := c.GetCollectionGraph()
		if err != nil {
			t.Fatalf("Failed to get collection graph: %v", err)
		}

		assert.Equal(t, updatedExpected, actualUpdatedCG)
	})
}
