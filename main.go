package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"service-weaver-app/components"
	"service-weaver-app/database"
	"service-weaver-app/models"
)

var inventory components.InventoryManagement
var orders components.OrderProcessing

func main() {
	// Initialize the database connection
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize components with database
	inventory = components.NewInventoryManagement(database.DB)
	orders = components.NewOrderProcessing(inventory, database.DB)

	// Define routes
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/add-product-form", addProductFormHandler)
	http.HandleFunc("/create-order-form", createOrderFormHandler)
	http.HandleFunc("/view-products", viewProductsHandler)
	http.HandleFunc("/view-orders", viewOrdersHandler)
	http.HandleFunc("/add-product", addProductHandler)
	http.HandleFunc("/create-order", createOrderHandler)


	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handlers
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	tmpl.Execute(w, nil)
}

// Add product handler for processing product addition
func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var product models.Product

	// Handle form data
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		product.ID = r.FormValue("id")
		product.Name = r.FormValue("name")
		product.Stock, err = strconv.Atoi(r.FormValue("stock"))
		if err != nil {
			http.Error(w, "Invalid stock value", http.StatusBadRequest)
			return
		}
		product.Price, err = strconv.ParseFloat(r.FormValue("price"), 64)
		if err != nil {
			http.Error(w, "Invalid price value", http.StatusBadRequest)
			return
		}
	} else {
		// Handle JSON payload
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Add product to the database via the inventory component
	if err := inventory.AddProduct(r.Context(), product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Product added successfully"})
}

// Create order handler for processing orders
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order

	// Handle form data
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		order.ProductID = r.FormValue("product_id")
		order.Quantity, err = strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			http.Error(w, "Invalid quantity value", http.StatusBadRequest)
			return
		}
	} else {
		// Handle JSON payload
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Create the order via the orders component
	createdOrder, err := orders.CreateOrder(r.Context(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdOrder)
}

func viewProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	products, err := inventory.GetProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch products: "+err.Error(), http.StatusInternalServerError)
		return
	}

	viewProductsTemplate.Execute(w, products)
}

func viewOrdersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	allOrders, err := orders.GetOrders(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch orders: "+err.Error(), http.StatusInternalServerError)
		return
	}

	viewOrdersTemplate.Execute(w, allOrders)
}

// Add product form handler
func addProductFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	addProductFormTemplate.Execute(w, nil)
}

// Create order form handler
func createOrderFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	createOrderFormTemplate.Execute(w, nil)
}

var viewOrdersTemplate = template.Must(template.New("viewOrders").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <title>View Orders</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <h1>Order List</h1>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>Order ID</th>
                    <th>Product ID</th>
                    <th>Quantity</th>
                    <th>Total</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.ProductID}}</td>
                    <td>{{.Quantity}}</td>
                    <td>{{.Total}}</td>
                    <td>{{.Status}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>
`))
var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Inventory Management System</title>
    <!-- Bootstrap CSS -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
        <div class="container-fluid">
            <a class="navbar-brand" href="/">Inventory Management</a>
        </div>
    </nav>
    <div class="container mt-4">
        <h1>Welcome to the Inventory Management System</h1>
        <p>Use the options below to navigate and interact with the system.</p>
        <div class="row">
            <div class="col-md-3">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Add Product</h5>
                        <p class="card-text">Add a new product to the inventory.</p>
                        <a href="/add-product-form" class="btn btn-primary">Add Product</a>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">Create Order</h5>
                        <p class="card-text">Place a new order for a product.</p>
                        <a href="/create-order-form" class="btn btn-primary">Create Order</a>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">View Products</h5>
                        <p class="card-text">Check the list of available products.</p>
                        <a href="/view-products" class="btn btn-primary">View Products</a>
                    </div>
                </div>
            </div>
            <div class="col-md-3">
                <div class="card">
                    <div class="card-body">
                        <h5 class="card-title">View Orders</h5>
                        <p class="card-text">Check the list of all orders.</p>
                        <a href="/view-orders" class="btn btn-primary">View Orders</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/js/bootstrap.bundle.min.js"></script>
</body>
</html>
`))

var addProductFormTemplate = template.Must(template.New("addProductForm").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Add Product</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <h1>Add Product</h1>
        <form action="/add-product" method="POST">
            <div class="mb-3">
                <label for="productID" class="form-label">Product ID</label>
                <input type="text" class="form-control" id="productID" name="id" required>
            </div>
            <div class="mb-3">
                <label for="productName" class="form-label">Product Name</label>
                <input type="text" class="form-control" id="productName" name="name" required>
            </div>
            <div class="mb-3">
                <label for="productStock" class="form-label">Stock</label>
                <input type="number" class="form-control" id="productStock" name="stock" required>
            </div>
            <div class="mb-3">
                <label for="productPrice" class="form-label">Price</label>
                <input type="number" step="0.01" class="form-control" id="productPrice" name="price" required>
            </div>
            <button type="submit" class="btn btn-primary">Add Product</button>
        </form>
    </div>
</body>
</html>
`))

var createOrderFormTemplate = template.Must(template.New("createOrderForm").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <title>Create Order</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <h1>Create Order</h1>
        <form action="/create-order" method="POST">
            <div class="mb-3">
                <label for="productID" class="form-label">Product ID</label>
                <input type="text" class="form-control" id="productID" name="product_id" required>
            </div>
            <div class="mb-3">
                <label for="quantity" class="form-label">Quantity</label>
                <input type="number" class="form-control" id="quantity" name="quantity" required>
            </div>
            <button type="submit" class="btn btn-primary">Create Order</button>
        </form>
    </div>
</body>
</html>
`))

var viewProductsTemplate = template.Must(template.New("viewProducts").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <title>View Products</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mt-4">
        <h1>Product List</h1>
        <table class="table table-striped">
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Stock</th>
                    <th>Price</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td>{{.ID}}</td>
                    <td>{{.Name}}</td>
                    <td>{{.Stock}}</td>
                    <td>{{.Price}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
</html>
`))
