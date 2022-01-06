package client

import (
	"fmt"
	"net/http"
)

func (c *Client) GetUsers() (Users, error) {
	url := fmt.Sprintf("%s/api/user", c.BaseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	users := Users{}
	if err != nil {
		return users, err
	}
	if err := c.sendRequest(req, &users); err != nil {
		return users, err
	}

	return users, nil
}
