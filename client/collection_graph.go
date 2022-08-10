package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type CollectionGraph struct {
	Revision int                          `json:"revision"`
	Groups   map[string]map[string]string `json:"groups"`
}

func (c *Client) GetCollectionGraph() (CollectionGraph, error) {
	url := fmt.Sprintf("%s/api/collection/graph", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	collectionGraph := CollectionGraph{}
	if err != nil {
		return collectionGraph, err
	}
	if err := c.sendRequest(req, &collectionGraph); err != nil {
		return collectionGraph, err
	}

	log.Printf("[DEBUG] Got collection graph '%+v'", collectionGraph)
	return collectionGraph, nil
}

func (c *Client) UpdateCollectionGraph(cg CollectionGraph) (CollectionGraph, error) {
	url := fmt.Sprintf("%s/api/collection/graph", c.BaseURL)

	retries := 5
	updated := CollectionGraph{}
	var err_rr error
	for retries >= 0 {
		current_cg, err := c.GetCollectionGraph()
		if err != nil {
			return cg, err
		}

		cg.Revision = current_cg.Revision
		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(cg)
		req, err := http.NewRequest(http.MethodPut, url, b)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			err_rr = err
			break
		}
		if err := c.sendRequest(req, &updated); err != nil {
			if strings.Contains(err.Error(), "collection_revision_pkey") || strings.Contains(err.Error(), "status code: 409") {
				time.Sleep(1500 * time.Millisecond)
				retries -= 1
			} else {
				err_rr = err
				break
			}
		} else {
			return updated, nil
		}
	}
	return cg, err_rr
}
