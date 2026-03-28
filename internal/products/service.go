package products

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

type EventPublisher interface {
	Publish(ctx context.Context, key, value []byte) error
}

type ProductEvent struct {
	Type    string   `json:"type"`
	Product *Product `json:"product,omitempty"`
	ID      int64    `json:"id,omitempty"`
}

type Service struct {
	repo      Repository
	publisher EventPublisher
}

func NewService(repo Repository, publisher EventPublisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	if err := req.Validate(); err != nil {
		return Product{}, err
	}

	product, err := s.repo.Create(ctx, req)
	if err != nil {
		return Product{}, err
	}

	s.publishEvent(ctx, ProductEvent{Type: "product_created", Product: &product})
	return product, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.publishEvent(ctx, ProductEvent{Type: "product_deleted", ID: id})
	return nil
}

func (s *Service) List(ctx context.Context, params ListProductsParams) (ListProductsResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	products, total, err := s.repo.List(ctx, params)
	if err != nil {
		return ListProductsResponse{}, err
	}

	return ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (s *Service) publishEvent(ctx context.Context, event ProductEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)
		return
	}

	key := []byte(fmt.Sprintf("%s:%d", event.Type, event.ID))
	if event.Product != nil {
		key = []byte(fmt.Sprintf("%s:%d", event.Type, event.Product.ID))
	}

	if err := s.publisher.Publish(ctx, key, data); err != nil {
		log.Printf("failed to publish event: %v", err)
	}
}
