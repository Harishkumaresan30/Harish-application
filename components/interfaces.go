package components

import (
	"context"
	"service-weaver-app/models"
)

// InventoryManagement defines methods for inventory operations.
// InventoryManagement defines the interface for inventory operations.
type InventoryManagement interface {
	AddProduct(ctx context.Context, product models.Product) error
	UpdateStock(ctx context.Context, productID string, quantity int) error
	CheckStock(ctx context.Context, productID string) (int, error)
	GetProduct(ctx context.Context, productID string) (models.Product, error)
	GetProducts(ctx context.Context) ([]models.Product, error)
}

type OrderProcessing interface {
	CreateOrder(ctx context.Context, order models.Order) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, status string) error
	GetOrders(ctx context.Context) ([]models.Order, error)
}

// Analytics defines methods for tracking metrics.
type Analytics interface {
	TrackMetric(ctx context.Context, name string, value float64) error
	GetMetrics(ctx context.Context) ([]models.Metric, error)
}
