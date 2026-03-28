package products

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, req CreateProductRequest) (Product, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, params ListProductsParams) ([]Product, int64, error)
}

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	var p Product
	err := r.pool.QueryRow(ctx,
		`INSERT INTO products (name, description, price)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, description, price, created_at`,
		req.Name, req.Description, req.Price,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt)
	if err != nil {
		return Product{}, fmt.Errorf("create product: %w", err)
	}
	return p, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) List(ctx context.Context, params ListProductsParams) ([]Product, int64, error) {
	var total int64
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	offset := (params.Page - 1) * params.PageSize
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, price, created_at
		 FROM products
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		params.PageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, total, nil
}
