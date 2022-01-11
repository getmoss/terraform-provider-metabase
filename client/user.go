package client

import (
	"fmt"
	"net/http"
)

type DeleteSuccess struct {
	Success bool `json:"success"`
}

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

func (c *Client) GetUser(id int) (User, error) {
	url := fmt.Sprintf("%s/api/user/%d", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	user := User{}
	if err != nil {
		return user, err
	}
	if err := c.sendRequest(req, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (c *Client) CreateUser(u User) (User, error) {
	url := fmt.Sprintf("%s/api/user", c.BaseURL)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return u, err
	}
	created := User{}
	if err := c.sendRequest(req, &created); err != nil {
		return u, err
	}

	return created, nil
}

func (c *Client) UpdateUser(u User) (User, error) {
	url := fmt.Sprintf("%s/api/user", c.BaseURL)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return u, err
	}
	created := User{}
	if err := c.sendRequest(req, &created); err != nil {
		return u, err
	}

	return created, nil
}

func (c *Client) DeleteUser(id int) (DeleteSuccess, error) {
	url := fmt.Sprintf("%s/api/user/%d", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	resp := DeleteSuccess{}
	if err != nil {
		return resp, err
	}

	if err := c.sendRequest(req, &resp); err != nil {
		return resp, err
	}

	return resp, nil
}
