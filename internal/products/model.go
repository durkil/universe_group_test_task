package products

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (r CreateProductRequest) Validate() error {
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}

type ListProductsParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListProductsResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}
