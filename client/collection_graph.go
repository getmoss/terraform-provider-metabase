package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CollectionGraph struct {
	Revision int                          `json:"revision"`
	Groups   map[string]map[string]string `json:"groups"`
}

func (c *Client) GetCollectionGraph() (CollectionGraph, error) {
	if c.collectionGraph != nil {
		return *c.collectionGraph, nil
	}
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
	c.collectionGraph = &collectionGraph
	return collectionGraph, nil
}

func (c *Client) UpdateCollectionGraph(cg CollectionGraph) (CollectionGraph, error) {
	url := fmt.Sprintf("%s/api/collection/graph", c.BaseURL)

	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(cg)
	req, err := http.NewRequest(http.MethodPut, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return cg, err
	}
	updated := CollectionGraph{}
	if err := c.sendRequest(req, &updated); err != nil {
		return cg, err
	}

	log.Printf("[INFO] Updated collection graph '%+v'", updated)
	return updated, nil
}
