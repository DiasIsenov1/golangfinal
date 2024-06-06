package routes

import (
	"api/pkg/product/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FineOne(ctx *gin.Context, client pb.ProductServiceClient, id int64) {
	req := &pb.FindOneRequest{
		Id: id,
	}

	res, err := client.FindOne(context.Background(), req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(int(res.Status), res)
}
