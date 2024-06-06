package service

import (
	"context"
	"order-service/pkg/client"
	"order-service/pkg/db"
	"order-service/pkg/models"
	"order-service/pkg/pb"
	"time"
)

type Server struct {
	H db.Handler
	pb.UnimplementedOrderServiceServer
	ProductSvc client.ProductServiceClient
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	order := models.Order{
		ProductID: req.ProductId,
		UserID:    req.UserId,
		Quantity:  req.Quantity,
		Price:     100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if result := s.H.DB.Create(&order); result.Error != nil {
		return &pb.CreateOrderResponse{
			Status: 500,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateOrderResponse{
		Status: 201,
		Id:     order.ID,
	}, nil
}

func (s *Server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	var order models.Order
	if result := s.H.DB.First(&order, req.Id); result.Error != nil {
		return &pb.GetOrderResponse{
			Status: 404,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.GetOrderResponse{
		Status: 200,
		Order: &pb.Order{
			Id:        order.ID,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Quantity:  order.Quantity,
			Price:     order.Price,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
			UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *Server) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.UpdateOrderResponse, error) {
	var order models.Order
	if result := s.H.DB.First(&order, req.Id); result.Error != nil {
		return &pb.UpdateOrderResponse{
			Status: 404,
			Error:  result.Error.Error(),
		}, nil
	}

	order.ProductID = req.ProductId
	order.Quantity = req.Quantity
	order.UserID = req.UserId
	order.UpdatedAt = time.Now()

	if result := s.H.DB.Save(&order); result.Error != nil {
		return &pb.UpdateOrderResponse{
			Status: 500,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.UpdateOrderResponse{
		Status: 200,
		Order: &pb.Order{
			Id:        order.ID,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Quantity:  order.Quantity,
			Price:     order.Price,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
			UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *Server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	if result := s.H.DB.Delete(&models.Order{}, req.Id); result.Error != nil {
		return &pb.DeleteOrderResponse{
			Status:  500,
			Error:   result.Error.Error(),
			Success: false,
		}, nil
	}

	return &pb.DeleteOrderResponse{
		Status:  200,
		Success: true,
	}, nil
}
