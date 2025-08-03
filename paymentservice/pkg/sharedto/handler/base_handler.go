package handler

import (
	"fmt"
	"specommerce/paymentservice/pkg/pagination"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	pageSizeMax = 1000
)

func ParsePagination(ctx *gin.Context, paging *pagination.Paging) error {
	errTemplate := "invalid pagination parameter: %s %w"
	paging.Number = 1
	paging.Size = 10

	if sizeStr := ctx.Query("size"); sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size == 0 || uint(size) > pageSizeMax {
			return fmt.Errorf(errTemplate, "size", err)
		}
		paging.Size = uint(size)
	}
	if numberStr := ctx.Query("page"); numberStr != "" {
		pageNumber, err := strconv.Atoi(numberStr)
		if err != nil || pageNumber < 0 {
			return fmt.Errorf(errTemplate, "page", err)
		}
		if pageNumber > 0 {
			paging.Number = uint(pageNumber)
		}
	}
	orders := pagination.Orders{}
	if sortQuery := ctx.Query("sort"); sortQuery != "" {
		sort := strings.Split(sortQuery, ",")
		for _, str := range sort {
			if strings.HasPrefix(str, "-") {
				if len(str) == 1 {
					continue
				}
				orders.Add(pagination.Order{Direction: pagination.DirectionDesc, ColumnName: str[1:]})
			} else {
				orders.Add(pagination.Order{Direction: pagination.DirectionAsc, ColumnName: str})
			}
		}
	}
	paging.Sort = orders
	return nil
}