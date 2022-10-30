package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/service/productservice"
	repo "github.com/terdia/mvp/mocks/repository"
)

type TestResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// createTestApplication is the application for test
func createTestApplication(t *testing.T, mockProductRepo bool) *application {

	cfg := new(config)
	logger := zerolog.New(&bytes.Buffer{})

	ctrl := gomock.NewController(t)
	productRepo := repo.NewMockProductRepository(ctrl)
	newProductService := productservice.NewProductService(
		productRepo,
	)

	if mockProductRepo {
		productRepo.EXPECT().Get(gomock.Any()).Return(&data.Product{
			ID:              1,
			Cost:            100,
			Name:            "Lemonade",
			Seller:          data.User{ID: 2},
			CreatedAt:       time.Now(),
			AmountAvailable: 20,
		}, nil)
	}

	return &application{
		wg:             new(sync.WaitGroup),
		config:         cfg,
		logger:         &logger,
		productService: newProductService,
	}
}

// Define a custom testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initializes and returns a new instance of testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	t.Helper()
	ts := httptest.NewServer(h)

	// Disable redirect-following for the client.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// get makes a request to given urlPath
func (ts *testServer) get(t *testing.T, urlPath string) TestResponse {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close() //nolint
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	return TestResponse{
		StatusCode: rs.StatusCode,
		Header:     rs.Header,
		Body:       body,
	}
}
