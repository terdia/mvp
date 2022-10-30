package dto

import (
	"time"

	"github.com/terdia/mvp/internal/data"
)

type (
	ProductRequest struct {
		Name     *string `json:"name"`
		Cost     int     `json:"cost"`
		Quantity int     `json:"amount_available"`
	}

	ProductResponse struct {
		Product APIProduct `json:"product"`
	}

	APIProduct struct {
		ID              int64     `json:"id"`
		Cost            int       `json:"cost"`
		Name            string    `json:"name"`
		CreatedAt       time.Time `json:"created_at"`
		AmountAvailable int       `json:"amount_available"`
	}

	ListProductRequest struct {
		Name    string
		Filters data.Filters
	}

	ListProductResponse struct {
		Metadata *data.Metadata `json:"metadata,omitempty"`
		Products []APIProduct   `json:"products"`
	}

	BuyProductResponse struct {
		AmountSpent int `json:"amount_spent"`
		Product     struct {
			Name     string `json:"name"`
			Cost     int    `json:"cost"`
			Quantity int    `json:"quantity_purchased"`
		} `json:"product_details"`
		Change []int `json:"change"`
	}
)
