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

	// Sometimes we try to upgrade the graph but the revision is old.
	// For this cases, wait and retry, fetching the latest revision
	retries := 5

	for retries > 0 {
		if err := c.sendRequest(req, &updated); err != nil {
			if strings.Contains(err.Error(), "collection_revision_pkey") {
				log.Println("[ERROR] DANDNDNDNDNDNDNDNDNDNDNND")
				retries -= 1
				time.Sleep(1500 * time.Millisecond)
				cg, _ := c.GetCollectionGraph()
				updated.Revision = cg.Revision
			}
			break
		} else {
			return updated, nil
		}
	}
	return cg, err
}
