package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Users struct {
	Data []User `json:"data"`
}

func (c *Client) GetUsers() (Users, error) {
	if c.users != nil {
		return *c.users, nil
	}
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
	c.users = &users
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
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(u)
	req, err := http.NewRequest(http.MethodPost, url, b)
	req.Header.Set("Content-Type", "application/json")
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

func (c *Client) UpdateUser(u User, id int) (User, error) {
	url := fmt.Sprintf("%s/api/user/%d", c.BaseURL, id)
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(u)
	req, err := http.NewRequest(http.MethodPut, url, b)
	req.Header.Set("Content-Type", "application/json")
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
