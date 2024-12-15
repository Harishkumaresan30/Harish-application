package components

import (
	"context"
	"database/sql"
	"fmt"
	"service-weaver-app/models"
)

// InventoryManagementImpl is the implementation of InventoryManagement.
type InventoryManagementImpl struct {
	db *sql.DB
}

// NewInventoryManagement initializes a new InventoryManagementImpl instance.
func NewInventoryManagement(db *sql.DB) *InventoryManagementImpl {
	return &InventoryManagementImpl{db: db}
}

// AddProduct adds a new product to the inventory.
func (im *InventoryManagementImpl) AddProduct(ctx context.Context, product models.Product) error {
	query := `INSERT INTO products (id, name, stock, price) VALUES ($1, $2, $3, $4)`
	_, err := im.db.ExecContext(ctx, query, product.ID, product.Name, product.Stock, product.Price)
	if err != nil {
		return fmt.Errorf("could not add product: %w", err)
	}
	return nil
}

// UpdateStock updates the stock level of an existing product.
func (im *InventoryManagementImpl) UpdateStock(ctx context.Context, productID string, quantity int) error {
	query := `UPDATE products SET stock = stock + $1 WHERE id = $2`
	_, err := im.db.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return fmt.Errorf("could not update stock: %w", err)
	}
	return nil
}

// CheckStock returns the current stock level of a product.
func (im *InventoryManagementImpl) CheckStock(ctx context.Context, productID string) (int, error) {
	var stock int
	query := `SELECT stock FROM products WHERE id = $1`
	err := im.db.QueryRowContext(ctx, query, productID).Scan(&stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("product not found")
		}
		return 0, fmt.Errorf("could not check stock: %w", err)
	}
	return stock, nil
}

// GetProduct retrieves the details of a specific product by its ID.
func (im *InventoryManagementImpl) GetProduct(ctx context.Context, productID string) (models.Product, error) {
	var product models.Product
	query := `SELECT id, name, stock, price FROM products WHERE id = $1`
	err := im.db.QueryRowContext(ctx, query, productID).Scan(&product.ID, &product.Name, &product.Stock, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Product{}, fmt.Errorf("product not found")
		}
		return models.Product{}, fmt.Errorf("could not fetch product: %w", err)
	}
	return product, nil
}

// GetProducts retrieves all products from the inventory.
func (im *InventoryManagementImpl) GetProducts(ctx context.Context) ([]models.Product, error) {
	query := `SELECT id, name, stock, price FROM products`
	rows, err := im.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not fetch products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Stock, &product.Price); err != nil {
			return nil, fmt.Errorf("could not scan product: %w", err)
		}
		products = append(products, product)
	}
	return products, nil
}
