package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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
	mu.Lock() // Lock the client to prevent concurrent updates as the whole graph has to be updated each time.

	defer mu.Unlock() // Unlock the client when the function returns
	url := fmt.Sprintf("%s/api/collection/graph", c.BaseURL)

	retries := 10
	updated := CollectionGraph{}
	var err_ret error
	for retries >= 0 {
		currentCG, err := c.GetCollectionGraph()
		if err != nil {
			return cg, err
		}
		// Update the revision which is incremented +1 by the server everytime the graph is updated
		cg.Revision = currentCG.Revision
		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(cg)
		req, err := http.NewRequest(http.MethodPut, url, b)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			err_ret = err
			break
		}
		if err := c.sendRequest(req, &updated); err != nil {
			if strings.Contains(err.Error(), "collection_revision_pkey") || strings.Contains(err.Error(), "status code: 409") {
				log.Printf("[ERROR] There were an error with the collection graph, retries left %d.", retries)
				time.Sleep(time.Duration(rand.Float32()*2) * time.Second)
				retries -= 1
			} else {
				err_ret = err
				break
			}
		} else {
			return updated, nil
		}
	}
	return cg, err_ret
}
