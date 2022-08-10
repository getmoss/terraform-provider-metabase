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
