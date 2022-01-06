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
	HTTPClient *http.Client

	sessionId string
	userAgent string
}

type LoginDetails struct {
	Host      string
	Username  string
	Password  string
	SessionId string
	UserAgent string
}

type LoginSuccess struct {
	Client    *Client
	SessionId string
}

func NewClient(l LoginDetails) (LoginSuccess, error) {
	log.Printf("[INFO] Creating new client for host '%s'", l.Host)
	httpClient := &http.Client{
		Timeout: time.Minute,
	}
	var sessionId string

	// Re-use existing sessionId if possible
	if l.SessionId != "" {
		log.Printf("[DEBUG] Checking if existing sessionId is valid")

		// if err != nil {
		// 	return nil, diag.FromErr(err)
		// }
		// return c, nil
	}

	// Login with username/password
	creds := map[string]string{"username": l.Username, "password": l.Password}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(creds)

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/session", l.Host), b)
	req.Header.Set("Content-Type", "application/json")

	res, err := httpClient.Do(req)
	if err != nil {
		log.Panicf("[ERROR] Error in authentication: %s", err)
		return LoginSuccess{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			log.Panicf("[ERROR] Error in request[%+v]: Got response[status=%+v, username=%s]", req.URL, res.Status, l.Username)
		}
		log.Panicf("[ERROR] unknown error, status code: %d", res.StatusCode)
	}

	var resp LoginResponse
	json.NewDecoder(res.Body).Decode(&resp)
	sessionId = resp.Id

	return LoginSuccess{
		Client: &Client{
			BaseURL:    l.Host,
			sessionId:  sessionId,
			HTTPClient: httpClient,
			userAgent:  l.UserAgent,
		},
		SessionId: sessionId,
	}, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("X-Metabase-Session", c.sessionId)
	req.Header.Set("User-Agent", c.userAgent)
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
