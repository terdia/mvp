package main

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {

	app := createTestApplication(t, false)
	ts := newTestServer(t, app.routes())

	rs := ts.get(t, "/v1/healthcheck")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	result := struct {
		Status  string
		Message string
	}{}

	_ = json.Unmarshal(rs.Body, &result)

	if result.Status != "success" {
		t.Errorf("want %s; got %s", "success", result.Status)
	}

	if result.Message != "available" {
		t.Errorf("want %s; got %s", "available", result.Message)
	}

}
