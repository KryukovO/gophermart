package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type MockAccrual struct {
	*httptest.Server
}

func NewMockAccrual() *MockAccrual {
	server := httptest.NewServer(http.HandlerFunc(ordersHandler))

	return &MockAccrual{
		Server: server,
	}
}

func (a MockAccrual) Close() {
	a.Server.Close()
}

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 3 || pathParts[0] != "api" || pathParts[1] != "orders" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Unexected path: %s", r.URL.Path)))
	}

	order := fmt.Sprintf(`{"order":"%s","status":"PROCESSED","accrual":500}`, pathParts[2])

	w.Write([]byte(order))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
