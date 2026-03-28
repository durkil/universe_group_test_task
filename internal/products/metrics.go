package products

import "github.com/prometheus/client_golang/prometheus"

var (
	ProductsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_created_total",
		Help: "Total number of products created",
	})

	ProductsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "products_deleted_total",
		Help: "Total number of products deleted",
	})
)

func init() {
	prometheus.MustRegister(ProductsCreated, ProductsDeleted)
}
