package dto

import (
	"time"

	"github.com/terdia/mvp/internal/data"
)

type ProductRequest struct {
	Name     *string `json:"name"`
	Cost     int     `json:"cost"`
	Quantity int     `json:"amount_available"`
}

type ProductResponse struct {
	Product APIProduct `json:"product"`
}

type APIProduct struct {
	ID              int64     `json:"id"`
	Cost            int       `json:"cost"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	AmountAvailable int       `json:"amount_available"`
}

type ListProductRequest struct {
	Name    string
	Filters data.Filters
}

type ListProductResponse struct {
	Metadata *data.Metadata `json:"metadata,omitempty"`
	Products []APIProduct   `json:"products"`
}
