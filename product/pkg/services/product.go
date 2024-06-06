package services

import (
	"context"
	"net/http"
	"strings"

	"product/pkg/db"
	"product/pkg/models"
	pb "product/pkg/pb"
)

type Server struct {
	H db.Handler
	pb.UnimplementedProductServiceServer
}

func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	var product models.Product

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if result := s.H.DB.Create(&product); result.Error != nil {
		return &pb.CreateProductResponse{
			Status: http.StatusConflict,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.CreateProductResponse{
		Status: http.StatusCreated,
		Id:     product.Id,
	}, nil
}

func (s *Server) FindOne(ctx context.Context, req *pb.FindOneRequest) (*pb.FindOneResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.Id); result.Error != nil {
		return &pb.FindOneResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &pb.FindOneData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.FindOneResponse{
		Status: http.StatusOK,
		Data:   data,
	}, nil
}

func (s *Server) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.DecreaseStockResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.Id); result.Error != nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	if product.Stock <= 0 {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock too low",
		}, nil
	}

	var log models.StockDecreaseLog

	if result := s.H.DB.Where(&models.StockDecreaseLog{OrderId: req.OrderId}).First(&log); result.Error == nil {
		return &pb.DecreaseStockResponse{
			Status: http.StatusConflict,
			Error:  "Stock already decreased",
		}, nil
	}

	product.Stock = product.Stock - 1

	s.H.DB.Save(&product)

	log.OrderId = req.OrderId
	log.ProductRefer = product.Id

	s.H.DB.Create(&log)

	return &pb.DecreaseStockResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var products []models.Product
	query := s.H.DB

	if req.Filter != "" {
		query = query.Where("name LIKE ?", "%"+req.Filter+"%")
	}

	if req.SortBy != "" {
		order := req.SortBy
		if req.SortOrder != "" {
			order = order + " " + strings.ToUpper(req.SortOrder)
		}
		query = query.Order(order)
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(int(offset)).Limit(int(req.PageSize))
	}

	if result := query.Find(&products); result.Error != nil {
		return &pb.ListProductsResponse{
			Status: http.StatusInternalServerError,
			Error:  result.Error.Error(),
		}, nil
	}

	var productData []*pb.ProductData
	for _, product := range products {
		productData = append(productData, &pb.ProductData{
			Id:    product.Id,
			Name:  product.Name,
			Stock: product.Stock,
			Price: product.Price,
		})
	}

	return &pb.ListProductsResponse{
		Status:   http.StatusOK,
		Products: productData,
		Total:    total,
	}, nil
}

func (s *Server) ReadProduct(ctx context.Context, req *pb.ReadProductRequest) (*pb.ReadProductResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.Id); result.Error != nil {
		return &pb.ReadProductResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	data := &pb.ProductData{
		Id:    product.Id,
		Name:  product.Name,
		Stock: product.Stock,
		Price: product.Price,
	}

	return &pb.ReadProductResponse{
		Status:  http.StatusOK,
		Product: data,
	}, nil
}

func (s *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	var product models.Product

	if result := s.H.DB.First(&product, req.Id); result.Error != nil {
		return &pb.UpdateProductResponse{
			Status: http.StatusNotFound,
			Error:  result.Error.Error(),
		}, nil
	}

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price

	if result := s.H.DB.Save(&product); result.Error != nil {
		return &pb.UpdateProductResponse{
			Status: http.StatusInternalServerError,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.UpdateProductResponse{
		Status: http.StatusOK,
	}, nil
}

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	if result := s.H.DB.Delete(&models.Product{}, req.Id); result.Error != nil {
		return &pb.DeleteProductResponse{
			Status: http.StatusInternalServerError,
			Error:  result.Error.Error(),
		}, nil
	}

	return &pb.DeleteProductResponse{
		Status: http.StatusOK,
	}, nil
}
