package product

import (
	"api/pkg/auth"
	"api/pkg/config"
	"api/pkg/product/routes"
	"github.com/gin-gonic/gin"
	"strconv"
)

func RegisterRoutes(r *gin.Engine, c *config.Config, authSvc *auth.ServiceClient) {
	a := auth.InitAuthMiddleware(authSvc)

	svc := &ServiceClient{
		Client: InitServiceClient(c),
	}

	productRoutes := r.Group("/product")
	productRoutes.Use(a.AuthRequired)
	productRoutes.POST("/", svc.CreateProduct)
	productRoutes.GET("/:id", svc.FindOne)
	productRoutes.GET("/", svc.ListProducts)
	productRoutes.PUT("/:id", svc.UpdateProduct)
	productRoutes.DELETE("/:id", svc.DeleteProduct)
}

func (svc *ServiceClient) FindOne(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}
	routes.FineOne(ctx, svc.Client, id)
}

func (svc *ServiceClient) CreateProduct(ctx *gin.Context) {
	routes.CreateProduct(ctx, svc.Client)
}

func (svc *ServiceClient) ListProducts(ctx *gin.Context) {
	routes.ListProducts(ctx, svc.Client)
}

func (svc *ServiceClient) UpdateProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}
	routes.UpdateProduct(ctx, svc.Client, id)
}

func (svc *ServiceClient) DeleteProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}
	routes.DeleteProduct(ctx, svc.Client, id)
}
