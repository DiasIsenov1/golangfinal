// delete_product.go
package routes

import (
	"api/pkg/product/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteProduct(ctx *gin.Context, client pb.ProductServiceClient, id int64) {
	req := &pb.DeleteProductRequest{
		Id: id,
	}

	res, err := client.DeleteProduct(context.Background(), req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), res)
}
