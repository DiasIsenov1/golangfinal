// update_product.go
package routes

import (
	"api/pkg/product/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UpdateProduct(ctx *gin.Context, client pb.ProductServiceClient, id int64) {
	var req pb.UpdateProductRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Id = id

	res, err := client.UpdateProduct(context.Background(), &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), res)
}
