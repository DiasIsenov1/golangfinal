package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"order-service/pkg/pb"
)

// Mock ProductServiceClient for testing
type mockProductServiceClient struct {
	pb.UnimplementedProductServiceServer
}

func (m *mockProductServiceClient) FindOne(ctx context.Context, req *pb.FindOneRequest, opts ...grpc.CallOption) (*pb.FindOneResponse, error) {
	// Mock response based on the request id
	if req.Id == 1 {
		return &pb.FindOneResponse{
			Status: 1,
			Error:  "",
			Data: &pb.FindOneData{
				Id:    req.Id,
				Name:  "Test Product",
				Stock: 10,
				Price: 100,
			},
		}, nil
	}
	// Mock product not found response
	return &pb.FindOneResponse{
		Status: 0,
		Error:  "Product not found",
	}, nil
}

func (m *mockProductServiceClient) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest, opts ...grpc.CallOption) (*pb.DecreaseStockResponse, error) {
	// Mock response based on the request id
	if req.Id == 1 {
		return &pb.DecreaseStockResponse{
			Status: 1,
			Error:  "",
		}, nil
	}
	// Mock decrease stock failure response
	return &pb.DecreaseStockResponse{
		Status: 0,
		Error:  "Failed to decrease stock",
	}, nil
}

func (m *mockProductServiceClient) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, opts ...grpc.CallOption) (*pb.CreateProductResponse, error) {
	// Mock response
	return &pb.CreateProductResponse{
		Status: 1,
		Error:  "",
		Id:     1,
	}, nil
}

func TestFindOne(t *testing.T) {
	mockClient := &mockProductServiceClient{}
	productClient := ProductServiceClient{Client: mockClient}

	res, err := productClient.FindOne(1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Data.Id)
	assert.Equal(t, "Test Product", res.Data.Name)
	assert.Equal(t, int64(10), res.Data.Stock)
	assert.Equal(t, int64(100), res.Data.Price)
}

func TestDecreaseStock(t *testing.T) {
	mockClient := &mockProductServiceClient{}
	productClient := ProductServiceClient{Client: mockClient}

	res, err := productClient.DecreaseStock(1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Status)
	assert.Empty(t, res.Error)
}

func TestCreateProduct(t *testing.T) {
	mockClient := &mockProductServiceClient{}
	productClient := ProductServiceClient{Client: mockClient}

	req := &pb.CreateProductRequest{
		Name:  "New Product",
		Stock: 20,
		Price: 200,
	}

	res, err := productClient.CreateProduct(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Status)
	assert.Empty(t, res.Error)
	assert.Equal(t, int64(1), res.Id)
}

func TestFindOne_NotFoun(t *testing.T) {
	mockClient := &mockProductServiceClient{}
	productClient := ProductServiceClient{Client: mockClient}

	res, err := productClient.FindOne(2)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), res.Status)
	assert.Equal(t, "Product not found", res.Error)
}

func TestDecreaseStock_Failure(t *testing.T) {
	mockClient := &mockProductServiceClient{}
	productClient := ProductServiceClient{Client: mockClient}

	res, err := productClient.DecreaseStock(2, 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), res.Status)
	assert.Equal(t, "Failed to decrease stock", res.Error)
}
