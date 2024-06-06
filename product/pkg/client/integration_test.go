package services_test

import (
	"context"
	"net/http"
	"product/pkg/db"
	"testing"

	"product/pkg/models"
	pb "product/pkg/pb"
	"product/pkg/services"
)

var handler db.Handler

func setup() {
	// Создаем соединение с вашей базой данных для тестов
	dbURL := "postgres://postgres:samsungtab7@localhost:5432/test_db"
	handler = db.Init(dbURL)

	// Автоматическое мигрирование моделей
	handler.DB.AutoMigrate(&models.Product{})
	handler.DB.AutoMigrate(&models.StockDecreaseLog{})
}

func teardown() {
	// Закрываем соединение с базой данных после завершения тестов
	handler.Close()
}

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func TestIntegrationCreateProduct(t *testing.T) {
	server := services.Server{H: handler}

	req := &pb.CreateProductRequest{
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	resp, err := server.CreateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if resp.Status != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, resp.Status)
	}

	if resp.Id == 0 {
		t.Error("Expected non-zero product ID")
	}

	var product models.Product
	if err := handler.DB.First(&product, resp.Id).Error; err != nil {
		t.Errorf("Expected product with ID %d to exist in database, but got error: %v", resp.Id, err)
	}

	if product.Name != req.Name || product.Stock != req.Stock || product.Price != req.Price {
		t.Errorf("Expected product with name %s, stock %d, and price %d, but got %+v", req.Name, req.Stock, req.Price, product)
	}
}

func TestIntegrationFindOneProduct(t *testing.T) {
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	handler.DB.Create(&testProduct)

	req := &pb.FindOneRequest{Id: testProduct.Id}

	resp, err := server.FindOne(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to find product: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	if resp.Data.Id != testProduct.Id || resp.Data.Name != testProduct.Name || resp.Data.Stock != testProduct.Stock || resp.Data.Price != testProduct.Price {
		t.Errorf("Expected product %+v, but got %+v", testProduct, resp.Data)
	}
}

func TestIntegrationUpdateProduct(t *testing.T) {
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	handler.DB.Create(&testProduct)

	req := &pb.UpdateProductRequest{
		Id:    testProduct.Id,
		Name:  "Updated Product",
		Stock: 20,
		Price: 200,
	}

	resp, err := server.UpdateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	var updatedProduct models.Product
	if err := handler.DB.First(&updatedProduct, testProduct.Id).Error; err != nil {
		t.Errorf("Failed to find updated product in database: %v", err)
	}

	if updatedProduct.Name != req.Name || updatedProduct.Stock != req.Stock || updatedProduct.Price != req.Price {
		t.Errorf("Expected updated product with name %s, stock %d, and price %d, but got %+v", req.Name, req.Stock, req.Price, updatedProduct)
	}
}

func TestIntegrationDeleteProduct(t *testing.T) {
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	handler.DB.Create(&testProduct)

	req := &pb.DeleteProductRequest{Id: testProduct.Id}

	resp, err := server.DeleteProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	var deletedProduct models.Product
	if err := handler.DB.First(&deletedProduct, testProduct.Id).Error; err == nil {
		t.Error("Expected product to be deleted, but found in database")
	}
}

func TestIntegrationListProducts(t *testing.T) {
	server := services.Server{H: handler}

	req := &pb.ListProductsRequest{}

	resp, err := server.ListProducts(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	if len(resp.Products) == 0 {
		t.Error("Expected non-empty product list")
	}
}
