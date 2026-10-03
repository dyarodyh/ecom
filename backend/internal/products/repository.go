package products

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (repo *Repository) List(ctx context.Context) ([]Product, error) {
	const query = `
	SELECT id, name, price, stock
	FROM products
	ORDER BY id`

	rows, err := repo.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query daftar produk: %w", err)
	}

	defer rows.Close()
	items := make([]Product, 0)

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
		)

		if err != nil {
			return nil, fmt.Errorf("membaca baris produk : %w", err)
		}

		items = append(items, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca daftar product : %w", err)

	}

	return items, nil
}

func (repo *Repository) Create(
	ctx context.Context,
	product Product,
) (Product, error) {
	const query = `
	INSERT INTO products (name, price, stock)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, stock`

	var createdProduct Product

	err := repo.pool.QueryRow(
		ctx,
		query,
		product.Name,
		product.Price,
		product.Stock,
	).Scan(
		&createdProduct.ID,
		&createdProduct.Name,
		&createdProduct.Price,
		&createdProduct.Stock,
	)

	if err != nil {
		return Product{}, fmt.Errorf("menambhakan produk: %w", err)
	}

	return createdProduct, nil
}
