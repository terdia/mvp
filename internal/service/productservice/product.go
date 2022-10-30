package productservice

import (
	"errors"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
	"github.com/terdia/mvp/pkg/dto"
	"github.com/terdia/mvp/pkg/validator"
)

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (p *productService) Create(product *data.Product) (map[string]string, error) {

	v := validator.New()

	if product.Validate(v); !v.Valid() {
		return v.Errors, nil
	}

	if err := p.repo.Insert(product); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateProductName):
			v.AddError("name", err.Error())
			return v.Errors, nil
		default:
			return nil, err
		}
	}

	return nil, nil
}

func (p *productService) List(r dto.ListProductRequest) ([]*data.Product, data.Metadata, error) {
	return p.repo.GetAll(r)
}

func (p *productService) GetOne(id int64) (*data.Product, error) {
	return p.repo.Get(id)
}

func (p *productService) Update(request data.Product) (*data.Product, map[string]string, error) {
	v := validator.New()
	request.Validate(v)
	if !v.Valid() {
		return nil, v.Errors, nil
	}

	product, err := p.getForUser(request)
	if err != nil {
		return nil, nil, err
	}

	if err = p.repo.Update(product); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateProductName):
			v.AddError("name", err.Error())
			return nil, v.Errors, nil
		default:
			return nil, nil, err
		}
	}

	return product, nil, nil
}

func (p *productService) Remove(request data.Product) error {
	product, err := p.getForUser(request)
	if err != nil {
		return err
	}

	return p.repo.Delete(product.ID)
}

func (p *productService) getForUser(request data.Product) (*data.Product, error) {
	product, err := p.GetOne(request.ID)
	if err != nil {
		return nil, err
	}

	if product.Seller.ID != request.Seller.ID {
		return nil, data.ErrNoPermission
	}

	product.Name = request.Name
	product.Cost = request.Cost
	product.AmountAvailable = request.AmountAvailable

	return product, nil
}
