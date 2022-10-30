package repositoryproduct

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/terdia/mvp/internal/data"
	"github.com/terdia/mvp/internal/repository"
	"github.com/terdia/mvp/pkg/dto"
)

type productRepository struct {
	*sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db}
}

func (repo *productRepository) Insert(product *data.Product) error {
	query := `INSERT INTO products (name, cost, quantity, seller_id)
			 VALUES($1, $2, $3, $4)
			 RETURNING id, name, cost, quantity, created_at`

	queryParams := []interface{}{product.Name, product.Cost, product.AmountAvailable, product.Seller.ID}

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	if err := repo.DB.QueryRowContext(ctx, query, queryParams...).Scan(
		&product.ID, &product.Name, &product.Cost, &product.AmountAvailable, &product.CreatedAt,
	); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "products_name_seller_id_key"`:
			return data.ErrDuplicateProductName
		default:
			return err
		}
	}

	return nil
}

func (repo *productRepository) Get(id int64) (*data.Product, error) {

	if id < 1 {
		return nil, data.ErrRecordNotFound
	}

	query := `SELECT id, name, cost, quantity, seller_id, created_at
			  FROM products
			  WHERE id = $1`

	var product data.Product

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Cost,
		&product.AmountAvailable,
		&product.Seller.ID,
		&product.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

func (repo *productRepository) Update(product *data.Product) error {
	query := `
			UPDATE products SET name = $1, cost = $2, quantity = $3
			WHERE id = $4 
			RETURNING name, cost, quantity`

	args := []interface{}{product.Name, product.Cost, product.AmountAvailable, product.ID}

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	err := repo.DB.QueryRowContext(ctx, query, args...).Scan(&product.Name, &product.Cost, &product.AmountAvailable)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.ErrRecordNotFound
		case err.Error() == `pq: duplicate key value violates unique constraint "products_name_seller_id_key"`:
			return data.ErrDuplicateProductName
		default:
			return err
		}
	}

	return nil
}

func (repo *productRepository) Delete(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}

	query := `DELETE FROM products WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	result, err := repo.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}

func (repo *productRepository) GetAll(r dto.ListProductRequest) ([]*data.Product, data.Metadata, error) {

	filters := r.Filters
	query := fmt.Sprintf(`
			SELECT count(*) OVER(), id, name, cost, quantity, seller_id, created_at
			FROM products
			WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
			ORDER BY %s %s, id ASC 
			LIMIT $2  OFFSET $3`, filters.SortColumn(), filters.SortDirection(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), repository.QueryTimeout)
	defer cancel()

	args := []interface{}{r.Name, filters.Limit(), filters.Offset()}

	rows, err := repo.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, data.Metadata{}, err
	}

	totalRecords := 0
	var products []*data.Product

	for rows.Next() {
		var product data.Product

		err = rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Name,
			&product.Cost,
			&product.AmountAvailable,
			&product.Seller.ID,
			&product.CreatedAt,
		)

		if err != nil {
			return nil, data.Metadata{}, err
		}

		products = append(products, &product)
	}

	metadata := data.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
}
