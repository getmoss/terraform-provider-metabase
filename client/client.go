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

	sessionId        string
	userAgent        string
	users            *Users
	permissionGroups *PermissionGroups
	collections      *Collections
	collectionGraph  *CollectionGraph
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

type DeleteSuccess struct {
	Success bool `json:"success"`
}

type LoginResponse struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
}

func NewClient(l LoginDetails) (LoginSuccess, error) {
	log.Printf("[INFO] Creating new client for host '%s'", l.Host)
	httpClient := &http.Client{
		Timeout: time.Minute,
	}
	var sessionId string

	// Re-use existing sessionId if possible
	// Session Id is valid
	sessionId = loginWithSessionId(l, httpClient, sessionId)

	if sessionId == "" { // Login with username/password
		log.Printf("[DEBUG] Logging in with username/password")
		creds := map[string]string{"username": l.Username, "password": l.Password}
		b := new(bytes.Buffer)
		_ = json.NewEncoder(b).Encode(creds)

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
		_ = json.NewDecoder(res.Body).Decode(&resp)
		sessionId = resp.Id
	}

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

func loginWithSessionId(l LoginDetails, httpClient *http.Client, sessionId string) string {
	if l.SessionId != "" {
		log.Printf("[DEBUG] Checking if existing sessionId is valid")
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/user/current", l.Host), nil)
		res, err := httpClient.Do(req)
		if err != nil {
			log.Printf("[WARN] Error fetching current user: %s", err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {

			sessionId = l.SessionId
		}
	}
	return sessionId
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
			errorMsg := fmt.Sprintf("errors='%+v', message='%s'", errRes.Errors, errRes.Message)
			log.Printf("[ERROR] Error in request[%+v]: Got response[status='%+v', errors='%s']", req.URL, res.Status, errorMsg)
			return errors.New(errorMsg)
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
