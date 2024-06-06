package client

import (
	"context"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"

	"github.com/stretchr/testify/mock"
	"order-service/pkg/pb"
)

type MockProductServiceClient struct {
	mock.Mock
}

func (m *MockProductServiceClient) CreateProduct(ctx context.Context, in *pb.CreateProductRequest, opts ...grpc.CallOption) (*pb.CreateProductResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CreateProductResponse), args.Error(1)
}

func (m *MockProductServiceClient) FindOne(ctx context.Context, in *pb.FindOneRequest, opts ...grpc.CallOption) (*pb.FindOneResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.FindOneResponse), args.Error(1)
}

func (m *MockProductServiceClient) DecreaseStock(ctx context.Context, in *pb.DecreaseStockRequest, opts ...grpc.CallOption) (*pb.DecreaseStockResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.DecreaseStockResponse), args.Error(1)
}

func TestFindOne_Success(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	serviceClient := ProductServiceClient{Client: mockClient}

	productID := int64(1)
	response := &pb.FindOneResponse{
		Status: 1,
		Data: &pb.FindOneData{
			Id:    productID,
			Name:  "Test Product",
			Stock: 100,
			Price: 50,
		},
	}

	mockClient.On("FindOne", context.Background(), &pb.FindOneRequest{Id: productID}).Return(response, nil)

	result, err := serviceClient.FindOne(productID)
	require.NoError(t, err)
	require.Equal(t, response, result)

	mockClient.AssertExpectations(t)
}

func TestFindOne_NotFound(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	serviceClient := ProductServiceClient{Client: mockClient}

	productID := int64(999)
	response := &pb.FindOneResponse{
		Status: 0,
		Error:  "Product not found",
	}

	mockClient.On("FindOne", context.Background(), &pb.FindOneRequest{Id: productID}).Return(response, nil)

	result, err := serviceClient.FindOne(productID)
	require.NoError(t, err)
	require.Equal(t, response, result)

	mockClient.AssertExpectations(t)
}

func TestDecreaseStock_Success(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	serviceClient := ProductServiceClient{Client: mockClient}

	productID := int64(1)
	orderID := int64(123)
	response := &pb.DecreaseStockResponse{Status: 1}

	mockClient.On("DecreaseStock", context.Background(), &pb.DecreaseStockRequest{Id: productID, OrderId: orderID}).Return(response, nil)

	result, err := serviceClient.DecreaseStock(productID, orderID)
	require.NoError(t, err)
	require.Equal(t, response, result)

	mockClient.AssertExpectations(t)
}

func TestDecreaseStock_OutOfStock(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	serviceClient := ProductServiceClient{Client: mockClient}

	productID := int64(1)
	orderID := int64(123)
	response := &pb.DecreaseStockResponse{
		Status: 0,
		Error:  "Out of stock",
	}

	mockClient.On("DecreaseStock", context.Background(), &pb.DecreaseStockRequest{Id: productID, OrderId: orderID}).Return(response, nil)

	result, err := serviceClient.DecreaseStock(productID, orderID)
	require.NoError(t, err)
	require.Equal(t, response, result)

	mockClient.AssertExpectations(t)
}

func TestCreateProduct_InvalidRequest(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	serviceClient := ProductServiceClient{Client: mockClient}

	request := &pb.CreateProductRequest{
		Name:  "",
		Stock: -1,
		Price: -1,
	}
	response := &pb.CreateProductResponse{
		Status: 0,
		Error:  "Invalid request",
	}

	mockClient.On("CreateProduct", context.Background(), request).Return(response, nil)

	result, err := serviceClient.CreateProduct(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, response, result)

	mockClient.AssertExpectations(t)
}

func TestFindOne_NoData(t *testing.T) {
	mockClient := new(MockProductServiceClient)
	productID := int64(1)
	response := &pb.FindOneResponse{
		Status: 1,
		Data:   nil,
	}

	mockClient.On("FindOne", context.Background(), &pb.FindOneRequest{Id: productID}).Return(response, nil)

}
