package products

import (
	"context"
	"sync"
)

type mockRepository struct {
	mu       sync.Mutex
	products []Product
	nextID   int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{nextID: 1}
}

func (m *mockRepository) Create(_ context.Context, req CreateProductRequest) (Product, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p := Product{
		ID:          m.nextID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}
	m.nextID++
	m.products = append(m.products, p)
	return p, nil
}

func (m *mockRepository) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.products {
		if p.ID == id {
			m.products = append(m.products[:i], m.products[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (m *mockRepository) List(_ context.Context, params ListProductsParams) ([]Product, int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	total := int64(len(m.products))
	start := (params.Page - 1) * params.PageSize
	if start >= int(total) {
		return []Product{}, total, nil
	}

	end := start + params.PageSize
	if end > int(total) {
		end = int(total)
	}

	result := make([]Product, end-start)
	copy(result, m.products[start:end])
	return result, total, nil
}

type mockPublisher struct {
	mu       sync.Mutex
	messages []publishedMessage
}

type publishedMessage struct {
	key   []byte
	value []byte
}

func newMockPublisher() *mockPublisher {
	return &mockPublisher{}
}

func (m *mockPublisher) Publish(_ context.Context, key, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, publishedMessage{key: key, value: value})
	return nil
}

func (m *mockPublisher) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.messages)
}
