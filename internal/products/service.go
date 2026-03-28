package products

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreateProductRequest) (Product, error) {
	if err := req.Validate(); err != nil {
		return Product{}, err
	}
	return s.repo.Create(ctx, req)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
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
