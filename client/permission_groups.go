package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (c *Client) GetPermissionGroups() (PermissionGroups, error) {
	url := fmt.Sprintf("%s/api/permissions/group", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	pg := PermissionGroups{}
	if err := c.sendRequest(req, &pg); err != nil {
		return nil, err
	}

	return pg, nil
}

func (c *Client) GetPermissionGroup(id int) (PermissionGroup, error) {
	url := fmt.Sprintf("%s/api/permissions/group/%d", c.BaseURL, id)
	pg := PermissionGroup{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return pg, err
	}

	if err := c.sendRequest(req, &pg); err != nil {
		return pg, err
	}

	log.Printf("[INFO] Got permissionGroup '%+v'", pg)
	return pg, nil
}

func (c *Client) CreatePermissionGroup(name string) (PermissionGroup, error) {
	url := fmt.Sprintf("%s/api/permissions/group", c.BaseURL)
	pg := PermissionGroup{Name: name}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(pg)
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return pg, err
	}
	if err := c.sendRequest(req, &pg); err != nil {
		return pg, err
	}

	log.Printf("[INFO] Created new permissionGroup '%+v'", pg)
	return pg, nil
}

func (c *Client) DeletePermissionGroup(id int) error {
	url := fmt.Sprintf("%s/api/permissions/group/%d", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	log.Printf("[INFO] Deleted permissionGroup with id[%d]", id)
	return nil
}
