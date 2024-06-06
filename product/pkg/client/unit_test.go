package services_test

import (
	"context"
	"net/http"
	"testing"

	"product/pkg/db"
	"product/pkg/models"
	pb "product/pkg/pb"
	"product/pkg/services"
)

func TestCreateProduct(t *testing.T) {
	// Создаем временную базу данных для тестов.
	// В реальной среде тесты могут использовать тестовую базу данных.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем данные для нового продукта.
	req := &pb.CreateProductRequest{
		Name:  "Test Product",
		Stock: 10,
		Price: 100,
	}

	// Создаем продукт с использованием CreateProduct службы.
	resp, err := server.CreateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Проверяем успешность создания продукта.
	if resp.Status != http.StatusCreated {
		t.Errorf("Expected status %d, but got %d", http.StatusCreated, resp.Status)
	}

	// Проверяем наличие ID в ответе.
	if resp.Id == 0 {
		t.Error("Expected non-zero product ID")
	}

	// Проверяем, что продукт был добавлен в базу данных.
	var product models.Product
	if err := tmpDB.DB.First(&product, resp.Id).Error; err != nil {
		t.Errorf("Expected product with ID %d to exist in database, but got error: %v", resp.Id, err)
	}

	// Проверяем соответствие данных продукта.
	if product.Name != req.Name || product.Stock != req.Stock || product.Price != req.Price {
		t.Errorf("Expected product with name %s, stock %d, and price %d, but got %+v", req.Name, req.Stock, req.Price, product)
	}
}
func TestFindOneProduct(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Создаем тестовый продукт в базе данных.
	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	tmpDB.DB.Create(&testProduct)

	// Задаем ID тестового продукта.
	req := &pb.FindOneRequest{Id: testProduct.Id}

	// Находим продукт по ID с использованием FindOne службы.
	resp, err := server.FindOne(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to find product: %v", err)
	}

	// Проверяем успешность поиска продукта.
	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем соответствие данных найденного продукта с тестовыми данными.
	if resp.Data.Id != testProduct.Id || resp.Data.Name != testProduct.Name || resp.Data.Stock != testProduct.Stock || resp.Data.Price != testProduct.Price {
		t.Errorf("Expected product %+v, but got %+v", testProduct, resp.Data)
	}
}
func TestUpdateProductInvalidID(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем данные для обновления продукта с недопустимым ID.
	req := &pb.UpdateProductRequest{
		Id:    999, // Недопустимый ID
		Name:  "Updated Product",
		Stock: 20,
		Price: 200,
	}

	// Пытаемся обновить продукт с недопустимым ID.
	resp, err := server.UpdateProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	// Проверяем, что была получена ошибка.
	if resp.Status != http.StatusNotFound {
		t.Errorf("Expected status %d, but got %d", http.StatusNotFound, resp.Status)
	}
}

