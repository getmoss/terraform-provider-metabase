package client

import (
	"fmt"
	"log"
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

	log.Printf("[DEBUG] Got users '%+v'", users)
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

	log.Printf("[INFO] Got user '%+v'", user)
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

	log.Printf("[INFO] Created user '%+v'", created)
	return created, nil
}

func (c *Client) UpdateUser(u User) (User, error) {
	url := fmt.Sprintf("%s/api/user", c.BaseURL)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return u, err
	}
	updated := User{}
	if err := c.sendRequest(req, &updated); err != nil {
		return u, err
	}

	log.Printf("[INFO] Updated user '%+v'", updated)
	return updated, nil
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

	log.Printf("[INFO] Deleted user by id='%d'", id)
	return resp, nil
}
