package productservice

import (
	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/pkg/dto"
)

type ProductService interface {
	Create(*data.Product) (map[string]string, error)
	GetOne(int64) (*data.Product, error)
	Update(data.Product) (*data.Product, map[string]string, error)
	Remove(data.Product) error
	List(dto.ListProductRequest) ([]*data.Product, data.Metadata, error)
}
