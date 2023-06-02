package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func server(url string, httpMethod string, expected interface{}) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case httpMethod:
			_ = json.NewEncoder(w).Encode(expected)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	svr := httptest.NewServer(mux)
	return svr
}

func graphServer() *httptest.Server {
	mux := http.NewServeMux()

	// Maintain an in-memory CollectionGraph that starts with a default value
	cg := CollectionGraph{
		Revision: 1,
		Groups: map[string]map[string]string{
			"1": {
				"1":  "read",
				"10": "write",
				"9":  "none",
			},
			"4": {
				"1":  "none",
				"10": "none",
				"9":  "none"},
		},
	}

	mux.HandleFunc("/api/collection/graph", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(cg)
		case http.MethodPut:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var newCg CollectionGraph
			err = json.Unmarshal(body, &newCg)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			cg = newCg
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	svr := httptest.NewServer(mux)
	return svr
}
