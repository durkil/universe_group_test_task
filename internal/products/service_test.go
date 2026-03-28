package products

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupService() (*Service, *mockRepository, *mockPublisher) {
	repo := newMockRepository()
	pub := newMockPublisher()
	svc := NewService(repo, pub)
	return svc, repo, pub
}

func TestService_Create(t *testing.T) {
	svc, _, pub := setupService()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		product, err := svc.Create(ctx, CreateProductRequest{
			Name:        "Test Product",
			Description: "A test product",
			Price:       29.99,
		})

		require.NoError(t, err)
		assert.Equal(t, int64(1), product.ID)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "A test product", product.Description)
		assert.Equal(t, 29.99, product.Price)
		assert.Equal(t, 1, pub.count())
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := svc.Create(ctx, CreateProductRequest{
			Name:  "",
			Price: 10.0,
		})

		assert.ErrorIs(t, err, ErrNameRequired)
	})

	t.Run("negative price", func(t *testing.T) {
		_, err := svc.Create(ctx, CreateProductRequest{
			Name:  "Product",
			Price: -5.0,
		})

		assert.ErrorIs(t, err, ErrInvalidPrice)
	})
}

func TestService_Delete(t *testing.T) {
	svc, _, pub := setupService()
	ctx := context.Background()

	product, err := svc.Create(ctx, CreateProductRequest{
		Name:  "To Delete",
		Price: 10.0,
	})
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := svc.Delete(ctx, product.ID)

		require.NoError(t, err)
		assert.Equal(t, 2, pub.count()) // 1 create + 1 delete
	})

	t.Run("not found", func(t *testing.T) {
		err := svc.Delete(ctx, 999)

		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestService_List(t *testing.T) {
	svc, _, _ := setupService()
	ctx := context.Background()

	for i := range 5 {
		_, err := svc.Create(ctx, CreateProductRequest{
			Name:  "Product " + string(rune('A'+i)),
			Price: float64(i+1) * 10,
		})
		require.NoError(t, err)
	}

	t.Run("first page", func(t *testing.T) {
		resp, err := svc.List(ctx, ListProductsParams{Page: 1, PageSize: 2})

		require.NoError(t, err)
		assert.Len(t, resp.Products, 2)
		assert.Equal(t, int64(5), resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 2, resp.PageSize)
	})

	t.Run("last page", func(t *testing.T) {
		resp, err := svc.List(ctx, ListProductsParams{Page: 3, PageSize: 2})

		require.NoError(t, err)
		assert.Len(t, resp.Products, 1)
		assert.Equal(t, int64(5), resp.Total)
	})

	t.Run("beyond last page", func(t *testing.T) {
		resp, err := svc.List(ctx, ListProductsParams{Page: 10, PageSize: 2})

		require.NoError(t, err)
		assert.Empty(t, resp.Products)
		assert.Equal(t, int64(5), resp.Total)
	})

	t.Run("default pagination", func(t *testing.T) {
		resp, err := svc.List(ctx, ListProductsParams{Page: 0, PageSize: 0})

		require.NoError(t, err)
		assert.Len(t, resp.Products, 5)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 20, resp.PageSize)
	})
}
