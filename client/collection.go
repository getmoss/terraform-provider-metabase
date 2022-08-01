package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Collection struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Archived bool   `json:"archived"`
}

type Collections []Collection

func (c *Client) GetCollections() (Collections, error) {
	if c.collections != nil {
		return *c.collections, nil
	}
	url := fmt.Sprintf("%s/api/collection", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	collections := Collections{}
	if err != nil {
		return collections, err
	}
	if err := c.sendRequest(req, &collections); err != nil {
		return collections, err
	}

	log.Printf("[DEBUG] Got collections '%+v'", collections)
	c.collections = &collections
	return collections, nil
}

func (c *Client) GetCollection(id string) (Collection, error) {
	url := fmt.Sprintf("%s/api/collection/%s", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	collection := Collection{}
	if err != nil {
		return collection, err
	}
	if err := c.sendRequest(req, &collection); err != nil {
		return collection, err
	}

	log.Printf("[INFO] Got collection '%+v'", collection)
	return collection, nil
}

func (c *Client) CreateCollection(col Collection) (Collection, error) {
	url := fmt.Sprintf("%s/api/collection", c.BaseURL)
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(col)
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return col, err
	}
	created := Collection{}
	if err := c.sendRequest(req, &created); err != nil {
		return col, err
	}

	log.Printf("[INFO] Created collection '%+v'", created)
	return created, nil
}

func (c *Client) UpdateCollection(col Collection) (Collection, error) {
	url := fmt.Sprintf("%s/api/collection/%v", c.BaseURL, col.Id)
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(col)
	req, err := http.NewRequest(http.MethodPut, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return col, err
	}
	updated := Collection{}
	if err := c.sendRequest(req, &updated); err != nil {
		return col, err
	}

	log.Printf("[INFO] Updated collection '%+v'", updated)
	return updated, nil
}
