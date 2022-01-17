package client

import (
	"encoding/json"
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
