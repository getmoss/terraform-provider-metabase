package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	sessionId  string
	HTTPClient *http.Client
}

func NewClient(baseURL, username, password string) (*Client, error) {
	log.Printf("[INFO] Creating new client for host '%s'", baseURL)
	httpClient := &http.Client{
		Timeout: time.Minute,
	}

	// Login
	creds := map[string]string{"username": username, "password": password}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(creds)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/session", baseURL), b)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Panicf("[ERROR] Error in authentication: %s", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			log.Panicf("[ERROR] Error in request[%+v]: Got response[status=%+v, username=%s]", req.URL, res.Status, username)
		}
		log.Panicf("[ERROR] unknown error, status code: %d", res.StatusCode)
	}

	var loginResponse LoginResponse
	json.NewDecoder(res.Body).Decode(&loginResponse)

	return &Client{
		BaseURL:    baseURL,
		sessionId:  loginResponse.Id,
		HTTPClient: httpClient,
	}, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("X-Metabase-Session", c.sessionId)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Not successful
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		b, _ := io.ReadAll(res.Body)
		if err = json.NewDecoder(bytes.NewReader(b)).Decode(&errRes); err == nil {
			log.Printf("[ERROR] Error in request[%+v]: Got response[status=%+v, errors=%+v]", req.URL, res.Status, errRes.Errors)
			return errors.New(fmt.Sprint(errRes.Errors))
		}
		return fmt.Errorf("status code: %d, error:%s", res.StatusCode, string(b))
	}

	// Successful but no body
	if res.StatusCode == http.StatusNoContent {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}
