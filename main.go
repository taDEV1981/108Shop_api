package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

type Product struct {
	ID           string `json:"id"`
	Category     string `json:"category"`
	Sellers      string `json:"sellers"`
	Name         string `json:"name"`
	CodeInternal string `json:"code_internal"`
	CodeExternal string `json:"code_external"`
	Type         string `json:"type"`
	Groups       string `json:"groups"`
	Details      string `json:"details"`
	Substitute   string `json:"substitute"`
	Assembleds   string `json:"assembleds"`
}

func main() {
	server := "eightzeroonefamr.database.windows.net"
	port := 1433
	user := "eightzeroone_famr"
	password := "data108@"
	database := "eightzeroone_famr"

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	defer db.Close() // Close the database connection when main() exits

	// Check if the connection is successful
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}

	fmt.Println("Connected to the database!")

	// Create a Gin router with default middleware stack
	router := gin.Default()

	// Add CORS middleware to allow requests from any origin
	router.Use(cors.Default())

	// Define your API routes
	router.GET("/products", getProducts)
	router.POST("/products", addProduct)
	router.GET("/products/:id", getProductByID)
	router.PUT("/products/:id", updateProduct)
	router.DELETE("/products/:id", deleteProductByID)

	// Start the HTTP server
	router.Run(":8080")
}

func getProducts(c *gin.Context) {
	var products []Product
	ctx := c.Request.Context()
	rows, err := db.QueryContext(ctx, "SELECT * FROM dbo.products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Category, &product.Sellers, &product.Name, &product.CodeInternal, &product.CodeExternal, &product.Type, &product.Groups, &product.Details, &product.Substitute, &product.Assembleds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, product)
	}
	c.JSON(http.StatusOK, products)
}

func getProductByID(c *gin.Context) {
	paramID := c.Param("id")
	var product Product
	ctx := c.Request.Context()
	err := db.QueryRowContext(ctx, "SELECT * FROM dbo.products WHERE product_id = ?", paramID).Scan(&product.ID, &product.Category, &product.Sellers, &product.Name, &product.CodeInternal, &product.CodeExternal, &product.Type, &product.Groups, &product.Details, &product.Substitute, &product.Assembleds)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
func addProduct(c *gin.Context) {
	var newProduct Product

	// Parse form data
	category := c.PostForm("category")
	sellers := c.PostForm("sellers")
	name := c.PostForm("name")
	codeInternal := c.PostForm("code_internal")
	codeExternal := c.PostForm("code_external")
	productType := c.PostForm("type")
	groups := c.PostForm("groups")
	details := c.PostForm("details")
	substitute := c.PostForm("substitute")
	assembleds := c.PostForm("assembleds")

	// Set the parsed form data to the newProduct struct
	newProduct.Category = category
	newProduct.Sellers = sellers
	newProduct.Name = name
	newProduct.CodeInternal = codeInternal
	newProduct.CodeExternal = codeExternal
	newProduct.Type = productType
	newProduct.Groups = groups
	newProduct.Details = details
	newProduct.Substitute = substitute
	newProduct.Assembleds = assembleds

	// Insert the new product into the database
	ctx := c.Request.Context()
	_, err := db.ExecContext(ctx, "INSERT INTO dbo.products (product_category, sellers, product_name, product_code_i, product_code_ex, product_type, product_groups, product_details, product_substitute, product_assembleds) VALUES (@category, @sellers, @name, @code_internal, @code_external, @type, @groups, @details, @substitute, @assembleds)",
		sql.Named("category", category),
		sql.Named("sellers", sellers),
		sql.Named("name", name),
		sql.Named("code_internal", codeInternal),
		sql.Named("code_external", codeExternal),
		sql.Named("type", productType),
		sql.Named("groups", groups),
		sql.Named("details", details),
		sql.Named("substitute", substitute),
		sql.Named("assembleds", assembleds),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
}

// func addProduct(c *gin.Context) {
// 	var newProduct Product
// 	if err := c.BindJSON(&newProduct); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// You can remove the product_id field from the insert statement
// 	ctx := c.Request.Context()
// 	_, err := db.ExecContext(ctx, "INSERT INTO dbo.products (product_category, sellers, product_name, product_code_i, product_code_ex, product_type, product_groups, product_details, product_substitute, product_assembleds) VALUES (@category, @sellers, @name, @code_internal, @code_external, @type, @groups, @details, @substitute, @assembleds)",
// 		sql.Named("category", newProduct.Category),
// 		sql.Named("sellers", newProduct.Sellers),
// 		sql.Named("name", newProduct.Name),
// 		sql.Named("code_internal", newProduct.CodeInternal),
// 		sql.Named("code_external", newProduct.CodeExternal),
// 		sql.Named("type", newProduct.Type),
// 		sql.Named("groups", newProduct.Groups),
// 		sql.Named("details", newProduct.Details),
// 		sql.Named("substitute", newProduct.Substitute),
// 		sql.Named("assembleds", newProduct.Assembleds),
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, newProduct)
// }

func updateProduct(c *gin.Context) {
	var updatedProduct Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	paramID := c.Param("id")
	ctx := c.Request.Context()
	_, err := db.ExecContext(ctx, "UPDATE dbo.products SET product_category=@category, sellers=@sellers, product_name=@name, product_code_i=@code_internal, product_code_ex=@code_external, product_type=@type, product_groups=@groups, product_details=@details, product_substitute=@substitute, product_assembleds=@assembleds WHERE product_id=@id",
		sql.Named("category", updatedProduct.Category),
		sql.Named("sellers", updatedProduct.Sellers),
		sql.Named("name", updatedProduct.Name),
		sql.Named("code_internal", updatedProduct.CodeInternal),
		sql.Named("code_external", updatedProduct.CodeExternal),
		sql.Named("type", updatedProduct.Type),
		sql.Named("groups", updatedProduct.Groups),
		sql.Named("details", updatedProduct.Details),
		sql.Named("substitute", updatedProduct.Substitute),
		sql.Named("assembleds", updatedProduct.Assembleds),
		sql.Named("id", paramID),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedProduct)
}

func deleteProductByID(c *gin.Context) {
	paramID := c.Param("id")
	ctx := c.Request.Context()
	_, err := db.ExecContext(ctx, "DELETE FROM dbo.products WHERE product_id=@id", sql.Named("id", paramID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
