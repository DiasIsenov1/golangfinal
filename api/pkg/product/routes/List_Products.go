// list_products.go
package routes

import (
	"api/pkg/product/pb"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ListProducts(ctx *gin.Context, c pb.ProductServiceClient) {
	var req pb.ListProductsRequest

	// Add query parameters for pagination, filtering, and sorting
	req.Filter = ctx.Query("filter")
	req.SortBy = ctx.Query("sort_by")
	req.SortOrder = ctx.Query("sort_order")
	req.Page, _ = strconv.ParseInt(ctx.Query("page"), 10, 32)
	req.PageSize, _ = strconv.ParseInt(ctx.Query("page_size"), 10, 32)

	res, err := c.ListProducts(context.Background(), &req)

	if err != nil || res.Status != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": res.Error})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