// TestDeleteProductInvalidID проверяет функцию DeleteProduct на обработку недопустимого ID продукта.
func TestDeleteProductInvalidID(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем ID для удаления, которого нет в базе данных.
	req := &pb.DeleteProductRequest{
		Id: 999, // Недопустимый ID
	}

	// Пытаемся удалить продукт с недопустимым ID.
	resp, err := server.DeleteProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	// Проверяем, что была получена ошибка.
	if resp.Status != http.StatusNotFound {
		t.Errorf("Expected status %d, but got %d", http.StatusNotFound, resp.Status)
	}
}
func TestListProductsPagination(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем запрос для списка продуктов с пагинацией.
	req := &pb.ListProductsRequest{
		Page:     1,
		PageSize: 10,
	}

	// Получаем список продуктов с использованием ListProducts службы.
	resp, err := server.ListProducts(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	// Проверяем успешность получения списка продуктов.
	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем корректность пагинации.
	// В данном тесте предполагается, что у вас есть данные в базе данных.
	// Проверьте, что количество элементов в ответе соответствует запрошенной странице и размеру страницы.
	// Проверьте, что данные продуктов в ответе соответствуют ожидаемым данным.
}

func TestListProductsFiltering(t *testing.T) {
	// Создаем временную базу данных для тестов.
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")

	// Инициализируем обработчик базы данных.
	handler := db.Handler{DB: tmpDB.DB}

	// Создаем новый сервер службы с использованием временной базы данных.
	server := services.Server{H: handler}

	// Задаем запрос для списка продуктов с фильтрацией по имени.
	req := &pb.ListProductsRequest{
		Filter: "Test",
	}

	// Получаем список продуктов с использованием ListProducts службы с фильтрацией.
	resp, err := server.ListProducts(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	// Проверяем успешность получения списка продуктов.
	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем корректность фильтрации.
	// В данном тесте предполагается, что у вас есть данные в базе данных.
	// Проверьте, что данные продуктов в ответе соответствуют ожидаемым данным.
	// Например, если у вас есть продукты с именами, содержащими "Test", то убедитесь, что они все присутствуют в ответе.
}
func TestDecreaseStock(t *testing.T) {
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")
	handler := db.Handler{DB: tmpDB.DB}
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	tmpDB.DB.Create(&testProduct)

	req := &pb.DecreaseStockRequest{
		Id:      testProduct.Id,
		OrderId: 1,
	}

	resp, err := server.DecreaseStock(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to decrease stock: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем, что запас был уменьшен на 1.
	var updatedProduct models.Product
	if err := tmpDB.DB.First(&updatedProduct, testProduct.Id).Error; err != nil {
		t.Fatalf("Failed to retrieve updated product: %v", err)
	}

	if updatedProduct.Stock != testProduct.Stock-1 {
		t.Errorf("Expected stock to be %d, but got %d", testProduct.Stock-1, updatedProduct.Stock)
	}
}

func TestUpdateProduct(t *testing.T) {
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")
	handler := db.Handler{DB: tmpDB.DB}
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	tmpDB.DB.Create(&testProduct)

	req := &pb.UpdateProductRequest{
		Id:    testProduct.Id,
		Name:  "Updated Test Product",
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

	// Проверяем, что продукт был обновлен в базе данных.
	var updatedProduct models.Product
	if err := tmpDB.DB.First(&updatedProduct, testProduct.Id).Error; err != nil {
		t.Fatalf("Failed to retrieve updated product: %v", err)
	}

	if updatedProduct.Name != req.Name || updatedProduct.Stock != req.Stock || updatedProduct.Price != req.Price {
		t.Errorf("Expected updated product to match request %+v, but got %+v", req, updatedProduct)
	}
}

func TestReadProduct(t *testing.T) {
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")
	handler := db.Handler{DB: tmpDB.DB}
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	tmpDB.DB.Create(&testProduct)

	req := &pb.ReadProductRequest{Id: testProduct.Id}

	resp, err := server.ReadProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to read product: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем, что полученные данные соответствуют тестовому продукту.
	if resp.Product == nil {
		t.Error("Expected product data, but got nil")
	}

	if resp.Product.Id != testProduct.Id || resp.Product.Name != testProduct.Name || resp.Product.Stock != testProduct.Stock || resp.Product.Price != testProduct.Price {
		t.Errorf("Expected product %+v, but got %+v", testProduct, resp.Product)
	}
}

func TestDeleteProduct(t *testing.T) {
	tmpDB := db.Init("postgres://postgres:samsungtab7@localhost:5432/product-service")
	handler := db.Handler{DB: tmpDB.DB}
	server := services.Server{H: handler}

	testProduct := models.Product{Name: "Test Product", Stock: 10, Price: 100}
	tmpDB.DB.Create(&testProduct)

	req := &pb.DeleteProductRequest{Id: testProduct.Id}

	resp, err := server.DeleteProduct(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	if resp.Status != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, resp.Status)
	}

	// Проверяем, что продукт был удален из базы данных.
	var deletedProduct models.Product
	if err := tmpDB.DB.First(&deletedProduct, testProduct.Id).Error; err == nil {
		t.Error("Expected product to be deleted, but found in database")
	}
}
