package main

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/terdia/mvp/pkg/dto"
)

func TestGetProduct(t *testing.T) {

	app := createTestApplication(t, true)
	ts := newTestServer(t, app.routes())

	rs := ts.get(t, "/v1/products/1")

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	result := struct {
		Status string
		Data   dto.ProductResponse
	}{}

	_ = json.Unmarshal(rs.Body, &result)

	if result.Status != "success" {
		t.Errorf("want %s; got %s", "success", result.Status)
	}

	if result.Data.Product.Name != "Lemonade" {
		t.Errorf("want %v; got %v", "Lemonade", result.Data.Product.Name)
	}

	if result.Data.Product.Cost != 100 {
		t.Errorf("want %v; got %v", 100, result.Data.Product.Cost)
	}

	if result.Data.Product.AmountAvailable != 20 {
		t.Errorf("want %v; got %v", 20, result.Data.Product.AmountAvailable)
	}
}
