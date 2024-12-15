package components

import (
	"context"
	"database/sql"
	"fmt"
	"service-weaver-app/models"
)

type OrderProcessingImpl struct {
	inventory InventoryManagement
	db        *sql.DB
}

func NewOrderProcessing(inventory InventoryManagement, db *sql.DB) *OrderProcessingImpl {
	return &OrderProcessingImpl{inventory: inventory, db: db}
}

func (op *OrderProcessingImpl) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {
	stock, err := op.inventory.CheckStock(ctx, order.ProductID)
	if err != nil {
		return models.Order{}, err
	}
	if stock < order.Quantity {
		return models.Order{}, fmt.Errorf("insufficient stock")
	}

	product, err := op.inventory.GetProduct(ctx, order.ProductID)
	if err != nil {
		return models.Order{}, err
	}

	order.Total = float64(order.Quantity) * product.Price

	err = op.inventory.UpdateStock(ctx, order.ProductID, -order.Quantity)
	if err != nil {
		return models.Order{}, err
	}

	query := `INSERT INTO orders (product_id, quantity, total, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err = op.db.QueryRowContext(ctx, query, order.ProductID, order.Quantity, order.Total, "Pending").Scan(&order.ID)
	if err != nil {
		return models.Order{}, fmt.Errorf("could not create order: %w", err)
	}

	return order, nil
}

func (op *OrderProcessingImpl) GetOrders(ctx context.Context) ([]models.Order, error) {
	query := `SELECT id, product_id, quantity, total, status FROM orders`
	rows, err := op.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch orders: %w", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.ProductID, &order.Quantity, &order.Total, &order.Status); err != nil {
			return nil, fmt.Errorf("could not scan order: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (op *OrderProcessingImpl) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := op.db.ExecContext(ctx, query, status, orderID)
	return err
}
