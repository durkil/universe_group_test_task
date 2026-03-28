package products

import "errors"

var (
	ErrNameRequired = errors.New("product name is required")
	ErrInvalidPrice = errors.New("price must be non-negative")
	ErrNotFound     = errors.New("product not found")
)
